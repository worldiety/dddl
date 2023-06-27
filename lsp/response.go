package lsp

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

type Response struct {
	Id     float64     `json:"id"`
	Result interface{} `json:"result"`
}

type Notification struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// async queue calls using SendNotification will corrupt SendResponse and other SendNotification calls.
var msgWriteLock sync.Mutex

// Send a response as one to the given request.
func SendResponse(response interface{}, requestId float64) error {
	msgWriteLock.Lock()
	defer msgWriteLock.Unlock()

	responseBytes, err := json.Marshal(Response{
		Id:     requestId,
		Result: response,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	responseData := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(responseBytes), responseBytes)
	log.Printf("Sending response: %s", responseBytes)
	fmt.Print(responseData)

	return nil
}

// Send a notification, which holds some information we push to the client.
func SendNotification(method string, notification interface{}) error {
	msgWriteLock.Lock()
	defer msgWriteLock.Unlock()

	responseBytes, err := json.Marshal(Notification{
		Method: method,
		Params: notification,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	responseData := fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(responseBytes), responseBytes)
	log.Printf("Sending notification: %s", responseBytes)
	fmt.Print(responseData)

	return nil
}
