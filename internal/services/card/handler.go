package card

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/internal/retry"
	"go.uber.org/zap"
)

func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]any{"error": msg})
}

type HandlerFactory struct {
	log    *zap.SugaredLogger
	mux    *http.ServeMux
	client v2.APIClientInterface
}

func NewHandlerFactory(
	logger *zap.SugaredLogger,
	client v2.APIClientInterface,
) *HandlerFactory {
	return &HandlerFactory{
		log:    logger,
		mux:    http.NewServeMux(),
		client: client,
	}
}

func (h *HandlerFactory) NewHandler(ctx context.Context) http.Handler {
	h.RegisterHealthzHandler(ctx)
	h.RegisterV1CardsHandler(ctx)
	return h.mux
}

func (h *HandlerFactory) RegisterHealthzHandler(ctx context.Context) {
	pattern := "/healthz"
	h.mux.HandleFunc(pattern, func(w http.ResponseWriter, _ *http.Request) {

		status := http.StatusOK
		h.log.Infow("request for", "service", ServiceName, "path", pattern,
			"status", status,
		)

		w.WriteHeader(status)
	})
}

func (h *HandlerFactory) RegisterV1CardsHandler(ctx context.Context) {
	pattern := "/v1/cards"
	h.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {

		cardID, err := ParseCardQuery(r)
		if err != nil {
			h.log.Infow("request bad",
				"service", ServiceName,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"error", err.Error(),
			)
			WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		result, retryable := h.client.MakeRetryable(cardID.ToURL())

		start := time.Now()
		err = retry.RunRetryable(ctx, retryable)
		end := time.Now()
		elapsed := end.Sub(start)

		if err != nil {
			h.log.Infow("client error",
				"service", ServiceName,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
				"error", err.Error(),
			)
			WriteError(w, http.StatusInternalServerError, err.Error())
			return
		}

		h.log.Infow("serving request",
			"service", ServiceName,
			"path", r.URL.Path,
			"cardID", cardID,
			"elapsed", elapsed.String(),
		)

		WriteJSON(w, http.StatusOK,
			map[string]any{
				"service": ServiceName,
				"ok":      true,
				"cardID":  cardID,
				"result":  json.RawMessage(result.Body),
			})
	})
}
