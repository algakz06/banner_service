package usecase

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/algakz/banner_service/models"
	"github.com/algakz/banner_service/pkg/auth"
	jwt "github.com/golang-jwt/jwt/v5"
)

type AuthClaims struct {
	jwt.RegisteredClaims
	User *models.User `json:"user"`
}

type AuthUseCase struct {
	userRepo       auth.UserRepository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthUseCase(
	userRepo auth.UserRepository,
	hashSalt string,
	signingKey []byte,
	tokenTTLSeconds time.Duration) *AuthUseCase {
	return &AuthUseCase{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (a *AuthUseCase) SignUp(ctx context.Context, username, password string) error {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))

	user := &models.User{
		Username: username,
		Password: fmt.Sprintf("%x", pwd.Sum(nil)),
	}

	return a.userRepo.CreateUser(ctx, user)
}

func (a *AuthUseCase) SignIn(ctx context.Context, username, password string) (string, error) {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(a.hashSalt))
	password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userRepo.GetUser(ctx, username, password)
	if err != nil {
		return "", auth.ErrUserNotFound
	}

	claims := AuthClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.signingKey)
}

func (a *AuthUseCase) ParseToken(ctx context.Context, accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		return a.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok {
		return claims.User, nil
	}

	return nil, auth.ErrInvalidAccessToken
}
