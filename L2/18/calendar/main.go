package main

import (
	"calendar/config"
	"calendar/internal/calendar"
	"calendar/internal/server"
)

func main() {
	cfg := config.Load()
	cal := calendar.NewCalendar()
	srv := server.NewServer(cfg.Port, cal)

	if err := srv.Start(); err != nil {
		panic(err)
	}
}
