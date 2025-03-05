package mock

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GenerateEventHash(payload []byte, sk string, t time.Time) (string, error) {
	mac := hmac.New(sha256.New, []byte(sk))

	_, err := mac.Write([]byte(fmt.Sprintf("%d", t.Unix())))
	if err != nil {
		return "", err
	}

	_, err = mac.Write([]byte("."))
	if err != nil {
		return "", err
	}

	_, err = mac.Write(payload)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(mac.Sum(nil)), nil
}

func DoRequest(
	ctx context.Context,
	method,
	url string,
	contentTypeAccept string,
	modifyRequest func(*http.Request),
	body interface{},
	response interface{},
) (*http.Response, error) {
	// Convert body to []byte
	var reqBody []byte
	if body != nil {
		var ok bool
		reqBody, ok = body.([]byte)
		if !ok {
			// Convert body to JSON
			jsonBody, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			reqBody = jsonBody
		}
	}
	// Create a new request object
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, err
	}
	// Set headers
	if modifyRequest != nil {
		modifyRequest(req)
	}
	if len(req.Header["Content-Type"]) < 1 {
		req.Header.Set("Content-Type", "application/json")
	}
	// Send the request and get the response
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if response == nil {
		return resp, nil
	}

	// Create a buffer to store the response content
	var responseBodyBuffer bytes.Buffer

	// Create a TeeReader to read and capture the response content
	teeReader := io.TeeReader(resp.Body, &responseBodyBuffer)

	// Read the response body
	respBody, err := io.ReadAll(teeReader)
	if err != nil {
		return nil, err
	}

	if contentTypeAccept == "text/json" ||
		contentTypeAccept == "" ||
		contentTypeAccept == "application/json" {
		err = json.Unmarshal(respBody, response)
	} else {
		err = xml.Unmarshal(respBody, response)
	}

	resp.Body = io.NopCloser(&responseBodyBuffer)

	return resp, err
}
