package callbackclient

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/google/uuid"
)

type CallbackRequestEvent struct {
	// Unique identifier for the event
	ServiceID uuid.UUID `json:"-"` // Excluded from JSON binding
	// Payload holds the event-specific data as a key-value map
	Payload map[string]interface{} `json:"payload,omitempty"`
	// CallbackURL is the endpoint where the event data will be sent
	CallbackURL string `json:"callback_url,omitempty" example:"https://service.com/callback"`
	// WebhookSecret is a security token used to verify the event source
	WebhookSecret string `json:"webhook_secret,omitempty"`
	// Method defines the HTTP request method (e.g., POST, PUT) used to send the callback
	Method string `json:"method,omitempty" example:"POST"`
	// MaxRetries specifies the maximum number of retry attempts if the callback fails
	MaxRetries int64 `json:"max_retries,omitempty" example:"50"`
}

func (c CallbackRequestEvent) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ServiceID,
			validation.Required.Error("service_id is required"),
			validation.By(func(value interface{}) error {
				id, ok := value.(uuid.UUID)
				if !ok || id == uuid.Nil {
					return fmt.Errorf("invalid service_id")
				}
				return nil
			}),
		),
		validation.Field(&c.Payload,
			validation.Required.Error("payload is required"),
			validation.By(func(value interface{}) error {
				// check if the value is empty
				if v, ok := value.(map[string]interface{}); ok {
					if len(v) == 0 {
						return fmt.Errorf("payload must not be empty")
					}
				}
				return nil
			}),
		),
		validation.Field(&c.CallbackURL,
			validation.Required.Error("callback url is required"),
			is.URL.Error("invalid callback url provided"),
		),
		validation.Field(&c.WebhookSecret, validation.Required.Error("webhook secret is required")),
	)
}

type CallbackServiceEventConfirmation struct {
	AcknowledgementID uuid.UUID `json:"acknowledgement_id,omitempty"`
}

type Event struct {
	// Unique identifier for the event
	ID uuid.UUID `json:"id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001" format:"uuid"`
	// Unique identifier for the service
	ServiceID uuid.UUID `json:"service_id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001" format:"uuid"`
	// Service contains details of the service that triggered the event
	Service Service `json:"service,omitempty"`
	// Payload holds the event-specific data as a key-value map
	Payload map[string]interface{} `json:"payload,omitempty"`
	// CallbackURL is the endpoint where the event data will be sent
	CallbackURL string `json:"callback_url,omitempty" example:"https://service.com/callback"`
	// WebhookSecret is a security token used to verify the event source
	WebhookSecret string `json:"webhook_secret,omitempty"`
	// Method defines the HTTP request method (e.g., POST, PUT) used to send the callback
	Method Method `json:"method,omitempty" example:"POST"`
	// Status indicates the current state of the event (e.g., ACTIVE, FAILED)
	Status Status `json:"status,omitempty" example:"ACTIVE"`
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
	// UpdatedAt when the event was last updated
	UpdatedAt time.Time `json:"updated_at,omitempty" example:"2023-09-11T14:30:00Z" format:"date-time"`
}

type Service struct {
	// Unique identifier for the service
	ID uuid.UUID `json:"id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
	// Name of the service
	Name string `json:"name,omitempty" example:"Example Service"`
	// Status of the service
	Status Status `json:"status,omitempty" example:"ACTIVE"`
	// Secret token (omitted in Swagger)
	SecretToken string `json:"secret_token,omitempty" swaggerignore:"true"`
	// CreatedAt when the service was created
	CreatedAt time.Time `json:"created_at,omitempty" example:"2023-09-11T14:30:00Z"`
	// UpdatedAt when the service was last updated
	UpdatedAt time.Time `json:"updated_at,omitempty" example:"2023-09-11T14:30:00Z"`
}

type EventList struct {
	Data     []Event  `json:"data"`
	MetaData MetaData `json:"meta_data,omitempty"`
}


