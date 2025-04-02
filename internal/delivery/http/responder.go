package http

import (
	"github.com/google/uuid"
	"net/http"

	json "github.com/json-iterator/go"

	"go.uber.org/zap"
)

type responder struct {
	logger *zap.Logger
}

func NewResponder(logger *zap.Logger) *responder {
	return &responder{
		logger: logger,
	}
}

type errorDetails struct {
	Fields  []string `json:"fields"`
	Message string   `json:"message"`
}

type meta struct {
	Code         int            `json:"code"`
	Message      string         `json:"message"`
	DebugId      string         `json:"debug_id"`
	Reason       string         `json:"reason,omitempty"`
	ErrorDetails []errorDetails `json:"details,omitempty"`
}

type myResponse struct {
	Meta meta `json:"meta"`
}

func (r *responder) With(statusCode int, w http.ResponseWriter, v interface{}) {
	if (statusCode >= 200 && statusCode < 300) || statusCode == 429 {
		r.writeResponse(w, v, statusCode)
	} else {
		r.withError(w, v.(string), statusCode)
	}

}
func (r *responder) WithOK(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusOK)
}
func (r *responder) WithNotFound(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusNotFound)
}

func (r *responder) WithCreated(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusCreated)
}

func (r *responder) WithBadRequest(w http.ResponseWriter, message string) string {
	return r.withError(w, message, http.StatusBadRequest)
}

func (r *responder) WithInternalError(w http.ResponseWriter, message string) string {
	return r.withError(w, message, http.StatusInternalServerError)
}

func (r *responder) WithUnauthorizedError(w http.ResponseWriter) string {
	return r.withError(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
}

func (r *responder) WithForbiddenError(w http.ResponseWriter) string {
	return r.withError(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
}

func (r *responder) WithTooManyRequests(w http.ResponseWriter) string {
	return r.withError(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
}

func (r *responder) withError(w http.ResponseWriter, message string, code int) string {

	errorUuid := uuid.New().String()

	res := myResponse{
		Meta: meta{
			Code:    code,
			Message: message,
			DebugId: errorUuid,
		},
	}

	r.writeResponse(w, res, code)
	return errorUuid
}

func (r *responder) writeResponse(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	data, err := json.Marshal(v)
	if err != nil {
		r.logger.Error("encoding response error: ", zap.String("reason", err.Error()))
		return
	}

	_, err = w.Write(data)
	if err != nil {
		r.logger.Error("writing response error: ", zap.String("reason", err.Error()))
	}
}
