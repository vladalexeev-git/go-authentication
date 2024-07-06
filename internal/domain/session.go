package domain

import (
	"go-authentication/internal/apperrors"
	"go-authentication/pkg/utils"
	"time"
)

type Session struct {
	ID        string    `json:"id" bson:"_id"`
	AccountID string    `json:"accountId" bson:"accountId"`
	Provider  string    `json:"provider" bson:"provider"`
	UserAgent string    `json:"userAgent" bson:"userAgent"`
	IP        string    `json:"ip" bson:"ip"`
	TTL       int       `json:"ttl" bson:"ttl"`
	ExpiresAt int64     `json:"expiresAt" bson:"expiresAt"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

func NewSession(aid, provider, userAgent, ip string, ttl time.Duration) (Session, error) {
	id, err := utils.UniqueString(32) // todo isn't uuid better?
	if err != nil {
		return Session{}, apperrors.ErrorSessionNotCreated
	}

	now := time.Now()

	return Session{
		ID:        id,
		AccountID: aid,
		Provider:  provider,
		UserAgent: userAgent,
		IP:        ip,
		TTL:       int(ttl.Seconds()),
		ExpiresAt: now.Add(ttl).Unix(),
		CreatedAt: now,
	}, nil
}
