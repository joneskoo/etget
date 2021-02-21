CREATE OR REPLACE VIEW kwh_by_temp_years AS (
	-- Select average kWh/h consumed based on Willab weather station temperature.
	-- Calculate the average temperature and average consumption of each day to
	-- smooth out noise, and group by temperature, segmented into 5°C ranges.
	SELECT
		FLOOR(temp.temp/5)*5 AS temp_range_lo,
		FLOOR(temp.temp/5)*5+5 AS temp_range_hi,
		EXTRACT(year FROM temp.date) AS year,
		ROUND(AVG(kwh.kwh)::numeric, 3) AS kwh_avg,
		COUNT(kwh.kwh) as days_averaged
	FROM
		(
		-- Average temperature by day
			SELECT
				DATE(measuretime AT TIME ZONE 'Europe/Helsinki') AS date,
				AVG(tempnow) AS temp
			FROM willab_weather
			GROUP BY 1
		) AS temp,
		(
			-- Average consumption by day
			SELECT
				DATE(ts AT TIME ZONE 'Europe/Helsinki') AS date,
				AVG(kwh) AS kwh
			FROM energiatili
			GROUP BY 1
		) AS kwh
	WHERE temp.date = kwh.date AND temp.date >= '2018-10-01'
	GROUP BY 1,3
	ORDER BY 1,3
);

-- SELECT
--         year,
--         MAX(CASE WHEN temp_range_lo = -25 THEN kwh_avg END) "-25°C…-20°C",
--         MAX(CASE WHEN temp_range_lo = -20 THEN kwh_avg END) "-20°C…-15°C",
--         MAX(CASE WHEN temp_range_lo = -15 THEN kwh_avg END) "-15°C…-10°C",
--         MAX(CASE WHEN temp_range_lo = -10 THEN kwh_avg END) "-10°C…-5°C",
--         MAX(CASE WHEN temp_range_lo = -5 THEN kwh_avg END) "-5°C…0°C",
--         MAX(CASE WHEN temp_range_lo = 0 THEN kwh_avg END) "+0°C…+5°C",
--         MAX(CASE WHEN temp_range_lo = 5 THEN kwh_avg END) "+5°C…+10°C",
--         MAX(CASE WHEN temp_range_lo = 10 THEN kwh_avg END) "+10°C…+15°C",
--         MAX(CASE WHEN temp_range_lo = 15 THEN kwh_avg END) "+15°C…+20°C",
--         MAX(CASE WHEN temp_range_lo = 20 THEN kwh_avg END) "+20°C…+25°C"
-- FROM kwh_by_temp_years
--  year | -25°C…-20°C | -20°C…-15°C | -15°C…-10°C | -10°C…-5°C | -5°C…0°C | +0°C…+5°C | +5°C…+10°C | +10°C…+15°C | +15°C…+20°C | +20°C…+25°C
-- ------+-------------+-------------+-------------+------------+----------+-----------+------------+-------------+-------------+-------------
--  2018 |             |       4.885 |             |      3.269 |    2.909 |     2.297 |      2.123 |       1.676 |       1.632 |
--  2019 |       5.080 |       4.554 |       3.834 |      3.163 |    2.578 |     2.169 |      1.583 |       1.136 |       1.064 |       0.969
--  2020 |             |             |       3.696 |      3.211 |    2.585 |     2.148 |      1.567 |       1.155 |       1.117 |       1.219
--  2021 |       4.785 |       4.556 |       3.683 |      3.353 |    2.831 |           |            |             |             |
-- (4 rows)


-- DROP TABLE IF EXISTS willab_weather;
-- CREATE TABLE willab_weather (
-- 	measuretime timestamp with time zone,
-- 	tempnow real,
-- 	temphi real,
-- 	templo real,
-- 	dewpoint real,
-- 	humidity real,
-- 	airpressure real,
-- 	windspeed real,
-- 	windspeedmax real,
-- 	winddir real,
-- 	precipitation1h real,
-- 	precipitation24h real,
-- 	solarrad real
-- );
-- -- http://www.ipv6.willab.fi/weather/weather_csv.zip
-- \copy willab_weather from '$tmp/all.csv' DELIMITER ',' CSV header
