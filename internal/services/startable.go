package services

import "context"

type StartableInterface interface {
	Start(ctx context.Context) (stop func())
	Done() <-chan struct{}
	Err() error
}
