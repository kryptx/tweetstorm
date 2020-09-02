package json

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// Writer is a type that can write JSON
type Writer struct{}

// WriteJSON sends the application/json header and writes the body as JSON
func (writer *Writer) Write(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Add("Content-type", "application/json")
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(body)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(statusCode)
	w.Write(buf.Bytes())
}
