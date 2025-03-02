package error_handler

import (
	"encoding/json"
	"net/http"
	"store-manager/internal/application/DTOs"
)

func WriteJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(DTOs.ErrorResponse{
		Code:    status,
		Message: message,
	})
}
