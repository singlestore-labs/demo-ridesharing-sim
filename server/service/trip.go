package service

import (
	"fmt"
	"server/database"
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
