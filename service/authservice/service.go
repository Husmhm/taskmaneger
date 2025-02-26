package authservice

import (
	"github.com/golang-jwt/jwt/v4"
	"strings"
	models "taskmaneger/model"
	"time"
)

type Config struct {
	SignKey               string        `koanf:"sign_key"`
	AccessExpirationTime  time.Duration `koanf:"access_expiration_time"`
	RefreshExpirationTime time.Duration `koanf:"refresh_expiration_time"`
	AccessTokenSubject    string        `koanf:"access_subject"`
	RefreshTokenSubject   string        `koanf:"refresh_subject"`
}

type Service struct {
	config Config
}

func New(cfg Config) Service {
	return Service{
		config: cfg,
	}

}

func (s Service) CreateAccessToken(user models.User) (string, error) {
	return s.createToken(user.ID, s.config.AccessExpirationTime, s.config.AccessTokenSubject)
}

func (s Service) CreateRefreshToken(user models.User) (string, error) {
	return s.createToken(user.ID, s.config.RefreshExpirationTime, s.config.RefreshTokenSubject)
}

func (s Service) ParseToken(bearerToken string) (*Claims, error) {

	tokenStr := strings.Replace(bearerToken, "Bearer ", "", 1)

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.SignKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func (s Service) createToken(userID uint, expireDuratipn time.Duration, subject string) (string, error) {

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuratipn)),
			Subject:   subject,
		},
		UserID: userID,
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := accessToken.SignedString([]byte(s.config.SignKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
