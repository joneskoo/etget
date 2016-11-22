DROP TABLE IF EXISTS energiatili;
CREATE TABLE energiatili (id SERIAL, ts timestamptz unique, kwh double precision, temp real);
