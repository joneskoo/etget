DROP TABLE power_contracts;
CREATE TABLE power_contracts (
    id SERIAL,
    valid_from TIMESTAMP WITH TIME ZONE NOT NULL,
    currency TEXT NOT NULL DEFAULT 'EUR',
    energy_monthly NUMERIC(6, 2) NOT NULL,
    transfer_monthly NUMERIC(6, 2) NOT NULL,
    transfer_night_per_kwh NUMERIC(7, 7) NOT NULL,
    transfer_day_per_kwh NUMERIC(7, 7) NOT NULL,
    energy_margin_per_kwh NUMERIC(7, 7) NOT NULL
);

INSERT INTO power_contracts
    (valid_from, currency, energy_monthly, transfer_monthly, transfer_night_per_kwh, transfer_day_per_kwh, energy_margin_per_kwh)
VALUES
    -- 143.44/a  // 1.57 + 2.35972 //  2.57 + 2.35972
    ('2000-01-01+02',          'EUR', 3.30, 11.95, 0.0392972, 0.0492972, 0.0017),
    -- 143.44/a // 1.57 + 2.79372 // 2.57 + 2.79372
    ('2015-01-01 00:00:00+02', 'EUR', 3.30, 11.95, 0.0436372, 0.0536372, 0.0017),
    -- 185,11e/a // 1.71 + 2,79372 // 2.79 + 2,79372
    ('2016-06-01 00:00:00+03', 'EUR', 3.30, 15.43, 0.0450372, 0.0558372, 0.0017);

CREATE OR REPLACE FUNCTION is_night(ts timestamptz) RETURNS bool AS $$
    SELECT EXTRACT(HOUR FROM ts) IN (22, 23, 00, 01, 02, 03, 04, 05, 06);
$$ LANGUAGE SQL;

CREATE OR REPLACE FUNCTION is_day(ts timestamptz) RETURNS bool AS $$
    SELECT EXTRACT(HOUR FROM ts) NOT IN (22, 23, 00, 01, 02, 03, 04, 05, 06);
$$ LANGUAGE SQL;

SET timezone = "Europe/Helsinki";

CREATE VIEW power_bills AS (
    WITH prices AS (
        SELECT DISTINCT ON (elspot.ts)
            elspot.ts,
            power_contracts.*,
            transfer_day_per_kwh * is_day(ts)::int + transfer_night_per_kwh * is_night(ts)::int AS transfer_e_kwh,
            fi/1000*1.24 + energy_margin_per_kwh AS spot_e_kwh
        FROM power_contracts
            LEFT JOIN elspot ON power_contracts.valid_from <= elspot.ts
        ORDER BY elspot.ts, power_contracts.valid_from DESC
    )
    SELECT
        DATE_TRUNC('month', ts)::date AS month,
        ROUND(SUM(kwh)) AS kwh,
        ROUND(SUM(kwh * is_night(ts)::int)) AS kwh_night,
        ROUND(SUM(kwh * is_day(ts)::int)) AS kwh_day,
        ROUND((AVG(energy_monthly) + SUM(kwh * spot_e_kwh))::numeric, 2) AS energy_eur,
        ROUND((AVG(transfer_monthly) + SUM(kwh * transfer_e_kwh))::numeric, 2) AS transfer_eur,
        ROUND(100*(SUM(kwh * is_night(ts)::int * spot_e_kwh)/SUM(kwh * is_night(ts)::int))::numeric, 2) AS energy_c_night_kwh,
        ROUND(100*(SUM(kwh * is_day(ts)::int * spot_e_kwh)/SUM(kwh * is_day(ts)::int))::numeric, 2) AS energy_c_day_kwh
    FROM energiatili INNER JOIN prices USING (ts)
    GROUP BY 1
    ORDER BY 1
);

CREATE OR REPLACE VIEW power_years AS (
    WITH prices AS (
        SELECT DISTINCT ON (elspot.ts)
            elspot.ts,
            power_contracts.*,
            transfer_day_per_kwh * is_day(ts)::int + transfer_night_per_kwh * is_night(ts)::int AS transfer_e_kwh,
            fi/1000*1.24 + energy_margin_per_kwh AS spot_e_kwh
        FROM power_contracts
            LEFT JOIN elspot ON power_contracts.valid_from <= elspot.ts
        ORDER BY elspot.ts, power_contracts.valid_from DESC
    )
    SELECT
        DATE_TRUNC('year', ts)::date AS month,
        ROUND(SUM(kwh)) AS kwh,
        ROUND(SUM(kwh * is_night(ts)::int)) AS kwh_night,
        ROUND(SUM(kwh * is_day(ts)::int)) AS kwh_day,
        ROUND((12*AVG(energy_monthly) + SUM(kwh * spot_e_kwh))::numeric, 2) AS energy_eur,
        ROUND((12*AVG(transfer_monthly) + SUM(kwh * transfer_e_kwh))::numeric, 2) AS transfer_eur,
        ROUND(100*(SUM(kwh * is_night(ts)::int * spot_e_kwh)/SUM(kwh * is_night(ts)::int))::numeric, 2) AS energy_c_night_kwh,
        ROUND(100*(SUM(kwh * is_day(ts)::int * spot_e_kwh)/SUM(kwh * is_day(ts)::int))::numeric, 2) AS energy_c_day_kwh
    FROM energiatili INNER JOIN prices USING (ts)
    GROUP BY 1
    ORDER BY 1
);
