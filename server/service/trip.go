package service

import (
	"fmt"
	"server/database"
)

func GetCurrentTripStatus(db string) map[string]interface{} {
	var result map[string]interface{}
	if db == "snowflake" {
		return result
	} else {
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
			return nil
		}

		result := make(map[string]interface{})
		for _, r := range results {
			key := fmt.Sprintf("%s_%s", r.Entity, r.Status)
			result[key] = r.Count
		}

		return result
	}
}
