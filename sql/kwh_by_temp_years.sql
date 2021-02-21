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
