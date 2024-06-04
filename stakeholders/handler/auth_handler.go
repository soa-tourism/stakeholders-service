package handler

import (
	"context"
	"stakeholders/proto/auth"
	"stakeholders/service"
)

type AuthHandler struct {
	AuthService *service.AuthService
	auth.UnimplementedAuthorizeServer
}

func (h AuthHandler) Login(ctx context.Context, request *auth.Credentials) (*auth.AuthenticationTokens, error) {
	return h.AuthService.Login(request)
}

func (h AuthHandler) Register(ctx context.Context, request *auth.AccountRegistrationRequest) (*auth.AccountRegistrationResponse, error) {
	return h.AuthService.Register(request)
}
