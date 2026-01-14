package handlers

import (
	"github.com/Alkindi42/probelet/internal/http/response"
	"net/http"
	"strconv"
)

// NewStatusHandler returns an HTTP handler that responds with the given
// HTTP status code.
func NewStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		codeStr := r.PathValue("code")

		httpStatusCode, err := strconv.Atoi(codeStr)

		if err != nil || httpStatusCode < 100 || httpStatusCode > 599 {
			response.JSONError(w, http.StatusBadRequest, "Invalid status code")
			return
		}

		response.JSONResponse(
			w,
			httpStatusCode,
			http.StatusText(httpStatusCode),
			map[string]int{"code": httpStatusCode},
		)
	})
}
