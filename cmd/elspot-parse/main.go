// elspot-parse imports a ".xls" market data file from Nordpool to PostgreSQL.
// Actually it's not an Excel document but a HTML table, that happens to load
// in Excel. Luckily it's easier to parse than the Excel file would have been.
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"time"

	"github.com/joneskoo/etget/htmltable"
	"github.com/lib/pq"
)

// location defines the time zone of the data being imported
var location *time.Location

func init() {
	var err error
	if location, err = time.LoadLocation("Europe/Paris"); err != nil {
		panic(fmt.Sprintf("ERROR loading timezone: %s", err))
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s FILENAME\n\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "where FILENAME is elspot 'xls' file\n")
	flag.PrintDefaults()
	os.Exit(1)
}

var traceTimings bool

func main() {
	flag.BoolVar(&traceTimings, "trace", false, "trace execution time")
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
	}

	progress := timer{time.Now()}

	f, err := os.OpenFile(flag.Arg(0), os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("ERROR opening data file: %s", err)
	}

	progress.Track("open file")

	tables, err := htmltable.Parse(f)
	if err != nil {
		log.Fatalf("ERROR parsing HTML table: %s", err)
	}

	progress.Track("parse html")

	records, err := parseTable(tables[0])
	if err != nil {
		log.Fatalf("ERROR parsing elspot table: %s", err)
	}

	progress.Track("parse table")

	rowsAffected, err := loadToPostgres(records)
	if err != nil {
		log.Fatalf("ERROR importing to PostgreSQL: %s", err)
	}

	progress.Track("load to postgres")

	fmt.Printf("OK! %d rows affected\n", rowsAffected)
}

type record struct {
	Time   time.Time
	Prices map[string]string
}

func parseTable(table htmltable.Table) (records []record, err error) {
	header := table.Headers[2]

	commaToPeriod := strings.NewReplacer(",", ".")

	for _, t := range table.Rows {
		prices := make(map[string]string, len(header)-1)
		for i, k := range header {
			prices[k] = commaToPeriod.Replace(t[i])
		}

		// Date is t[0], and hour is t[1], in CET timezone (Europe/Paris)
		ts, err := time.ParseInLocation("02-01-2006 15", fmt.Sprintf("%s %s", t[0], t[1][0:2]), location)
		if err != nil {
			return nil, err
		}

		records = append(records, record{
			Time:   ts,
			Prices: prices,
		})
	}
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

func loadToPostgres(records []record) (rowsAffected int64, err error) {
	progress := timer{time.Now()}

	db, err := sql.Open("postgres", "sslmode=disable")
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
		_, err = stmt.Exec(r.Time, r.Prices["FI"])
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
