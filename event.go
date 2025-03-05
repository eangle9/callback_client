package callbackclient

import (
	"context"
	"net/http"
	"strings"

	"github.com/avast/retry-go"
)

func (c *callbackClient) SendCallbackEvent(ctx context.Context, param CallbackRequestEvent) (*CallbackServiceEventConfirmation, error) {
	errorResponse := &ErrorResponse{}

	var successResponse struct {
		OK   bool                              `json:"ok"`
		Data *CallbackServiceEventConfirmation `json:"data,omitempty"`
	}

	err := retry.Do(func() error {
		if err := DoRequest(
			ctx,
			http.MethodPost,
			c.URL+"/v1/send_callback",
			func(r *http.Request) {
				r.Header.Set("Authorization", c.SecretKey)
			},
			param,
			&successResponse,
			errorResponse,
		); err != nil {
			if strings.Contains(err.Error(), "request failed") {
				return errorResponse
			}

			return err
		}
		return nil
	}, c.RetryOptions...)
	if err != nil {
		if errorResponse.CallbackError.Message != "" {
			return nil, errorResponse
		}

		return nil, err
	}

	return successResponse.Data, nil
}

func (c *callbackClient) GetEventDetailByID(ctx context.Context, eventID string) (*Event, error) {
	errorResponse := &ErrorResponse{}

	var successResponse struct {
		OK   bool   `json:"ok"`
		Data *Event `json:"data,omitempty"`
	}

	err := retry.Do(func() error {
		if err := DoRequest(
			ctx,
			http.MethodPost,
			c.URL+"/v1/event/"+eventID,
			func(r *http.Request) {
				r.Header.Set("Authorization", c.SecretKey)
			},
			nil,
			&successResponse,
			errorResponse,
		); err != nil {
			if strings.Contains(err.Error(), "request failed") {
				return errorResponse
			}

			return err
		}
		return nil
	}, c.RetryOptions...)
	if err != nil {
		if errorResponse.CallbackError.Message != "" {
			return nil, errorResponse
		}

		return nil, err
	}

	return successResponse.Data, nil
}

func (c *callbackClient) GetListOfEvents(ctx context.Context, filter string) (*EventList, error) {
	errorResponse := &ErrorResponse{}

	var successResponse struct {
		Ok bool `json:"ok"`
		EventList
	}

	err := retry.Do(func() error {
		if err := DoRequest(
			ctx,
			http.MethodGet,
			c.URL+"/v1/events?"+filter,
			func(r *http.Request) {
				r.Header.Set("Authorization", c.SecretKey)
			},
			nil,
			&successResponse,
			errorResponse,
		); err != nil {
			if strings.Contains(err.Error(), "request failed") {
				return errorResponse
			}

			return err
		}
		return nil
	}, c.RetryOptions...)
	if err != nil {
		if errorResponse.CallbackError.Message != "" {
			return nil, errorResponse
		}

		return nil, err
	}

	return &successResponse.EventList, nil
}
