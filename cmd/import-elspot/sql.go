package main

const (
	targetTable = "elspot"
	tmpTable    = "_elspot_tmp"

	createTableSQL = `CREATE TABLE IF NOT EXISTS elspot (
    id      SERIAL,
    ts      TIMESTAMPTZ UNIQUE,
    FI      REAL
    );`
)
