package server

import (
	"encoding/json"
	"net/http"
	"strconv"

	"calendar/internal/calendar"
)

type Handler struct {
	calendar *calendar.Calendar
}

func NewHandler(calendar *calendar.Calendar) *Handler {
	return &Handler{calendar: calendar}
}

type EventRequest struct {
	EventID int    `json:"event_id"`
	UserID  int    `json:"user_id"`
	Date    string `json:"date"`
	Text    string `json:"text"`
}

func (h *Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req EventRequest
	if err := decodeJSONRequest(r, &req); err != nil {
		sendError(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	if req.UserID == 0 || req.Date == "" || req.Text == "" {
		sendError(w, "Обязательные поля отсутствуют", http.StatusBadRequest)
		return
	}

	event, err := h.calendar.CreateEvent(req.UserID, req.Date, req.Text)
	if err != nil {
		handleCalendarError(w, err)
		return
	}

	sendSuccess(w, map[string]interface{}{"event": event})
}

func (h *Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req EventRequest
	if err := decodeJSONRequest(r, &req); err != nil {
		sendError(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	if req.EventID == 0 || req.UserID == 0 || req.Date == "" || req.Text == "" {
		sendError(w, "Обязательные поля отсутствуют", http.StatusBadRequest)
		return
	}

	event, err := h.calendar.UpdateEvent(req.EventID, req.UserID, req.Date, req.Text)
	if err != nil {
		handleCalendarError(w, err)
		return
	}

	sendSuccess(w, map[string]interface{}{"event": event})
}

func (h *Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		EventID int `json:"event_id"`
		UserID  int `json:"user_id"`
	}

	if err := decodeJSONRequest(r, &req); err != nil {
		sendError(w, "Неверное тело запроса", http.StatusBadRequest)
		return
	}

	if req.EventID == 0 || req.UserID == 0 {
		sendError(w, "Обязательные поля отсутствуют", http.StatusBadRequest)
		return
	}

	err := h.calendar.DeleteEvent(req.EventID, req.UserID)
	if err != nil {
		handleCalendarError(w, err)
		return
	}

	sendSuccess(w, map[string]interface{}{"message": "Событие удалено"})
}

func (h *Handler) EventsForDay(w http.ResponseWriter, r *http.Request) {
	h.handleEventsQuery(w, r, h.calendar.GetEventsForDay)
}

func (h *Handler) EventsForWeek(w http.ResponseWriter, r *http.Request) {
	h.handleEventsQuery(w, r, h.calendar.GetEventsForWeek)
}

func (h *Handler) EventsForMonth(w http.ResponseWriter, r *http.Request) {
	h.handleEventsQuery(w, r, h.calendar.GetEventsForMonth)
}

func (h *Handler) handleEventsQuery(w http.ResponseWriter, r *http.Request, getEvents func(int, string) ([]calendar.Event, error)) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || userID == 0 {
		sendError(w, "Неверный user_id", http.StatusBadRequest)
		return
	}

	date := r.URL.Query().Get("date")
	if date == "" {
		sendError(w, "Дата обязательна", http.StatusBadRequest)
		return
	}

	events, err := getEvents(userID, date)
	if err != nil {
		handleCalendarError(w, err)
		return
	}

	sendSuccess(w, map[string]interface{}{"events": events})
}

func decodeJSONRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func sendSuccess(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"result": data})
}

func sendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func handleCalendarError(w http.ResponseWriter, err error) {
	switch err {
	case calendar.ErrEventNotFound:
		sendError(w, err.Error(), http.StatusServiceUnavailable)
	case calendar.ErrInvalidDate:
		sendError(w, err.Error(), http.StatusBadRequest)
	default:
		sendError(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}
