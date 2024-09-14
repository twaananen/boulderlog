package services

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/twaananen/boulderlog/models"
	"github.com/twaananen/boulderlog/utils"
	"golang.org/x/crypto/bcrypt"
)

// Remove the MockDatabase definition and its methods

func TestAuthenticateUser(t *testing.T) {
	mockDB := new(MockDatabase)
	userService := NewUserService(mockDB)

	// Test case 1: Successful authentication of existing user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	mockDB.On("GetUserByUsername", "existinguser").Return(&models.User{
		Username: "existinguser",
		Password: string(hashedPassword),
	}, nil).Once()

	token, err := userService.AuthenticateUser("existinguser", "password123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test case 2: Failed authentication due to wrong password
	mockDB.On("GetUserByUsername", "existinguser").Return(&models.User{
		Username: "existinguser",
		Password: string(hashedPassword),
	}, nil).Once()

	token, err = userService.AuthenticateUser("existinguser", "wrongpassword")
	assert.Error(t, err)
	assert.Empty(t, token)

	// Test case 3: New user creation
	mockDB.On("GetUserByUsername", "newuser").Return(nil, nil).Once()
	mockDB.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil).Once()

	token, err = userService.AuthenticateUser("newuser", "newpassword")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Test case 4: Database error
	mockDB.On("GetUserByUsername", "erroruser").Return(nil, assert.AnError).Once()

	token, err = userService.AuthenticateUser("erroruser", "password")
	assert.Error(t, err)
	assert.Empty(t, token)

	mockDB.AssertExpectations(t)
}

func TestIsUserLoggedIn(t *testing.T) {
	mockDB := new(MockDatabase)
	userService := NewUserService(mockDB)

	// Test case 1: User is logged in
	utils.JWTSecret = []byte("test_secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(utils.JWTSecret)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})

	assert.True(t, userService.IsUserLoggedIn(req))

	// Test case 2: User is not logged in (no cookie)
	req, _ = http.NewRequest("GET", "/", nil)
	assert.False(t, userService.IsUserLoggedIn(req))

	// Test case 3: Invalid token
	req, _ = http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid_token"})
	assert.False(t, userService.IsUserLoggedIn(req))
}

func TestGetUsernameFromToken(t *testing.T) {
	mockDB := new(MockDatabase)
	userService := NewUserService(mockDB)

	// Test case 1: Valid token
	utils.JWTSecret = []byte("test_secret")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": "testuser",
		"exp":      time.Now().Add(time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString(utils.JWTSecret)

	req, _ := http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: tokenString})

	username, err := userService.GetUsernameFromToken(req)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", username)

	// Test case 2: No token cookie
	req, _ = http.NewRequest("GET", "/", nil)
	username, err = userService.GetUsernameFromToken(req)
	assert.Error(t, err)
	assert.Empty(t, username)

	// Test case 3: Invalid token
	req, _ = http.NewRequest("GET", "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: "invalid_token"})
	username, err = userService.GetUsernameFromToken(req)
	assert.Error(t, err)
	assert.Empty(t, username)
}

func TestCreateUser(t *testing.T) {
	mockDB := new(MockDatabase)
	userService := NewUserService(mockDB)

	// Test case 1: Successful user creation
	mockDB.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil).Once()

	err := userService.CreateUser("newuser", "password123")
	assert.NoError(t, err)

	// Test case 2: Database error
	mockDB.On("CreateUser", mock.AnythingOfType("*models.User")).Return(assert.AnError).Once()

	err = userService.CreateUser("erroruser", "password123")
	assert.Error(t, err)

	mockDB.AssertExpectations(t)
}

// Other test functions remain the same
// ...
