package main

const createTable = `CREATE TABLE IF NOT EXISTS energiatili (
    id SERIAL,
    ts timestamptz unique,
    kwh double precision,
    temp real);`
