package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pingan/bastion/internal/model"
	"github.com/pingan/bastion/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AuthService struct {
	userRepo  *repository.UserRepo
	jwtSecret []byte
}

func NewAuthService(userRepo *repository.UserRepo, jwtSecret string) *AuthService {
	return &AuthService{userRepo: userRepo, jwtSecret: []byte(jwtSecret)}
}

func (s *AuthService) Login(username, password string) (string, *model.User, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return "", nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}
	token, err := s.generateToken(user)
	if err != nil {
		return "", nil, err
	}
	return token, user, nil
}

func (s *AuthService) ValidateToken(tokenStr string) (int64, string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, "", err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", errors.New("invalid token")
	}
	userID := int64(claims["user_id"].(float64))
	role := claims["role"].(string)
	return userID, role, nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
