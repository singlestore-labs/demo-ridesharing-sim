package service

import (
	"server/database"
	"strings"
)

func GetMinuteAvgWaitTimeLastHour(db, city string) []map[string]interface{} {
	var query string
	var args []interface{}

	if db == "snowflake" {
		database.SetupSnowflakeQuery()
		query = `
			WITH minute_avg_wait AS (
				SELECT 
					DATE_TRUNC('MINUTE', request_time) AS minute_interval,
					AVG(DATEDIFF('SECOND', request_time, accept_time)) AS avg_wait_time
				FROM 
					trips
				WHERE 
					request_time >= DATEADD(HOUR, -1, CURRENT_TIMESTAMP())
					AND accept_time IS NOT NULL
					{{ city_filter }}
				GROUP BY 
					minute_interval
			)
			SELECT 
				TO_CHAR(w.minute_interval, 'YYYY-MM-DD HH24:MI:00') AS minute_interval,
				ROUND(w.avg_wait_time, 2) AS avg_wait_time,
				COALESCE(
					ROUND(
						(w.avg_wait_time - LAG(w.avg_wait_time) OVER (ORDER BY w.minute_interval)) / 
						NULLIF(LAG(w.avg_wait_time) OVER (ORDER BY w.minute_interval), 0) * 100,
						2
					),
					0
				) AS percent_change
			FROM 
				minute_avg_wait w
			ORDER BY 
				w.minute_interval;
		`
	} else {
		query = `
			WITH minute_avg_wait AS (
				SELECT 
					DATE_FORMAT(request_time, '%Y-%m-%d %H:%i:00') AS minute_interval,
					AVG(TIMESTAMPDIFF(SECOND, request_time, accept_time)) AS avg_wait_time
				FROM 
					trips
				WHERE 
					request_time >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
					AND accept_time IS NOT NULL
					{{ city_filter }}
				GROUP BY 
					minute_interval
			)
			SELECT 
				w.minute_interval,
				ROUND(w.avg_wait_time, 2) AS avg_wait_time,
				COALESCE(
					ROUND(
						(w.avg_wait_time - LAG(w.avg_wait_time) OVER (ORDER BY w.minute_interval)) / 
						NULLIF(LAG(w.avg_wait_time) OVER (ORDER BY w.minute_interval), 0) * 100,
						2
					),
					0
				) AS percent_change
			FROM 
				minute_avg_wait w
			ORDER BY 
				w.minute_interval;
		`
	}

	// Replace placeholders based on whether city is provided
	if city != "" {
		query = strings.ReplaceAll(query, "{{ city_filter }}", "AND city = ?")
		args = append(args, city)
	} else {
		query = strings.ReplaceAll(query, "{{ city_filter }}", "")
	}

	var results = make([]map[string]interface{}, 0)

	if db == "snowflake" {
		rows, err := database.SnowflakeDB.Query(query, args...)
		if err != nil {
			return nil
		}
		defer rows.Close()

		for rows.Next() {
			var minuteInterval string
			var avgWaitTime float64
			var percentChange float64

			if err := rows.Scan(&minuteInterval, &avgWaitTime, &percentChange); err != nil {
				return nil
			}

			result := map[string]interface{}{
				"minute_interval": minuteInterval,
				"avg_wait_time":   avgWaitTime,
				"percent_change":  percentChange,
			}
			results = append(results, result)
		}

		if err = rows.Err(); err != nil {
			return nil
		}
	} else {
		if err := database.SingleStoreDB.Raw(query, args...).Scan(&results).Error; err != nil {
			return nil
		}
	}

	return results
}
