package main

const (
	dropTable = `DROP TABLE IF EXISTS energiatili;`

	createTable = `CREATE TABLE IF NOT EXISTS energiatili (
    id SERIAL,
    ts timestamptz unique,
    kwh double precision,
    temp real);`
)
