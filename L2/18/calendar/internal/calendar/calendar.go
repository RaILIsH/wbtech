package calendar

import (
	"errors"
	"sort"
	"sync"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidDate   = errors.New("invalid date format")
)

type Event struct {
	ID     int       `json:"id"`
	UserID int       `json:"user_id"`
	Date   time.Time `json:"date"`
	Text   string    `json:"text"`
}

type Calendar struct {
	mu      sync.RWMutex
	events  []Event
	counter int
}

func NewCalendar() *Calendar {
	return &Calendar{events: make([]Event, 0)}
}

func (c *Calendar) CreateEvent(userID int, dateStr, text string) (Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Event{}, ErrInvalidDate
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	event := Event{
		ID:     c.counter + 1,
		UserID: userID,
		Date:   date,
		Text:   text,
	}

	c.events = append(c.events, event)
	c.counter++

	return event, nil
}

func (c *Calendar) UpdateEvent(eventID, userID int, dateStr, text string) (Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return Event{}, ErrInvalidDate
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for i, event := range c.events {
		if event.ID == eventID && event.UserID == userID {
			updatedEvent := Event{ID: eventID, UserID: userID, Date: date, Text: text}
			c.events[i] = updatedEvent
			return updatedEvent, nil
		}
	}

	return Event{}, ErrEventNotFound
}

func (c *Calendar) DeleteEvent(eventID, userID int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i, event := range c.events {
		if event.ID == eventID && event.UserID == userID {
			c.events = append(c.events[:i], c.events[i+1:]...)
			return nil
		}
	}

	return ErrEventNotFound
}

func (c *Calendar) GetEventsForDay(userID int, dateStr string) ([]Event, error) {
	return c.getEventsByFilter(userID, dateStr, func(eventDate, targetDate time.Time) bool {
		return isSameDay(eventDate, targetDate)
	})
}

func (c *Calendar) GetEventsForWeek(userID int, dateStr string) ([]Event, error) {
	return c.getEventsByFilter(userID, dateStr, func(eventDate, targetDate time.Time) bool {
		ey, ew := eventDate.ISOWeek()
		ty, tw := targetDate.ISOWeek()
		return ey == ty && ew == tw
	})
}

func (c *Calendar) GetEventsForMonth(userID int, dateStr string) ([]Event, error) {
	return c.getEventsByFilter(userID, dateStr, func(eventDate, targetDate time.Time) bool {
		return eventDate.Year() == targetDate.Year() && eventDate.Month() == targetDate.Month()
	})
}

func (c *Calendar) getEventsByFilter(userID int, dateStr string, filter func(time.Time, time.Time) bool) ([]Event, error) {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, ErrInvalidDate
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	var result []Event
	for _, event := range c.events {
		if event.UserID == userID && filter(event.Date, date) {
			result = append(result, event)
		}
	}

	sortEvents(result)
	return result, nil
}

func isSameDay(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func sortEvents(events []Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Date.Before(events[j].Date)
	})
}
