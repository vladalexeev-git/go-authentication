package service

import (
	"context"
	"fmt"
	"go-authentication/config"
	"go-authentication/pkg/utils"
	"net/url"
	"strings"
)

type socialAuthService struct {
	cfg *config.Config

	sessionService Session
	accountService Account
}

func NewSocialAuth(
	cfg *config.Config,
	a Account,
	s Session) *socialAuthService {

	return &socialAuthService{
		cfg:            cfg,
		sessionService: s,
		accountService: a,
	}
}

func (s *socialAuthService) AuthorizationURL(ctx context.Context, provider string) (*url.URL, error) {
	const op = "service.AuthorizationURL"
	provider = strings.ToLower(provider)

	scope, err := utils.UniqueString(32)
	if err != nil {
		return nil, fmt.Errorf("%s: gen unique string for scope: %w", op, err)
	}

	u, err := url.Parse(s.cfg.SocialAuth.Endpoints()[provider].AuthURL)
	if err != nil {
		return nil, fmt.Errorf("%s: parsing auth url error: %w", op, err)
	}

	q := u.Query()
	q.Set("client_id", s.cfg.ClientIDs()[provider])
	q.Set("scope", s.cfg.Scopes()[provider])
	q.Set("state", scope)
	u.RawQuery = q.Encode()

	return u, nil
}
