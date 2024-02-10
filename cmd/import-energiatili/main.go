// The import-energiatili command imports data from www.energiatili.fi.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"encoding/json"

	"github.com/joneskoo/etget/energiatili"
	"github.com/joneskoo/etget/keyring"
	"github.com/lib/pq"
)

func main() {
	credfile := flag.String("credfile", "./credentials.json", "File username/password are saved in (plaintext)")
	connstring := flag.String("connstring", "sslmode=disable", "https://www.postgresql.org/docs/current/static/libpq-connect.html#LIBPQ-CONNSTRING")
	flag.Parse()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cs := keyring.CredentialStore{
		File:   *credfile,
		Domain: "www.energiatili.fi",
	}

	client := &energiatili.Client{
		UsernamePasswordFunc: cs.UsernamePassword,
	}

	consumptionReportFile := "consumptionreport.json"

	// Download data from API
	var f *os.File
	var err error

	f, err = os.Open(consumptionReportFile)
	if os.IsNotExist(err) {
		log.Println("Downloading consumption data…")
		f, err = os.Create(consumptionReportFile)
		if err != nil {
			panic(err)
		}
		if err := client.ConsumptionReport(ctx, f); err != nil {
			log.Fatalln(err)
		}
		f.Seek(0, 0)
	} else {
		log.Println("Using cached consumption data:", consumptionReportFile)
	}
	defer f.Close()

	if err != nil {
		panic(err)
	}

	var consumptionreport energiatili.ConsumptionReport
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&consumptionreport)
	if err != nil {
		log.Fatalf("ERROR parsing JSON structure: %s", err)
	}
	points, err := consumptionreport.Records()
	if err != nil {
		log.Fatalf("ERROR parsing data: %s", err)
	}

	rowsAffected, err := importPoints(*connstring, points)
	if err != nil {
		log.Fatalf("ERROR importing to database: %s", err)
	}

	log.Printf("Loaded %d new rows", rowsAffected)
}

func importPoints(connstring string, points []energiatili.Record) (rowsAffected int64, err error) {
	targetTable := "energiatili"
	tmpTable := fmt.Sprintf("_%s_tmp", targetTable)

	db, err := sql.Open("postgres", connstring)
	if err != nil {
		return 0, fmt.Errorf("connect to database: %s", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		return 0, fmt.Errorf("test database connection: %s", err)
	}

	// Ensure table exists
	_, err = db.Exec(createTable)
	if err != nil {
		return 0, fmt.Errorf("ensure table exists: %s", err)
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
		_, err = stmt.Exec(point.Timestamp.UTC(), point.Value)
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
