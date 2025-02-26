package config

import "time"

const (
	AccessTokenExpireDuration  = time.Hour * 24
	RefreshTokenExpireDuration = time.Hour * 24 * 7
	AccessTokenSubject         = "at"
	RefreshTokenSubject        = "rt"
	AuthMiddleWareContextKey   = "claim"
)
