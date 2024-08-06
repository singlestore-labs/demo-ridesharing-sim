package database

import "simulator/models"

var Local *LocalStore

type LocalStore struct {
	Riders  map[string]models.Rider
	Drivers map[string]models.Driver
	Trips   map[string]models.Trip
}

func InitializeLocal() {
	Local = &LocalStore{
		Riders:  make(map[string]models.Rider, 0),
		Drivers: make(map[string]models.Driver, 0),
		Trips:   make(map[string]models.Trip, 0),
	}
}
