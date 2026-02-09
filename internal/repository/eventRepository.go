package repository

import (
	"context"
	"dailyPlanner/internal/database"
	"dailyPlanner/internal/models"
	"fmt"
	"github.com/google/uuid"
)

type eventRepository struct {
	db *database.DB
}

func NewEventRepository(db *database.DB) *eventRepository {
	return &eventRepository{db}
}

func (vr eventRepository) CreateEvent(ctx context.Context, userId string, event *models.Event) error {
	event.EventId = uuid.New().String()
	event.UserId = userId
	event.Completed = false
	if event.Color == "" {
		event.Color = "red"
	}

	query := `INSERT INTO events (event_id, user_id, title_event, date_event, completed, color)
			   VALUES (:event_id, :user_id, :title_event, :date_event, :completed, :color)`

	_, err := vr.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("error when creating a event when accessing the database: %w", err)
	}

	return nil
}

func (vr eventRepository) GetEventById(ctx context.Context, eventId string) (*models.Event, error) {
	query := `SELECT * FROM events WHERE event_id = $1`

	var event models.Event
	err := vr.db.GetContext(ctx, &event, query, eventId)
	if err != nil {
		return nil, fmt.Errorf("error when getting the event by id: %w", err)
	}
	return &event, nil
}

func (vr eventRepository) GetEventsByUserIdAndDate(ctx context.Context, userId, date string) ([]*models.Event, error) {
	query := `SELECT * FROM events WHERE user_id = $1 AND date_event = $2`

	var events []*models.Event
	err := vr.db.SelectContext(ctx, &events, query, userId, date)
	if err != nil {
		return nil, fmt.Errorf("error when getting the events by user id: %w", err)
	}
	return events, nil
}

func (vr eventRepository) UpdateEvent(ctx context.Context, eventId, newTitle, color string) error {
	event, err := vr.GetEventById(ctx, eventId)
	if err != nil {
		return fmt.Errorf("error when getting the event by id in update event: %w", err)
	}

	event.TitleEvent = newTitle
	if color != "" {
		event.Color = color
	}

	query := `UPDATE events
			SET title_event = :title_event, color = :color
			WHERE event_id = :event_id`

	result, err := vr.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("error when updating the event: %w", err)
	}

	if !CheckUpdate(result) {
		return fmt.Errorf("error when updating the event in check update: %w", err)
	}

	return nil
}

func (vr eventRepository) CompleteEvent(ctx context.Context, eventId string) error {
	event, err := vr.GetEventById(ctx, eventId)
	if err != nil {
		return fmt.Errorf("error when getting the event by id: %w", err)
	}

	event.Completed = true

	query := `UPDATE events
			SET completed = :completed 
			WHERE event_id = :event_id`

	result, err := vr.db.NamedExecContext(ctx, query, event)
	if err != nil {
		return fmt.Errorf("error when updating the event for completed: %w", err)
	}
	if !CheckUpdate(result) {
		return fmt.Errorf("error when updating the event for completed: %w", err)
	}
	return nil
}

func (vr eventRepository) DeleteEvent(ctx context.Context, eventId string) error {
	query := `DELETE FROM events WHERE event_id = $1`

	result, err := vr.db.ExecContext(ctx, query, eventId)
	if err != nil {
		return fmt.Errorf("error when deleting the event: %w", err)
	}
	if !CheckUpdate(result) {
		return fmt.Errorf("error when deleting the event in check update: %w", err)
	}
	return nil
}
