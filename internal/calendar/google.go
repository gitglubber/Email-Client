package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/gitglubber/email-client/internal/models"
	"google.golang.org/api/calendar/v3"
)

type GoogleCalendarService struct {
	service *calendar.Service
}

func NewGoogleCalendarService(service *calendar.Service) *GoogleCalendarService {
	return &GoogleCalendarService{service: service}
}

func (g *GoogleCalendarService) ListEvents(ctx context.Context, startDate, endDate time.Time) ([]*models.CalendarEvent, error) {
	events, err := g.service.Events.List("primary").
		TimeMin(startDate.Format(time.RFC3339)).
		TimeMax(endDate.Format(time.RFC3339)).
		SingleEvents(true).
		OrderBy("startTime").
		Do()

	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	calendarEvents := make([]*models.CalendarEvent, 0, len(events.Items))
	for _, event := range events.Items {
		calEvent := g.convertToCalendarEvent(event)
		calendarEvents = append(calendarEvents, calEvent)
	}

	return calendarEvents, nil
}

func (g *GoogleCalendarService) CreateEvent(ctx context.Context, req *models.CalendarEventRequest) (*models.CalendarEvent, error) {
	event := &calendar.Event{
		Summary:     req.Title,
		Description: req.Description,
		Location:    req.Location,
		Start: &calendar.EventDateTime{
			DateTime: req.Start.Format(time.RFC3339),
			TimeZone: "UTC",
		},
		End: &calendar.EventDateTime{
			DateTime: req.End.Format(time.RFC3339),
			TimeZone: "UTC",
		},
	}

	// Add attendees
	if len(req.Attendees) > 0 {
		event.Attendees = make([]*calendar.EventAttendee, 0, len(req.Attendees))
		for _, email := range req.Attendees {
			event.Attendees = append(event.Attendees, &calendar.EventAttendee{
				Email: email,
			})
		}
	}

	createdEvent, err := g.service.Events.Insert("primary", event).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return g.convertToCalendarEvent(createdEvent), nil
}

func (g *GoogleCalendarService) convertToCalendarEvent(event *calendar.Event) *models.CalendarEvent {
	calEvent := &models.CalendarEvent{
		ID:          event.Id,
		Title:       event.Summary,
		Description: event.Description,
		Location:    event.Location,
	}

	if event.Start != nil {
		if event.Start.DateTime != "" {
			startTime, _ := time.Parse(time.RFC3339, event.Start.DateTime)
			calEvent.Start = startTime
		} else if event.Start.Date != "" {
			startTime, _ := time.Parse("2006-01-02", event.Start.Date)
			calEvent.Start = startTime
			calEvent.IsAllDay = true
		}
	}

	if event.End != nil {
		if event.End.DateTime != "" {
			endTime, _ := time.Parse(time.RFC3339, event.End.DateTime)
			calEvent.End = endTime
		} else if event.End.Date != "" {
			endTime, _ := time.Parse("2006-01-02", event.End.Date)
			calEvent.End = endTime
		}
	}

	if event.Organizer != nil {
		calEvent.Organizer = event.Organizer.Email
	}

	calEvent.Attendees = make([]string, 0, len(event.Attendees))
	for _, attendee := range event.Attendees {
		calEvent.Attendees = append(calEvent.Attendees, attendee.Email)
	}

	return calEvent
}
