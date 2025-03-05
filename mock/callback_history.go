package mock

import (
	"context"
	"fmt"

	callback "dev.azure.com/2f-capital/go-packages/callback-client.git"
	"github.com/google/uuid"
)

func (c *callbackClient) GetCallbackHistoryByEventID(ctx context.Context, eventID, filter string) (*callback.CallbackHistoryList, error) {
	for _, e := range c.Service.Events {
		if e.ID == eventID {
			var callbackHistory []callback.CallbackHistory
			if len(e.CallbackHistory) == 0 {
				return &callback.CallbackHistoryList{Data: []callback.CallbackHistory{}}, nil
			}

			event, err := c.GetEventDetailByID(ctx, eventID)
			if err != nil {
				return nil, err
			}
			for _, ch := range e.CallbackHistory {
				callbackHistory = append(callbackHistory, callback.CallbackHistory{
					ID:           uuid.MustParse(ch.ID),
					Event:        *event,
					Status:       callback.Status(ch.Status),
					ResponseCode: ch.ResponseCode,
					ReasonFailed: ch.ReasonFailed,
					CreatedAt:    ch.CreatedAt,
				})
			}

			return &callback.CallbackHistoryList{Data: callbackHistory}, nil
		}
	}

	return nil, fmt.Errorf("callback history with eventID %s not found", eventID)
}
