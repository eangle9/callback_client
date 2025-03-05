package callbackclient

import (
	"time"

	"github.com/google/uuid"
)

type CallbackHistory struct {
	ID           uuid.UUID `json:"id,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
	Event        Event     `json:"event,omitempty"`
	Status       Status    `json:"status,omitempty" example:"FAILED"`
	ResponseCode int64     `json:"response_code,omitempty" example:"200"`
	ReasonFailed string    `json:"reason_failed,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty" example:"2023-09-11T14:30:00Z"`
	UpdatedAt    time.Time `json:"updated_at,omitempty" example:"2023-09-11T14:30:00Z"`
}

type CallbackHistoryList struct {
	Data     []CallbackHistory `json:"data"`
	MetaData MetaData          `json:"meta_data,omitempty"`
}
