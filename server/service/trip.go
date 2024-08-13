package service

import (
	"fmt"
	"server/database"
	"strings"
)

func GetCities(db string) []string {
	var cities []string
	if db == "snowflake" {
		rows, err := database.SnowflakeDB.Query("SELECT DISTINCT city FROM trips")
		if err != nil {
			return []string{}
		}
		defer rows.Close()

		for rows.Next() {
			var city string
			if err := rows.Scan(&city); err != nil {
				continue
			}
			cities = append(cities, city)
		}

		if err = rows.Err(); err != nil {
			fmt.Println("Error iterating over rows:", err)
			return []string{}
		}
	} else {
		err := database.SingleStoreDB.Raw("SELECT DISTINCT city FROM trips").Scan(&cities).Error
		if err != nil {
			return []string{}
		}
	}
	return cities
}

func GetCurrentTripStatus(db string) map[string]interface{} {
	result := map[string]interface{}{
		"trips_requested":     0,
		"trips_accepted":      0,
		"trips_en_route":      0,
		"riders_idle":         0,
		"riders_requested":    0,
		"riders_waiting":      0,
		"riders_in_progress":  0,
		"drivers_available":   0,
		"drivers_in_progress": 0,
	}

	if db == "snowflake" {
		query := `
			SELECT 'trips' as entity, status, COUNT(*) as count
				FROM trips
				GROUP BY status
				UNION ALL
				SELECT 'riders' as entity, status, COUNT(*) as count
				FROM riders
				GROUP BY status
				UNION ALL
				SELECT 'drivers' as entity, status, COUNT(*) as count
				FROM drivers
				GROUP BY status
				ORDER BY entity, status;
		`

		rows, err := database.SnowflakeDB.Query(query)
		if err != nil {
			return result
		}
		defer rows.Close()

		for rows.Next() {
			var entity, status string
			var count int
			if err := rows.Scan(&entity, &status, &count); err != nil {
				continue
			}
			key := fmt.Sprintf("%s_%s", entity, status)
			if _, exists := result[key]; exists {
				result[key] = count
			}
		}

		if err = rows.Err(); err != nil {
			fmt.Println("Error iterating over rows:", err)
		}

		return result
	}

	query := `
		SELECT 'trips' as entity, status, COUNT(*) as count
			FROM trips
			GROUP BY status
			UNION ALL
			SELECT 'riders' as entity, status, COUNT(*) as count
			FROM riders
			GROUP BY status
			UNION ALL
			SELECT 'drivers' as entity, status, COUNT(*) as count
			FROM drivers
			GROUP BY status
			ORDER BY entity, status;
	`

	var results []struct {
		Entity string
		Status string
		Count  int
	}

	err := database.SingleStoreDB.Raw(query).Scan(&results).Error
	if err != nil {
		return result
	}

	for _, r := range results {
		key := fmt.Sprintf("%s_%s", r.Entity, r.Status)
		if _, exists := result[key]; exists {
			result[key] = r.Count
		}
	}

	return result
}

func GetCurrentTripStatusByCity(db string, city string) map[string]interface{} {
	result := map[string]interface{}{
		"trips_requested":     0,
		"trips_accepted":      0,
		"trips_en_route":      0,
		"riders_idle":         0,
		"riders_requested":    0,
		"riders_waiting":      0,
		"riders_in_progress":  0,
		"drivers_available":   0,
		"drivers_in_progress": 0,
	}

	if db == "snowflake" {
		query := `
			SELECT 'trips' as entity, status, COUNT(*) as count
				FROM trips
				WHERE city = ?
				GROUP BY status
				UNION ALL
				SELECT 'riders' as entity, status, COUNT(*) as count
				FROM riders
				WHERE location_city = ?
				GROUP BY status
				UNION ALL
				SELECT 'drivers' as entity, status, COUNT(*) as count
				FROM drivers
				WHERE location_city = ?
				GROUP BY status
				ORDER BY entity, status;
		`

		rows, err := database.SnowflakeDB.Query(query, city, city, city)
		if err != nil {
			return result
		}
		defer rows.Close()

		for rows.Next() {
			var entity, status string
			var count int
			if err := rows.Scan(&entity, &status, &count); err != nil {
				continue
			}
			key := fmt.Sprintf("%s_%s", entity, status)
			if _, exists := result[key]; exists {
				result[key] = count
			}
		}

		if err = rows.Err(); err != nil {
			fmt.Println("Error iterating over rows:", err)
		}

		return result
	}

	query := `
		SELECT 'trips' as entity, status, COUNT(*) as count
			FROM trips
			WHERE city = ?
			GROUP BY status
			UNION ALL
			SELECT 'riders' as entity, status, COUNT(*) as count
			FROM riders
			WHERE location_city = ?
			GROUP BY status
			UNION ALL
			SELECT 'drivers' as entity, status, COUNT(*) as count
			FROM drivers
			WHERE location_city = ?
			GROUP BY status
			ORDER BY entity, status;
	`
	var results []struct {
		Entity string
		Status string
		Count  int
	}

	err := database.SingleStoreDB.Raw(query, city, city, city).Scan(&results).Error
	if err != nil {
		return result
	}

	for _, r := range results {
		key := fmt.Sprintf("%s_%s", r.Entity, r.Status)
		if _, exists := result[key]; exists {
			result[key] = r.Count
		}
	}

	return result
}

func GetMinuteTripCountsLastHour(db, city string) []map[string]interface{} {
	var query string
	var args []interface{}

	if db == "snowflake" {
		query = `
			WITH minute_counts AS (
				SELECT 
					DATE_TRUNC('MINUTE', request_time) AS minute_interval,
					COUNT(*) AS trip_count
				FROM 
					trips
				WHERE 
					request_time >= DATEADD(HOUR, -1, CURRENT_TIMESTAMP())
					{{ city_filter }}
				GROUP BY 
					minute_interval
			)
			SELECT 
				TO_CHAR(c.minute_interval, 'YYYY-MM-DD HH24:MI:00') AS minute_interval,
				c.trip_count,
				COALESCE(
					ROUND(
						(c.trip_count - LAG(c.trip_count) OVER (ORDER BY c.minute_interval)) / 
						NULLIF(LAG(c.trip_count) OVER (ORDER BY c.minute_interval), 0) * 100,
						2
					),
					0
				) AS percent_change
			FROM 
				minute_counts c
			ORDER BY 
				c.minute_interval;
		`
	} else {
		query = `
			WITH minute_counts AS (
				SELECT 
					DATE_FORMAT(request_time, '%Y-%m-%d %H:%i:00') AS minute_interval,
					COUNT(*) AS trip_count
				FROM 
					trips
				WHERE 
					request_time >= DATE_SUB(NOW(), INTERVAL 1 HOUR)
					{{ city_filter }}
				GROUP BY 
					minute_interval
			)
			SELECT 
				c.minute_interval,
				c.trip_count,
				COALESCE(
					ROUND(
						(c.trip_count - LAG(c.trip_count) OVER (ORDER BY c.minute_interval)) / 
						NULLIF(LAG(c.trip_count) OVER (ORDER BY c.minute_interval), 0) * 100,
						2
					),
					0
				) AS percent_change
			FROM 
				minute_counts c
			ORDER BY 
				c.minute_interval;
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
			var tripCount int
			var percentChange float64

			if city != "" {
				if err := rows.Scan(&minuteInterval, &tripCount, &percentChange); err != nil {
					return nil
				}
			} else {
				if err := rows.Scan(&minuteInterval, &tripCount, &percentChange); err != nil {
					return nil
				}
			}

			result := map[string]interface{}{
				"interval":       minuteInterval,
				"trip_count":     tripCount,
				"percent_change": percentChange,
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

func GetHourlyTripCountsLastDay(db, city string) []map[string]interface{} {
	var query string
	var args []interface{}

	if db == "snowflake" {
		query = `
			WITH hourly_counts AS (
				SELECT 
					DATE_TRUNC('HOUR', request_time) AS hourly_interval,
					COUNT(*) AS trip_count
				FROM 
					trips
				WHERE 
					request_time >= DATEADD(HOUR, -24, CURRENT_TIMESTAMP())
					{{ city_filter }}
				GROUP BY 
					hourly_interval
			)
			SELECT 
				TO_CHAR(c.hourly_interval, 'YYYY-MM-DD HH24:00:00') AS hourly_interval,
				c.trip_count,
				COALESCE(
					ROUND(
						(c.trip_count - LAG(c.trip_count) OVER (ORDER BY c.hourly_interval)) / 
						NULLIF(LAG(c.trip_count) OVER (ORDER BY c.hourly_interval), 0) * 100,
						2
					),
					0
				) AS percent_change
			FROM 
				hourly_counts c
			ORDER BY 
				c.hourly_interval;
		`
	} else {
		query = `
			WITH hourly_counts AS (
				SELECT 
					DATE_FORMAT(request_time, '%Y-%m-%d %H:00:00') AS hourly_interval,
					COUNT(*) AS trip_count
				FROM 
					trips
				WHERE 
					request_time >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
					{{ city_filter }}
				GROUP BY 
					hourly_interval
			)
			SELECT 
				c.hourly_interval,
				c.trip_count,
				COALESCE(
					ROUND(
						(c.trip_count - LAG(c.trip_count) OVER (ORDER BY c.hourly_interval)) / 
						NULLIF(LAG(c.trip_count) OVER (ORDER BY c.hourly_interval), 0) * 100,
						2
					),
					0
				) AS percent_change
			FROM 
				hourly_counts c
			ORDER BY 
				c.hourly_interval;
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
			var hourlyInterval string
			var tripCount int
			var percentChange float64

			if err := rows.Scan(&hourlyInterval, &tripCount, &percentChange); err != nil {
				return nil
			}

			result := map[string]interface{}{
				"interval":       hourlyInterval,
				"trip_count":     tripCount,
				"percent_change": percentChange,
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

func GetDailyTripCountsLastWeek(db, city string) []map[string]interface{} {
	var query string
	var args []interface{}

	if db == "snowflake" {
		query = `
			WITH daily_counts AS (
				SELECT 
					DATE(request_time) AS daily_interval,
					COUNT(*) AS trip_count
				FROM 
					trips
				WHERE 
					request_time >= DATEADD(DAY, -7, CURRENT_DATE())
					{{ city_filter }}
				GROUP BY 
					daily_interval
			)
			SELECT 
				TO_CHAR(c.daily_interval, 'YYYY-MM-DD') AS daily_interval,
				c.trip_count,
				COALESCE(
					ROUND(
						(c.trip_count - LAG(c.trip_count) OVER (ORDER BY c.daily_interval)) / 
						NULLIF(LAG(c.trip_count) OVER (ORDER BY c.daily_interval), 0) * 100,
						2
					),
					0
				) AS percent_change
			FROM 
				daily_counts c
			ORDER BY 
				c.daily_interval;
		`
	} else {
		query = `
			WITH daily_counts AS (
				SELECT 
					DATE(request_time) AS daily_interval,
					COUNT(*) AS trip_count
				FROM 
					trips
				WHERE 
					request_time >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)
					{{ city_filter }}
				GROUP BY 
					daily_interval
			)
			SELECT 
				DATE_FORMAT(c.daily_interval, '%Y-%m-%d') AS daily_interval,
				c.trip_count,
				COALESCE(
					ROUND(
						(c.trip_count - LAG(c.trip_count) OVER (ORDER BY c.daily_interval)) / 
						NULLIF(LAG(c.trip_count) OVER (ORDER BY c.daily_interval), 0) * 100,
						2
					),
					0
				) AS percent_change
			FROM 
				daily_counts c
			ORDER BY 
				c.daily_interval;
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
			var dailyInterval string
			var tripCount int
			var percentChange float64

			if err := rows.Scan(&dailyInterval, &tripCount, &percentChange); err != nil {
				return nil
			}

			result := map[string]interface{}{
				"interval":       dailyInterval,
				"trip_count":     tripCount,
				"percent_change": percentChange,
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
