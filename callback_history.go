package callbackclient

import (
	"context"
	"net/http"
	"strings"

	"github.com/avast/retry-go"
)

func (c *callbackClient) GetCallbackHistoryByEventID(ctx context.Context, eventID, filter string) (*CallbackHistoryList, error) {
	errorResponse := &ErrorResponse{}

	var successResponse struct {
		Ok bool `json:"ok"`
		CallbackHistoryList
	}

	err := retry.Do(func() error {
		if err := DoRequest(
			ctx,
			http.MethodGet,
			c.URL+"/v1/callback_history/"+eventID+"?"+filter,
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

	return &successResponse.CallbackHistoryList, nil
}
