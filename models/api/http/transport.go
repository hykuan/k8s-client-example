package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/hykuan/k8s-client-example"
	log "github.com/hykuan/k8s-client-example/logger"
	"github.com/hykuan/k8s-client-example/models"
)

const contentType = "application/json"

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	logger                    log.Logger
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc models.Service, l log.Logger) http.Handler {
	logger = l

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	mux := bone.New()

	mux.Post("/training", kithttp.NewServer(
		startTrainingEndpoint(svc),
		decodeTrainingReq,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/version", quai.Version("models"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func decodeTrainingReq(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		logger.Warn("Invalid or missing content type.")
		return nil, errUnsupportedContentType
	}

	var training models.Training
	if err := json.NewDecoder(r.Body).Decode(&training); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode persistent volume: %s", err))
		return nil, err
	}

	return trainingReq{training: training}, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", contentType)

	if ar, ok := response.(quai.Response); ok {
		for k, v := range ar.Headers() {
			w.Header().Set(k, v)
		}

		w.WriteHeader(ar.Code())

		if ar.Empty() {
			return nil
		}
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", contentType)

	switch err {
	case models.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case models.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case models.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	default:
		if statusError, isStatus := err.(*k8sErrors.StatusError); isStatus {
			w.WriteHeader(int(statusError.Status().Code))
		}
		switch err.(type) {
		case *json.SyntaxError:
			w.WriteHeader(http.StatusBadRequest)
		case *json.UnmarshalTypeError:
			w.WriteHeader(http.StatusBadRequest)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
