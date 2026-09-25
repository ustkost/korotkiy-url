package handler

import (
	"fmt"
	"net/http"
	"strconv"
)

const defaultLimit = 20

func parseIntParam(r *http.Request, key string, fallback int) (int, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return fallback, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid integer", key)
	}
	return n, nil
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
