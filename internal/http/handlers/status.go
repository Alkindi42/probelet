package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Alkindi42/probelet/internal/http/response"
)

func parseStatusCode(codeStr string) (int, bool) {
	code, err := strconv.Atoi(codeStr)

	if err != nil || code < 100 || code > 599 {
		return 0, false
	}

	return code, true
}

// NewStatusHandler returns an HTTP handler that responds with the given
// HTTP status code.
func NewStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		codeStr := r.PathValue("code")

		httpStatusCode, ok := parseStatusCode(codeStr)
		if !ok {
			response.JSONError(w, http.StatusBadRequest, "invalid status code")
			return
		}

		response.JSON(
			w,
			httpStatusCode,
			strings.ToLower(http.StatusText(httpStatusCode)),
			map[string]int{"code": httpStatusCode},
		)
	})
}
