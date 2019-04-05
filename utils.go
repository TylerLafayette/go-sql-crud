package sqlcrud

import (
	"encoding/json"
	"net/http"
)

// httpWrite writes a json response and sets the code in the http header.
func httpWrite(w http.ResponseWriter, r Response) {
	// Set Content-Type header to application/json.
	w.Header().Set("Content-Type", "application/json")
	// Write the Response HTTP code into the header.
	w.WriteHeader(r.Code)
	// Create a new JSON encoder based on the ResponseWriter.
	j := json.NewEncoder(w)
	// Encode/send the Response data as JSON.
	j.Encode(r)
}
