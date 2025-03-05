package mock

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	callback "dev.azure.com/2f-capital/go-packages/callback-client.git"
	"github.com/google/uuid"
)

func (c *callbackClient) SendCallbackEvent(ctx context.Context, param callback.CallbackRequestEvent) (*callback.CallbackServiceEventConfirmation, error) {
	eventData := &Event{
		ID:              uuid.NewString(),
		Payload:         param.Payload,
		CallbackURL:     param.CallbackURL,
		WebhookSecret:   param.WebhookSecret,
		Method:          callback.Method(param.Method),
		Status:          callback.StatusActive,
		MaxRetries:      param.MaxRetries,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		CallbackHistory: make(map[string]*CallbackHistory),
	}
	c.Service.Events[eventData.ID] = eventData

	ht := time.Now()
	payload, err := json.Marshal(param.Payload)
	if err != nil {
		return nil, err
	}

	hash, err := GenerateEventHash(payload, param.WebhookSecret, ht)
	if err != nil {
		return nil, err
	}

	res, err := DoRequest(
		ctx,
		param.Method,
		param.CallbackURL,
		"application/json",
		func(r *http.Request) {
			r.Header.Set("Content-Type", "application/json")
			r.Header.Set("X-MP-SIGNATURE", hash)
			r.Header.Set("X-MP-Time", fmt.Sprintf("%d", ht.Unix()))
		},
		param.Payload,
		nil,
	)
	if res != nil {
		res.Body.Close()
	}

	if err != nil {
		callbackHistory := &CallbackHistory{
			ID:           uuid.NewString(),
			Status:       string(callback.StatusFailed),
			ReasonFailed: err.Error(),
			CreatedAt:    time.Now(),
		}
		c.Service.Events[eventData.ID].CallbackHistory[callbackHistory.ID] = callbackHistory
		c.Service.Events[eventData.ID].ReasonFailed = err.Error()
		c.Service.Events[eventData.ID].Status = callback.StatusFailed
		c.Service.Events[eventData.ID].RetryCount++
		c.Service.Events[eventData.ID].UpdatedAt = time.Now()
		return nil, err
	} else if res.StatusCode != http.StatusOK {
		err := fmt.Errorf("webhook rejected by service with statuscode %d", res.StatusCode)
		callbackHistory := &CallbackHistory{
			ID:           uuid.NewString(),
			Status:       string(callback.StatusFailed),
			ResponseCode: int64(res.StatusCode),
			ReasonFailed: err.Error(),
			CreatedAt:    time.Now(),
		}
		c.Service.Events[eventData.ID].CallbackHistory[callbackHistory.ID] = callbackHistory
		c.Service.Events[eventData.ID].ReasonFailed = err.Error()
		c.Service.Events[eventData.ID].Status = callback.StatusFailed
		c.Service.Events[eventData.ID].LastResponseCode = int64(res.StatusCode)
		c.Service.Events[eventData.ID].RetryCount++
		c.Service.Events[eventData.ID].UpdatedAt = time.Now()
		return nil, err
	}

	c.Service.Events[eventData.ID].Status = callback.StatusSucceeded
	c.Service.Events[eventData.ID].LastResponseCode = int64(res.StatusCode)
	c.Service.Events[eventData.ID].RetryCount++
	c.Service.Events[eventData.ID].UpdatedAt = time.Now()

	callbackHistory := &CallbackHistory{
		ID:           uuid.NewString(),
		Status:       string(callback.StatusSucceeded),
		ResponseCode: int64(res.StatusCode),
		CreatedAt:    time.Now(),
	}
	c.Service.Events[eventData.ID].CallbackHistory[callbackHistory.ID] = callbackHistory

	eventID, err := uuid.Parse(eventData.ID)
	if err != nil {
		return nil, err
	}
	response := &callback.CallbackServiceEventConfirmation{
		AcknowledgementID: eventID,
	}
	return response, nil
}

func (c *callbackClient) GetEventDetailByID(ctx context.Context, eventID string) (*callback.Event, error) {
	for _, e := range c.Service.Events {
		if e.ID == eventID {
			event := &callback.Event{
				ID:               uuid.MustParse(e.ID),
				Payload:          e.Payload,
				CallbackURL:      e.CallbackURL,
				WebhookSecret:    e.WebhookSecret,
				Method:           e.Method,
				Status:           e.Status,
				MaxRetries:       e.MaxRetries,
				RetryCount:       e.RetryCount,
				NextRetryAt:      e.NextRetryAt,
				LastResponseCode: e.LastResponseCode,
				ReasonFailed:     e.ReasonFailed,
				CreatedAt:        e.CreatedAt,
				UpdatedAt:        e.UpdatedAt,
			}
			return event, nil
		}
	}

	return nil, fmt.Errorf("event not found")
}

func (c *callbackClient) GetListOfEvents(ctx context.Context, filter string) (*callback.EventList, error) {
	var events []callback.Event

	for _, e := range c.Service.Events {
		event := callback.Event{
			ID:               uuid.MustParse(e.ID),
			Payload:          e.Payload,
			CallbackURL:      e.CallbackURL,
			WebhookSecret:    e.WebhookSecret,
			Method:           e.Method,
			Status:           e.Status,
			MaxRetries:       e.MaxRetries,
			RetryCount:       e.RetryCount,
			NextRetryAt:      e.NextRetryAt,
			LastResponseCode: e.LastResponseCode,
			ReasonFailed:     e.ReasonFailed,
			CreatedAt:        e.CreatedAt,
			UpdatedAt:        e.UpdatedAt,
		}
		events = append(events, event)
	}
	return &callback.EventList{Data: events}, nil
}
