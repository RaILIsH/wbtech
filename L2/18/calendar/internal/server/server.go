package server

import (
	"fmt"
	"net/http"

	"calendar/internal/calendar"
)

type Server struct {
	port int
	cal  *calendar.Calendar
}

func NewServer(port int, cal *calendar.Calendar) *Server {
	return &Server{port: port, cal: cal}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	handler := NewHandler(s.cal)

	mux.HandleFunc("/create_event", handler.CreateEvent)
	mux.HandleFunc("/update_event", handler.UpdateEvent)
	mux.HandleFunc("/delete_event", handler.DeleteEvent)
	mux.HandleFunc("/events_for_day", handler.EventsForDay)
	mux.HandleFunc("/events_for_week", handler.EventsForWeek)
	mux.HandleFunc("/events_for_month", handler.EventsForMonth)

	fmt.Printf("Сервер запущен на порту %d\n", s.port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.port), LoggingMiddleware(mux))
}
