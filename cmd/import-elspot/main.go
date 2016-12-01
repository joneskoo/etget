// The import-elspot command loads market data file from Nordpool.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"time"

	"github.com/joneskoo/etget/htmltable"
	"github.com/joneskoo/etget/notz"
	"github.com/lib/pq"
)

const timeLayout = "02-01-2006 15"

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s ELSPOT\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "   ELSPOT	elspot 'xls' file name or URL\n\n")
	flag.PrintDefaults()
	os.Exit(1)
}

var traceTimings bool

func main() {
	connstring := flag.String("connstring", "sslmode=disable", "https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING")
	flag.BoolVar(&traceTimings, "trace", false, "trace execution time")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	progress := timer{time.Now()}

	var src io.ReadCloser

	// If ELSPOT is a URL, download it. If not, assume it's a file.
	if u, err := url.Parse(flag.Arg(0)); err == nil && (u.Scheme == "http" || u.Scheme == "https") {
		resp, err := http.Get(flag.Arg(0))
		if err != nil {
			log.Fatalf("ERROR opening URL: %s", err)
		}
		if resp.StatusCode != http.StatusOK {
			log.Fatalf("ERROR got HTTP status code: %d, want %d", resp.StatusCode, http.StatusOK)
		}
		src = resp.Body
	} else {
		src, err = os.OpenFile(flag.Arg(0), os.O_RDONLY, 0)
		if err != nil {
			log.Fatalf("ERROR opening data file: %s", err)
		}
	}
	defer src.Close()

	progress.Track("open file")

	tables, err := htmltable.Parse(src)
	if err != nil {
		log.Fatalf("ERROR parsing HTML table: %s", err)
	}

	progress.Track("parse html")

	records, err := parseTable(tables[0])
	if err != nil {
		log.Fatalf("ERROR parsing elspot table: %s", err)
	}

	progress.Track("parse table")

	rowsAffected, err := loadToPostgres(*connstring, records)
	if err != nil {
		log.Fatalf("ERROR importing to PostgreSQL: %s", err)
	}

	progress.Track("load to postgres")

	fmt.Printf("OK! %d rows affected\n", rowsAffected)
}

type record struct {
	Timestamp time.Time
	Prices    map[string]string
}

// records implements notz.Interface for notz.FixDST.
type records []record

func (r records) Len() int                     { return len(r) }
func (r records) Time(i int) time.Time         { return r[i].Timestamp }
func (r records) SetTime(i int, new time.Time) { r[i].Timestamp = new }

func parseTable(table htmltable.Table) (data []record, err error) {
	var loc *time.Location
	loc, err = time.LoadLocation("Europe/Paris")
	if err != nil {
		return nil, err

	}
	header := table.Headers[2]

	commaToPeriod := strings.NewReplacer(",", ".")

	for _, t := range table.Rows {
		prices := make(map[string]string, len(header)-1)
		for i, k := range header {
			prices[k] = commaToPeriod.Replace(t[i])
		}
		if prices["SYS"] == "" {
			continue
		}

		// Date is t[0], and hour is first two bytes of t[1]
		ts, err := time.ParseInLocation(timeLayout, fmt.Sprintf("%s %s", t[0], t[1][0:2]), loc)
		if err != nil {
			return nil, fmt.Errorf("parsing timestamp: %s", err)
		}

		data = append(data, record{
			Timestamp: ts,
			Prices:    prices,
		})
	}
	notz.FixDST(records(data))
	return
}

type timer struct{ time.Time }

func (t *timer) Track(msg string) {
	if !traceTimings {
		return
	}
	if t.IsZero() {
		t.Time = time.Now()
	}
	fmt.Println(msg, "took", time.Now().Sub(t.Time))
	t.Time = time.Now()
}

func loadToPostgres(connstring string, records []record) (rowsAffected int64, err error) {
	progress := timer{time.Now()}

	db, err := sql.Open("postgres", connstring)
	if err != nil {
		return 0, fmt.Errorf("connect to database: %s", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return 0, fmt.Errorf("test database connection: %s", err)
	}

	progress.Track("connect to database")

	// Ensure table exists
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return 0, fmt.Errorf("ensure table exists: %s", err)
	}

	progress.Track("table exists")

	txn, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %s", err)
	}

	progress.Track("begin transaction")

	// Create an empty temporary table identical to target
	_, err = txn.Exec(fmt.Sprintf("CREATE TEMP TABLE %s ON COMMIT DROP AS SELECT * FROM %s WITH NO DATA", pq.QuoteIdentifier(tmpTable), pq.QuoteIdentifier(targetTable)))
	if err != nil {
		return 0, fmt.Errorf("create temporary table: %s", err)
	}

	progress.Track("create temp table")

	// Load data into temporary table
	stmt, err := txn.Prepare(pq.CopyIn(tmpTable, "ts", "fi"))
	if err != nil {
		return 0, fmt.Errorf("copy data into temporary table: %s", err)
	}
	for _, r := range records {
		if r.Prices["FI"] == "" {
			continue
		}
		_, err = stmt.Exec(r.Timestamp, r.Prices["FI"])
		if err != nil {
			return 0, fmt.Errorf("insert data into temporary table: %s", err)
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return 0, fmt.Errorf("flush after loading data: %s", err)
	}
	err = stmt.Close()
	if err != nil {
		return
	}

	progress.Track("load data into temp table")

	// Copy data from temporary table into target
	res, err := txn.Exec(fmt.Sprintf("INSERT INTO %s (ts, FI) SELECT ts, FI FROM %s ON CONFLICT DO NOTHING", pq.QuoteIdentifier(targetTable), pq.QuoteIdentifier(tmpTable)))
	if err != nil {
		return 0, fmt.Errorf("load data from temporary table: %s", err)
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	progress.Track("copy data to target table")

	err = txn.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit transaction: %s", err)
	}

	progress.Track("commit transaction")

	return
}
