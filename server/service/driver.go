package service

import (
	"server/database"
	"server/model"
)

func GetAllDrivers(db string) []model.Driver {
	var drivers []model.Driver
	if db == "snowflake" {
		rows, err := database.SnowflakeDB.Query("SELECT * FROM drivers")
		if err != nil {
			return nil
		}
		defer rows.Close()
		for rows.Next() {
			var driver model.Driver
			err := rows.Scan(
				&driver.ID,
				&driver.FirstName,
				&driver.LastName,
				&driver.Email,
				&driver.PhoneNumber,
				&driver.DateOfBirth,
				&driver.CreatedAt,
				&driver.LocationCity,
				&driver.LocationLat,
				&driver.LocationLong,
				&driver.Status,
			)
			if err != nil {
				continue
			}
			drivers = append(drivers, driver)
		}
		if err = rows.Err(); err != nil {
			return nil
		}
	} else {
		database.SingleStoreDB.Find(&drivers)
		return drivers
	}
	return drivers
}
