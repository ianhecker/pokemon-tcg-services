package card

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ianhecker/pokemon-tcg-services/internal/pokemonpricetracker"
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
	client pokemonpricetracker.APIClientInterface
}

func NewHandlerFactory(
	logger *zap.SugaredLogger,
	client pokemonpricetracker.APIClientInterface,
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

		card, err := ParseCardQuery(r)
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

		start := time.Now()
		err = h.client.GetPricing(ctx, &card)
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
			"ID", card.ID,
			"elapsed", elapsed.String(),
		)

		WriteJSON(w, http.StatusOK,
			map[string]any{
				"service": ServiceName,
				"ok":      true,
				"ID":      card.ID,
				"result":  card,
			})
	})
}
