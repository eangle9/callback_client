package callbackclient

import (
	"context"

	"github.com/avast/retry-go"
)

type callbackClient struct {
	URL          string         // callback server url
	SecretKey    string         // service secret key
	RetryOptions []retry.Option // Retry all errors, not just the last error
}

type Client interface {
	SendCallbackEvent(ctx context.Context, param CallbackRequestEvent) (*CallbackServiceEventConfirmation, error)
	GetEventDetailByID(ctx context.Context, eventID string) (*Event, error)
	GetListOfEvents(ctx context.Context, filter string) (*EventList, error)
	GetCallbackHistoryByEventID(ctx context.Context, eventID, filter string) (*CallbackHistoryList, error)
}

func NewAccountClient(url string, secret string, retryOptions []retry.Option) Client {
	return &callbackClient{
		URL:          url,
		SecretKey:    secret,
		RetryOptions: retryOptions,
	}
}
