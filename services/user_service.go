package services

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/twaananen/boulderlog/db"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/utils"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db db.Database
}

func NewUserService(db db.Database) *UserService {
	return &UserService{db: db}
}

func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	// utils.LogInfo("Login request for user: " + username)

	user, err := s.db.GetUserByUsername(username)
	if err != nil {
		return "", err
	}

	if user == nil {
		// User not found, create a new one
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			utils.LogError("Error hashing password", err)
			return "", err
		}

		newUser := &models.User{
			Username: username,
			Password: string(hashedPassword),
		}

		if err := s.db.CreateUser(newUser); err != nil {
			utils.LogError("Error creating user", err)
			return "", err
		}
		utils.LogInfo("New user created: " + username)
	} else {
		// User found, check password
		// utils.LogInfo("Checking password for user: " + username)
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			// utils.LogInfo("Password does not match for user: " + username)
			return "", models.ErrInvalidCredentials
		}
		// utils.LogInfo("Password matches for user: " + username)
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(14 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(utils.JWTSecret)
	if err != nil {
		utils.LogError("Error generating token", err)
		return "", err
	}

	utils.LogInfo("User logged in: " + username)
	return tokenString, nil
}

func (s *UserService) IsUserLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return utils.JWTSecret, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return false
	}

	if time.Now().Unix() > int64(exp) {
		return false
	}

	return true
}

func (s *UserService) GetUsernameFromToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return utils.JWTSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return "", fmt.Errorf("invalid expiration claim")
	}

	if time.Now().Unix() > int64(exp) {
		return "", fmt.Errorf("token has expired")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in token")
	}

	return username, nil
}

func (s *UserService) CreateUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return s.db.CreateUser(user)
}

func (s *UserService) RefreshToken(r *http.Request, w http.ResponseWriter) error {
	username, err := s.GetUsernameFromToken(r)
	if err != nil {
		return err
	}

	// Generate new token with refreshed expiration
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(14 * 24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(utils.JWTSecret)
	if err != nil {
		return err
	}

	// Set the new cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(14 * 24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
		Secure:   true,
	})

	return nil
}
