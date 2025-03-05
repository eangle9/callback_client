package mock

import (
	"time"

	callback "dev.azure.com/2f-capital/go-packages/callback-client.git"
)

type callbackClient struct {
	Service Service
}

type Service struct {
	// ID is the unique identifier for the service.
	// It is automatically generated when the service is created.
	ID string `json:"id,omitempty"`
	// Status is the current status of the service.
	// It is set to active by default.
	Status callback.Status `json:"status,omitempty"`
	// SecretToken is the secret the service uses to authenticate itself.
	// It is automatically generated when the service is created.
	SecretToken string `json:"secret_token,omitempty"`
	Events      map[string]*Event
}

type Event struct {
	// Unique identifier for the event
	ID string `json:"id,omitempty"`
	// Payload holds the event-specific data as a key-value map
	Payload map[string]interface{} `json:"payload,omitempty"`
	// CallbackURL is the endpoint where the event data will be sent
	CallbackURL string `json:"callback_url,omitempty" example:"https://service.com/callback"`
	// WebhookSecret is a security token used to verify the event source
	WebhookSecret string `json:"webhook_secret,omitempty"`
	// Method defines the HTTP request method (e.g., POST, PUT) used to send the callback
	Method callback.Method `json:"method,omitempty" example:"POST"`
	// Status indicates the current state of the event (e.g., ACTIVE, FAILED)
	Status callback.Status `json:"status,omitempty" example:"ACTIVE"`
	// MaxRetries specifies the maximum number of retry attempts if the callback fails
	MaxRetries int64 `json:"max_retries,omitempty" example:"50"`
	// RetryCount tracks the number of times the event has been retried
	RetryCount int64 `json:"retry_count,omitempty" example:"10"`
	// NextRetryAt specifies the scheduled time for the next retry attempt
	NextRetryAt time.Time `json:"next_retry_at,omitempty" example:"2023-09-11T14:30:00Z" format:"date-time"`
	// LastResponseCode stores the HTTP response code from the last callback attempt
	LastResponseCode int64 `json:"last_response_code,omitempty" example:"200"`
	// reason stores an error for failed callback attempt
	ReasonFailed string `json:"reason_failed,omitempty"`
	// CreatedAt when the event was created
	CreatedAt time.Time `json:"created_at,omitempty" example:"2023-09-11T14:30:00Z" format:"date-time"`
	// UpdatedAt when the service was last updated
	UpdatedAt       time.Time `json:"updated_at,omitempty" example:"2023-09-11T14:30:00Z" format:"date-time"`
	CallbackHistory map[string]*CallbackHistory
}

type CallbackHistory struct {
	// id is the unique identifier for the callback history
	ID string `json:"id,omitempty"`
	// Status indicates the status of the callback attempt
	Status string `json:"status,omitempty" example:"FAILED"`
	// ResponseCode stores the HTTP response code from the callback attempt
	ResponseCode int64 `json:"response_code,omitempty" example:"200"`
	// reason stores an error for failed callback attempt
	ReasonFailed string `json:"reason_failed,omitempty"`
	// CreatedAt when the callback history was created
	CreatedAt time.Time `json:"created_at,omitempty" example:"2023-09-11T14:30:00Z" format:"date-time"`
}
