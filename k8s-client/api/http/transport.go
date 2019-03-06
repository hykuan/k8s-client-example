package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hykuan/k8s-client-example"
	"github.com/hykuan/k8s-client-example/k8s-client"
	"io"
	"net/http"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/go-zoo/bone"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/qeek-dev/quaistudio/logger"
	k8sErrors "k8s.io/apimachinery/pkg/api/errors"
)

const contentType = "application/json"

var (
	errUnsupportedContentType = errors.New("unsupported content type")
	logger                    log.Logger
)

// MakeHandler returns a HTTP handler for API endpoints.
func MakeHandler(svc k8s_client.Service, l log.Logger) http.Handler {
	logger = l

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorEncoder(encodeError),
	}

	mux := bone.New()

	mux.Post("/pv", kithttp.NewServer(
		createNFSPVEndpoint(svc),
		decodeNFSPersistentVolume,
		encodeResponse,
		opts...,
	))

	mux.Post("/pvc", kithttp.NewServer(
		createPVCEndpoint(svc),
		decodePersistentVolumeClaim,
		encodeResponse,
		opts...,
	))

	mux.Post("/deployment", kithttp.NewServer(
		createDeploymentEndpoint(svc),
		decodeDeployment,
		encodeResponse,
		opts...,
	))

	mux.GetFunc("/version", quai.Version("k8s-client"))
	mux.Handle("/metrics", promhttp.Handler())

	return mux
}

func decodeNFSPersistentVolume(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		logger.Warn("Invalid or missing content type.")
		return nil, errUnsupportedContentType
	}

	var pv k8s_client.NFSPersistentVolume
	if err := json.NewDecoder(r.Body).Decode(&pv); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode persistent volume: %s", err))
		return nil, err
	}

	return nfsPVReq{pv}, nil
}

func decodePersistentVolumeClaim(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		logger.Warn("Invalid or missing content type.")
		return nil, errUnsupportedContentType
	}

	var pvc k8s_client.PersistentVolumeClaim
	if err := json.NewDecoder(r.Body).Decode(&pvc); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode persistent volume: %s", err))
		return nil, err
	}

	return pvcReq{pvc}, nil
}

func decodeDeployment(_ context.Context, r *http.Request) (interface{}, error) {
	if r.Header.Get("Content-Type") != contentType {
		logger.Warn("Invalid or missing content type.")
		return nil, errUnsupportedContentType
	}

	var deployment k8s_client.Deployment
	if err := json.NewDecoder(r.Body).Decode(&deployment); err != nil {
		logger.Warn(fmt.Sprintf("Failed to decode persistent volume: %s", err))
		return nil, err
	}

	return deploymentReq{deployment}, nil
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
	case k8s_client.ErrMalformedEntity:
		w.WriteHeader(http.StatusBadRequest)
	case k8s_client.ErrUnauthorizedAccess:
		w.WriteHeader(http.StatusForbidden)
	case k8s_client.ErrConflict:
		w.WriteHeader(http.StatusConflict)
	case errUnsupportedContentType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
	case io.ErrUnexpectedEOF:
		w.WriteHeader(http.StatusBadRequest)
	case io.EOF:
		w.WriteHeader(http.StatusBadRequest)
	default:
		if  statusError, isStatus := err.(*k8sErrors.StatusError); isStatus {
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
