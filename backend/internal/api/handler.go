// Package api exposes the calculator over HTTP. It only translates between
// JSON and the calculator package; it contains no calculation logic.
package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"mime"
	"net/http"

	"sezzle-calculator/backend/internal/calculator"
)

const maxBodyBytes = 4 << 10

type calculateRequest struct {
	// Pointer so a missing field can be told apart from an empty string.
	Expression *string `json:"expression"`
}

type calculateResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

type handler struct {
	evaluate func(string) (float64, error)
}

// NewHandler returns the HTTP handler serving the API routes.
func NewHandler() http.Handler {
	return newHandler(calculator.Evaluate)
}

func newHandler(evaluate func(string) (float64, error)) http.Handler {
	h := &handler{evaluate: evaluate}
	mux := http.NewServeMux()
	mux.HandleFunc("/api/calculate", h.calculate)
	mux.HandleFunc("/api/health", health)
	return recoverPanics(mux)
}

func (h *handler) calculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	if mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported media type")
		return
	}

	var req calculateRequest
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes))
	dec.DisallowUnknownFields()
	err := dec.Decode(&req)
	if err == nil {
		// Anything after the first JSON value makes the request invalid.
		err = dec.Decode(&struct{}{})
		if err == io.EOF {
			err = nil
		} else if err == nil {
			err = errors.New("trailing data")
		}
	}
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, "request too large")
			return
		}
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.Expression == nil {
		writeError(w, http.StatusBadRequest, "invalid request")
		return
	}

	result, err := h.evaluate(*req.Expression)
	if err != nil {
		status, message := mapEvaluationError(err)
		writeError(w, status, message)
		return
	}
	writeJSON(w, http.StatusOK, calculateResponse{Result: result})
}

// mapEvaluationError turns calculator errors into stable client messages.
// Unrecognized errors are reported generically so internals never leak.
func mapEvaluationError(err error) (int, string) {
	switch {
	case errors.Is(err, calculator.ErrEmpty):
		return http.StatusBadRequest, "empty expression"
	case errors.Is(err, calculator.ErrSyntax):
		return http.StatusBadRequest, "invalid expression"
	case errors.Is(err, calculator.ErrDivisionByZero):
		return http.StatusBadRequest, "division by zero"
	case errors.Is(err, calculator.ErrInvalidSqrt):
		return http.StatusBadRequest, "invalid square root"
	case errors.Is(err, calculator.ErrInvalidNumber):
		return http.StatusBadRequest, "invalid number"
	case errors.Is(err, calculator.ErrTooComplex):
		return http.StatusBadRequest, "expression too complex"
	}
	log.Printf("unexpected evaluation error: %v", err)
	return http.StatusInternalServerError, "internal error"
}

func health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func recoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic serving %s %s: %v", r.Method, r.URL.Path, rec)
				writeError(w, http.StatusInternalServerError, "internal error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("encoding response: %v", err)
	}
}
