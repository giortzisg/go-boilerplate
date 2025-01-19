package server

import "context"

type Server interface {
	Start(context.Context) error
}
