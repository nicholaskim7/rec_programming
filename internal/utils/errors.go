package utils

import (
	"net/http"
    "encoding/json"
)
//http.Error() sets content type to plain text (problem in the future when frontend clients crash)
// instead send structured json response
func RespondWithError(w http.ResponseWriter, code int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}