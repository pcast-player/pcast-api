package service_interface

import "context"

type OAuth interface {
	GetGoogleAuthURL(state string) (string, error)
	HandleGoogleCallback(ctx context.Context, code, state string) (string, error)
}
