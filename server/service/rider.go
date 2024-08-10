package service

import (
	"server/database"
	"server/model"
)

func GetAllRiders(store string) []model.Rider {
	var riders []model.Rider
	if store == "snowflake" {
		rows, err := database.SnowflakeDB.Query("SELECT * FROM riders")
		if err != nil {
			return nil
		}
		defer rows.Close()
		for rows.Next() {
			var rider model.Rider
			err := rows.Scan(
				&rider.ID,
				&rider.FirstName,
				&rider.LastName,
				&rider.Email,
				&rider.PhoneNumber,
				&rider.DateOfBirth,
				&rider.CreatedAt,
				&rider.LocationCity,
				&rider.LocationLat,
				&rider.LocationLong,
				&rider.Status,
			)
			if err != nil {
				continue
			}
			riders = append(riders, rider)
		}
		if err = rows.Err(); err != nil {
			return nil
		}
	} else {
		database.SingleStoreDB.Find(&riders)
		return riders
	}
	return riders
}
