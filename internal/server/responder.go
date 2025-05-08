package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"net/http"

	"github.com/ntquang98/go-rkinetics-service/internal/pkg/common"
)

type httpResponder struct {
	hostname string
	start    time.Time
}

func NewHTTPResponder(hostname string) *httpResponder {
	return &httpResponder{
		hostname: hostname,
		start:    time.Now(),
	}
}

func (res *httpResponder) Respond(w http.ResponseWriter, r *http.Request, _response common.Response) error {
	if _response == nil {
		http.Error(w, "invalid response", http.StatusInternalServerError)
		return nil
	}

	response := _response.Interface()

	dif := float64(time.Since(res.start).Milliseconds())
	w.Header().Set("X-Execution-Time", fmt.Sprintf("%.4f ms", dif))
	w.Header().Set("X-Hostname", res.hostname)

	if id, ok := r.Context().Value(common.K_XRequestID).(uint64); ok {
		w.Header().Set(common.K_XRequestID, strconv.FormatUint(id, 10))
	}

	// Map status to HTTP code
	statusCode := http.StatusOK
	switch response.Status {
	case "NO_CONTENT":
		w.WriteHeader(http.StatusNoContent)
		return nil
	case common.APIStatus.Ok:
		statusCode = http.StatusOK
	case common.APIStatus.Error:
		statusCode = http.StatusInternalServerError
	case common.APIStatus.Invalid:
		statusCode = http.StatusBadRequest
	case common.APIStatus.NotFound:
		statusCode = http.StatusNotFound
	case common.APIStatus.Forbidden:
		statusCode = http.StatusForbidden
	case common.APIStatus.Unauthorized:
		statusCode = http.StatusUnauthorized
	case common.APIStatus.Existed:
		statusCode = http.StatusConflict
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(response)
}
