package mock

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	callback "dev.azure.com/2f-capital/go-packages/callback-client.git"
	"github.com/google/uuid"
)

var secretKey = "test webhook secret key"

func InitTestCallbackServer() *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/v1/callback":
			verifyCallbacks(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	return server
}

func verifyCallbacks(w http.ResponseWriter, r *http.Request) {
	hash := r.Header.Get("X-MP-SIGNATURE")
	t := r.Header.Get("X-MP-Time")

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mac := hmac.New(sha256.New, []byte(secretKey))

	if _, err := mac.Write([]byte(t)); err != nil {
		return
	}

	if _, err := mac.Write([]byte(".")); err != nil {
		return
	}

	if _, err := mac.Write([]byte(payload)); err != nil {
		return
	}

	expectedHash := hex.EncodeToString(mac.Sum(nil))

	if expectedHash != hash {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func TestSendCallbackEvent(t *testing.T) {
	server := InitTestCallbackServer()
	defer server.Close()
	cb := callbackClient{
		Service: Service{
			Status: callback.StatusActive,
			Events: make(map[string]*Event),
		},
	}

	type args struct {
		arg callback.CallbackRequestEvent
	}

	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "successfully send callback",
			args: args{
				arg: callback.CallbackRequestEvent{
					Payload: map[string]interface{}{
						"event":          "payment_success",
						"transaction_id": "txn_123456",
						"amount":         100.50,
						"currency":       "USD",
						"status":         "completed",
					},
					CallbackURL:   server.URL + "/v1/callback",
					WebhookSecret: secretKey,
					Method:        http.MethodPost,
				},
			},
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			confirmation, err := cb.SendCallbackEvent(context.Background(), tt.args.arg)
			if (err != nil) == !tt.wantErr {
				t.Errorf("expected to get nil error, but got %v", err)
				return
			}

			got := confirmation.AcknowledgementID != uuid.Nil
			if got != tt.want {
				t.Errorf("expected to get %v, but got %v", tt.want, got)
				return
			}
		})
	}
}

func TestGetEventDetailByID(t *testing.T) {
	cb := callbackClient{
		Service: Service{
			Status: callback.StatusActive,
			Events: make(map[string]*Event),
		},
	}
	mockEvent := &Event{
		ID: uuid.NewString(),
		Payload: map[string]interface{}{
			"test": "test callback payload",
		},
		CallbackURL:      "http://test_callback_url",
		WebhookSecret:    "test secret key",
		Method:           http.MethodPost,
		Status:           callback.StatusSucceeded,
		RetryCount:       4,
		LastResponseCode: http.StatusOK,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		CallbackHistory:  make(map[string]*CallbackHistory),
	}
	cb.Service.Events[mockEvent.ID] = mockEvent

	type args struct {
		eventID string
	}

	tests := []struct {
		name    string
		args    args
		want    callback.Event
		wantErr bool
	}{
		{
			name: "succefully get event detail",
			args: args{
				eventID: mockEvent.ID,
			},
			want: callback.Event{
				ID:               uuid.MustParse(mockEvent.ID),
				Payload:          mockEvent.Payload,
				CallbackURL:      mockEvent.CallbackURL,
				WebhookSecret:    mockEvent.WebhookSecret,
				Method:           mockEvent.Method,
				Status:           mockEvent.Status,
				RetryCount:       mockEvent.RetryCount,
				LastResponseCode: mockEvent.LastResponseCode,
				CreatedAt:        mockEvent.CreatedAt,
				UpdatedAt:        mockEvent.UpdatedAt,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		got, err := cb.GetEventDetailByID(context.Background(), tt.args.eventID)
		if (err != nil) != tt.wantErr {
			t.Errorf("expected to get nil error, but got %v", err)
			return
		}

		if !reflect.DeepEqual(got.ID, tt.want.ID) || !reflect.DeepEqual(got.Payload, tt.want.Payload) || !reflect.DeepEqual(got.CallbackURL, tt.want.CallbackURL) {
			t.Errorf("expected to get %+v, but got %+v", tt.want, got)
			return
		}
	}
}
