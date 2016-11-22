package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joneskoo/etget/energiatili"
	"github.com/lib/pq"
)

func main() {
	ignoreMissing := flag.Bool("ignore-missing", false, "ignore missing records")
	input := flag.String("input", "power.json", "input file name")
	flag.Parse()

	f, err := os.OpenFile(*input, os.O_RDONLY, 0)
	if err != nil {
		log.Fatalf("ERROR opening data file: %s", err)
	}

	cr, err := energiatili.FromJSON(f)
	if err != nil {
		log.Fatalf("ERROR parsing JSON structure: %s", err)
	}

	points, err := cr.DataPoints()
	switch err {
	case nil: // OK
	case energiatili.ErrorMissingRecord:
		if !*ignoreMissing {
			log.Printf("ERROR parsing data: %s", err)
			log.Printf("To ignore error, use --ignore-missing")
			os.Exit(1)
		}
	default:
		log.Fatalf("ERROR parsing data: %s", err)
	}

	rowsAffected, err := importPoints(points)
	if err != nil {
		log.Fatalf("ERROR importing to database: %s", err)
	}

	fmt.Printf("Loaded %d new rows\n", rowsAffected)
}

func importPoints(points []energiatili.DataPoint) (rowsAffected int64, err error) {
	targetTable := "energiatili"
	tmpTable := fmt.Sprintf("_%s_tmp", targetTable)
	db, err := sql.Open("postgres", "sslmode=disable")
	if err != nil {
		return 0, fmt.Errorf("connect to database: %s", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return 0, fmt.Errorf("test database connection: %s", err)
	}

	txn, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %s", err)
	}

	// Create an empty temporary table identical to target
	_, err = txn.Exec(fmt.Sprintf("CREATE TEMP TABLE %s ON COMMIT DROP AS SELECT * FROM %s WITH NO DATA", pq.QuoteIdentifier(tmpTable), pq.QuoteIdentifier(targetTable)))
	if err != nil {
		return 0, fmt.Errorf("create temporary table: %s", err)
	}

	// Load data into temporary table
	stmt, err := txn.Prepare(pq.CopyIn(tmpTable, "ts", "kwh"))
	if err != nil {
		return 0, fmt.Errorf("copy data into temporary table: %s", err)
	}
	for _, point := range points {
		_, err = stmt.Exec(point.Time.UTC(), point.Kwh)
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

	// Copy data from temporary table into target
	res, err := txn.Exec(fmt.Sprintf("INSERT INTO %s (ts, kwh) SELECT ts, kwh FROM %s ON CONFLICT DO NOTHING", pq.QuoteIdentifier(targetTable), pq.QuoteIdentifier(tmpTable)))
	if err != nil {
		return 0, fmt.Errorf("load data from temporary table: %s", err)
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	err = txn.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit transaction: %s", err)
	}
	return
}
