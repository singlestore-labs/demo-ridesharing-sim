package database

import (
	"simulator/models"

	cmap "github.com/orcaman/concurrent-map/v2"
)

var Local *LocalStore

type LocalStore struct {
	Riders  cmap.ConcurrentMap[string, models.Rider]
	Drivers cmap.ConcurrentMap[string, models.Driver]
	Trips   cmap.ConcurrentMap[string, models.Trip]
}

func InitializeLocal() {
	Local = &LocalStore{
		Riders:  cmap.New[models.Rider](),
		Drivers: cmap.New[models.Driver](),
		Trips:   cmap.New[models.Trip](),
	}
}
