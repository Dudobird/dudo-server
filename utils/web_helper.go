package utils

import (
	"encoding/json"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

//ContextToken for context key type
type ContextToken string

// Message response json data
type Message struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewMessage create a map data for send to Response function
func NewMessage(status int, message string) *Message {
	return &Message{
		Status:  status,
		Message: message,
	}
}

// JSONMessageWithData send back status and data
func JSONMessageWithData(w http.ResponseWriter, status int, text string, data interface{}) {
	message := NewMessage(status, text)
	message.Data = data
	JSONResonseWithMessage(w, message)
}

// JSONRespnseWithErr response with custom error
func JSONRespnseWithErr(w http.ResponseWriter, err *CustomError) {
	JSONRespnseWithTextMessage(w, err.Code(), err.Error())
}

// JSONRespnseWithTextMessage will send back with status and a simple text message
func JSONRespnseWithTextMessage(w http.ResponseWriter, status int, text string) {
	message := NewMessage(status, text)
	JSONResonseWithMessage(w, message)
}

// JSONResonseWithMessage will create a map data and send back to user
func JSONResonseWithMessage(w http.ResponseWriter, message *Message) {
	data, err := json.Marshal(message)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(message.Status)
	w.Write(data)
}

// ValidateUUID validate a uuid string
// return true when is valid, or false
func ValidateUUID(id string) bool {
	_, err := uuid.FromString(id)
	if err != nil {
		return false
	}
	return true
}
