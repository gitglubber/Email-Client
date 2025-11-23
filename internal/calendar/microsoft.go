package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/gitglubber/email-client/internal/models"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
)

type MicrosoftCalendarService struct {
	client *msgraphsdk.GraphServiceClient
}

func NewMicrosoftCalendarService(client *msgraphsdk.GraphServiceClient) *MicrosoftCalendarService {
	return &MicrosoftCalendarService{client: client}
}

func (m *MicrosoftCalendarService) ListEvents(ctx context.Context, startDate, endDate time.Time) ([]*models.CalendarEvent, error) {
	startDateStr := startDate.Format(time.RFC3339)
	endDateStr := endDate.Format(time.RFC3339)

	filter := fmt.Sprintf("start/dateTime ge '%s' and end/dateTime le '%s'", startDateStr, endDateStr)

	config := &users.ItemCalendarEventsRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.ItemCalendarEventsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	events, err := m.client.Me().Calendar().Events().Get(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}

	calendarEvents := make([]*models.CalendarEvent, 0)
	for _, event := range events.GetValue() {
		calendarEvent := m.convertToCalendarEvent(event)
		calendarEvents = append(calendarEvents, calendarEvent)
	}

	return calendarEvents, nil
}

func (m *MicrosoftCalendarService) CreateEvent(ctx context.Context, req *models.CalendarEventRequest) (*models.CalendarEvent, error) {
	event := graphmodels.NewEvent()

	subject := req.Title
	event.SetSubject(&subject)

	if req.Description != "" {
		body := graphmodels.NewItemBody()
		content := req.Description
		body.SetContent(&content)
		contentType := graphmodels.TEXT_BODYTYPE
		body.SetContentType(&contentType)
		event.SetBody(body)
	}

	// Set start time
	start := graphmodels.NewDateTimeTimeZone()
	startStr := req.Start.Format(time.RFC3339)
	start.SetDateTime(&startStr)
	timezone := "UTC"
	start.SetTimeZone(&timezone)
	event.SetStart(start)

	// Set end time
	end := graphmodels.NewDateTimeTimeZone()
	endStr := req.End.Format(time.RFC3339)
	end.SetDateTime(&endStr)
	end.SetTimeZone(&timezone)
	event.SetEnd(end)

	if req.Location != "" {
		location := graphmodels.NewLocation()
		displayName := req.Location
		location.SetDisplayName(&displayName)
		event.SetLocation(location)
	}

	// Set attendees
	if len(req.Attendees) > 0 {
		attendees := make([]graphmodels.Attendeeable, 0, len(req.Attendees))
		for _, email := range req.Attendees {
			attendee := graphmodels.NewAttendee()
			emailAddr := graphmodels.NewEmailAddress()
			address := email
			emailAddr.SetAddress(&address)
			attendee.SetEmailAddress(emailAddr)
			attendeeType := graphmodels.REQUIRED_ATTENDEETYPE
			attendee.SetTypeEscaped(&attendeeType)
			attendees = append(attendees, attendee)
		}
		event.SetAttendees(attendees)
	}

	createdEvent, err := m.client.Me().Calendar().Events().Post(ctx, event, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return m.convertToCalendarEvent(createdEvent), nil
}

func (m *MicrosoftCalendarService) convertToCalendarEvent(event graphmodels.Eventable) *models.CalendarEvent {
	calEvent := &models.CalendarEvent{
		ID:    *event.GetId(),
		Title: *event.GetSubject(),
	}

	if event.GetBody() != nil && event.GetBody().GetContent() != nil {
		calEvent.Description = *event.GetBody().GetContent()
	}

	if event.GetStart() != nil && event.GetStart().GetDateTime() != nil {
		startTime, _ := time.Parse(time.RFC3339, *event.GetStart().GetDateTime())
		calEvent.Start = startTime
	}

	if event.GetEnd() != nil && event.GetEnd().GetDateTime() != nil {
		endTime, _ := time.Parse(time.RFC3339, *event.GetEnd().GetDateTime())
		calEvent.End = endTime
	}

	if event.GetLocation() != nil && event.GetLocation().GetDisplayName() != nil {
		calEvent.Location = *event.GetLocation().GetDisplayName()
	}

	if event.GetOrganizer() != nil && event.GetOrganizer().GetEmailAddress() != nil {
		calEvent.Organizer = *event.GetOrganizer().GetEmailAddress().GetAddress()
	}

	attendees := event.GetAttendees()
	calEvent.Attendees = make([]string, 0, len(attendees))
	for _, attendee := range attendees {
		if attendee.GetEmailAddress() != nil {
			calEvent.Attendees = append(calEvent.Attendees, *attendee.GetEmailAddress().GetAddress())
		}
	}

	if event.GetIsAllDay() != nil {
		calEvent.IsAllDay = *event.GetIsAllDay()
	}

	return calEvent
}
