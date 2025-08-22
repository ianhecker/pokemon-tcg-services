package card

import (
	"context"
	"net/http"
	"time"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
	"github.com/ianhecker/pokemon-tcg-services/internal/services"
	"go.uber.org/zap"
)

const ServiceName = "card"
const ServerShutdown = 10 * time.Second

type CardServiceInterface interface {
	services.ServiceInterface
}

type CardService struct {
	log      *zap.SugaredLogger
	client   v2.APIClientInterface
	srv      *http.Server
	done     chan struct{}
	err      error
	shutdown time.Duration
}

func NewService(
	logger *zap.SugaredLogger,
	client v2.APIClientInterface,
	addr string,
) CardServiceInterface {
	srv := &http.Server{
		Addr: addr,
	}
	handlerFactory := NewHandlerFactory(logger)
	handler := handlerFactory.NewHandler()

	return NewServiceFromRaw(
		logger,
		client,
		srv,
		handler,
		make(chan struct{}, 1),
		nil,
		ServerShutdown,
	)
}

func NewServiceFromRaw(
	logger *zap.SugaredLogger,
	client v2.APIClientInterface,
	srv *http.Server,
	handler http.Handler,
	done chan struct{},
	err error,
	shutdown time.Duration,
) CardServiceInterface {
	srv.Handler = handler
	return &CardService{
		log:      logger,
		client:   client,
		srv:      srv,
		done:     done,
		err:      err,
		shutdown: shutdown,
	}
}

func (svc *CardService) Start(ctx context.Context) (stop func()) {

	svc.log.Infow("server start", "address", svc.srv.Addr)
	go func() {
		err := svc.srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			svc.log.Errorw("server error", "err", err)
			svc.err = err
		}
		close(svc.done)
	}()
	go func() {
		<-ctx.Done()
		svc.log.Errorw("server context shutdown", "timeout", svc.shutdown)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), svc.shutdown)
		defer cancel()
		_ = svc.srv.Shutdown(shutdownCtx)
	}()
	return func() {
		svc.log.Infow("server manual shutdown", "timeout", svc.shutdown)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), svc.shutdown)
		defer cancel()
		_ = svc.srv.Shutdown(shutdownCtx)
	}
}

func (svc *CardService) Done() <-chan struct{} {
	return svc.done
}

func (svc *CardService) Err() error {
	return svc.err
}
