package service

import (
	"fmt"
	"server/database"
)

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
		return result
	}

	query := `
		SELECT entity, status, COUNT(*) as count
		FROM (
			SELECT 'trips' as entity, status FROM trips
			UNION ALL
			SELECT 'riders' as entity, status FROM riders
			UNION ALL
			SELECT 'drivers' as entity, status FROM drivers
		) AS combined
		GROUP BY entity, status
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
