package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/joneskoo/etget/energiatili"
	"github.com/lib/pq"
)

func main() {
	f, err := os.OpenFile("power.json", os.O_RDONLY, 0)
	if err != nil {
		log.Fatalln(err)
	}
	points, err := parseHours(f)
	if err != nil {
		log.Fatalln(err)
	}
	rowsAffected, err := importPoints(points)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Loaded %d new rows\n", rowsAffected)
}

type point struct {
	Time time.Time
	Kwh  float64
}

func parseHours(r io.Reader) (points []point, err error) {
	cr, err := energiatili.FromJSON(r)
	if err != nil {
		return nil, err
	}
	if len(cr.Hours.Consumptions) != 2 {
		return nil, fmt.Errorf("len(Consumptions) = %d, want 2", len(cr.Hours.Consumptions))
	}
	meterA := cr.Hours.Consumptions[0].Series.Data
	meterB := cr.Hours.Consumptions[1].Series.Data
	points = make([]point, len(meterA)+len(meterB))
	a := 0
	b := 0
	for i := 0; i < len(meterA)+len(meterB); i++ {
		if b >= len(meterB) || a < len(meterA) && meterA[a].Time.Before(meterB[b].Time) {
			points[i] = point{
				Time: meterA[a].Time,
				Kwh:  meterA[a].Kwh,
			}
			a++
		} else {
			points[i] = point{
				Time: meterB[b].Time,
				Kwh:  meterB[b].Kwh,
			}
			b++
		}
	}
	return points, nil
}

func importPoints(points []point) (rowsAffected int64, err error) {
	targetTable := "energiatili"
	tmpTable := fmt.Sprintf("_%s_tmp", targetTable)
	db, err := sql.Open("postgres", "sslmode=disable")
	if err != nil {
		return
	}
	defer db.Close()

	txn, err := db.Begin()
	if err != nil {
		return
	}

	// Create an empty temporary table identical to target
	_, err = txn.Exec(fmt.Sprintf("CREATE TEMP TABLE %s ON COMMIT DROP AS SELECT * FROM %s WITH NO DATA", pq.QuoteIdentifier(tmpTable), pq.QuoteIdentifier(targetTable)))
	if err != nil {
		return
	}

	// Load data into temporary table
	stmt, err := txn.Prepare(pq.CopyIn(tmpTable, "ts", "kwh"))
	if err != nil {
		return
	}
	for _, point := range points {
		_, err = stmt.Exec(point.Time.UTC(), point.Kwh)
		if err != nil {
			return
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return
	}
	err = stmt.Close()
	if err != nil {
		return
	}

	// Copy data from temporary table into target
	res, err := txn.Exec(fmt.Sprintf("INSERT INTO %s (ts, kwh) SELECT ts, kwh FROM %s ON CONFLICT DO NOTHING", pq.QuoteIdentifier(targetTable), pq.QuoteIdentifier(tmpTable)))
	if err != nil {
		return
	}
	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return
	}

	err = txn.Commit()
	if err != nil {
		return
	}
	return
}
