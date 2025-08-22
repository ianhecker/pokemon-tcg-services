package card

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HandlerFactory struct {
	log *zap.SugaredLogger
	mux *http.ServeMux
}

func NewHandlerFactory(
	logger *zap.SugaredLogger,
) *HandlerFactory {
	return &HandlerFactory{
		log: logger,
		mux: http.NewServeMux(),
	}
}

func (h *HandlerFactory) NewHandler() http.Handler {
	h.RegisterHealthzHandler()
	h.RegisterV1PingHandler()
	return h.mux
}

func (h *HandlerFactory) RegisterHealthzHandler() {

	pattern := "/healthz"
	h.mux.HandleFunc(pattern, func(w http.ResponseWriter, _ *http.Request) {

		status := http.StatusOK
		h.log.Infow("service request",
			"service", ServiceName,
			"path", pattern,
			"status", status,
		)

		w.WriteHeader(status)
	})
}

func (h *HandlerFactory) RegisterV1PingHandler() {

	pattern := "/v1/cards"
	h.mux.HandleFunc(pattern, func(w http.ResponseWriter, _ *http.Request) {

		response := fmt.Sprintf(`{"service":"%s","ok":true}`, ServiceName)
		h.log.Infow("request",
			"service", ServiceName,
			"path", pattern,
			"response", response,
		)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(response))
	})
}
