package handlers

import (
	"net/http"
	"strconv"
)

func NewStatusHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		codeStr := r.PathValue("code")

		httpStatusCode, err := strconv.Atoi(codeStr)

		if err != nil || httpStatusCode < 100 || httpStatusCode > 599 {
			http.Error(w, "Invalid status code", http.StatusBadRequest)
			return
		}

		w.WriteHeader(httpStatusCode)
		w.Write([]byte(http.StatusText(httpStatusCode)))
	})
}
