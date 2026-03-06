package json

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// WriteError writes a standardized error response.
// It automatically determines the status code if the error implements StatusError.
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	var statusErr StatusError
	status := http.StatusInternalServerError

	if errors.As(err, &statusErr) {
		status = statusErr.Status()
	}

	Write(w, r, status, err.Error())
}

type Pagination struct {
	Page       int   `json:"page"`
	Size       int   `json:"size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

type Response struct {
	ResponseCode    string      `json:"responseCode"`
	ResponseMessage string      `json:"responseMessage"`
	Data            any         `json:"data"`
	Pagination      *Pagination `json:"pagination,omitempty"`
	TraceID         string      `json:"traceId"`
	Timestamp       string      `json:"timestamp"`
}

func NewPagination(page, size int, total int64) Pagination {
	totalPages := int(total) / size
	if int(total)%size > 0 {
		totalPages++
	}

	return Pagination{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: totalPages,
	}
}

func Write(w http.ResponseWriter, r *http.Request, status int, data any) {
	write(w, r, status, data, nil)
}

func WriteWithPagination(w http.ResponseWriter, r *http.Request, status int, data any, pag Pagination) {
	write(w, r, status, data, &pag)
}

func write(w http.ResponseWriter, r *http.Request, status int, data any, pag *Pagination) {
	resp := Response{
		ResponseCode:    "00",
		ResponseMessage: "Success",
		Data:            data,
		Pagination:      pag,
		TraceID:         middleware.GetReqID(r.Context()),
		Timestamp:       time.Now().Format(time.RFC3339),
	}

	if status >= 400 {
		resp.ResponseCode = "99"
		resp.ResponseMessage = "Error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
