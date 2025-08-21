package service

import (
	"context"

	v2 "github.com/ianhecker/pokemon-tcg-services/internal/pokemontcgio/v2"
)

type ServiceInterface interface {
	Start(ctx context.Context)
	Stop() error
}

type CardServiceInterface interface {
	ServiceInterface
}

type CardService struct {
}

func NewCardService(
	client v2.APIClientInterface,
) CardServiceInterface {
	return &CardService{}
}

func (svc CardService) Start(ctx context.Context) {}

func (svc CardService) Stop() error {
	return nil
}
