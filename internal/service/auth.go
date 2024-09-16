package service

import (
	"context"
	"errors"
	"time"

	"github.com/bojackodin/notes/internal/entity"
	"github.com/bojackodin/notes/internal/repository"
	"github.com/bojackodin/notes/internal/repository/repositoryerror"
	"github.com/bojackodin/notes/internal/service/serviceerror"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
}

type AuthService struct {
	userRepository repository.User
	secret         string
	tokenTTL       time.Duration
}

func NewAuthService(userRepository repository.User, secret string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepository: userRepository,
		secret:         secret,
		tokenTTL:       tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, username, password string) (int64, error) {
	hash, err := generatePasswordHash(password)
	if err != nil {
		return 0, err
	}

	user := entity.User{
		Username: username,
		Password: hash,
	}

	userId, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, repositoryerror.ErrDuplicate) {
			return 0, serviceerror.ErrUserDuplicate
		}
		return 0, err
	}
	return userId, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, username, password string) (string, error) {
	user, err := s.userRepository.GetUserByUsername(ctx, username)
	if err != nil {
		return "", nil
	}

	match, err := matches(user.Password, password)
	if err != nil {
		return "", err
	}
	if !match {
		return "", errors.New("don't match")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
	})

	return token.SignedString([]byte(s.secret))
}

func (s *AuthService) ParseToken(accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserID, nil
}

func generatePasswordHash(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func matches(hash []byte, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}
