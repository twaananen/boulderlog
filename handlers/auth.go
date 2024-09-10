package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/twaananen/boulderlog/components"
	"github.com/twaananen/boulderlog/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	userFile    = "data/users.csv"
	tokenExpiry = 24 * time.Hour
)

var jwtSecret []byte

// InitJWTSecret initializes the JWT secret from the environment variable
func InitJWTSecret() error {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return fmt.Errorf("JWT_SECRET environment variable is not set")
	}
	jwtSecret = []byte(secret)
	return nil
}

type User struct {
	Username string
	Password string
}

func AuthStatus(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := IsUserLoggedIn(r)

	if isLoggedIn {
		components.AuthStatusLoggedIn().Render(r.Context(), w)
	} else {
		components.AuthStatusLoggedOut().Render(r.Context(), w)
	}
}

func IsUserLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		exp := claims["exp"].(float64)
		return int64(exp) >= time.Now().Unix()
	}

	return false
}

func GetUsernameFromSession(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), nil
	}

	return "", fmt.Errorf("invalid token")
}

func LoginModal(w http.ResponseWriter, r *http.Request) {
	components.LoginModal().Render(r.Context(), w)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	components.Login("", "").Render(r.Context(), w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := findUser(username)
	if err != nil {
		// User not found, create a new one
		if err := createUser(username, password); err != nil {
			utils.LogError("Error creating user", err)
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
	} else {
		// User found, check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			components.LoginForm("Username and password do not match", username).Render(r.Context(), w)
			return
		}
	}

	// Generate JWT token
	token, err := generateToken(username)
	if err != nil {
		utils.LogError("Error generating token", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Set token as cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(tokenExpiry),
		HttpOnly: true,
		Path:     "/",
	})

	utils.LogInfo(fmt.Sprintf("User %s logged in", username))
	// Redirect to home page
	w.Header().Set("HX-Redirect", "/")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// Remove the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Redirect to home page with a full page reload
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func findUser(username string) (*User, error) {
	// Ensure the data directory exists
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		utils.LogError("Failed to create data directory", err)
		return nil, err
	}

	file, err := os.Open(userFile)
	if err != nil {
		if os.IsNotExist(err) {
			utils.LogError("User file not found", err)
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		utils.LogError("Failed to read user file", err)
		return nil, err
	}
	for _, record := range records {
		if record[0] == username {
			utils.LogInfo(fmt.Sprintf("User %s found", username))
			return &User{Username: record[0], Password: record[1]}, nil
		}
	}

	utils.LogInfo(fmt.Sprintf("User %s not found", username))
	return nil, fmt.Errorf("user not found")
}

func createUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		utils.LogError("Failed to hash password", err)
		return err
	}

	// Ensure the data directory exists
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		utils.LogError("Failed to create data directory", err)
		return err
	}

	file, err := os.OpenFile(userFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		utils.LogError("Failed to open user file", err)
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	utils.LogInfo(fmt.Sprintf("User %s created", username))
	return writer.Write([]string{username, string(hashedPassword)})
}

func generateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(tokenExpiry).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	// Check if user is logged in
	if !IsUserLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Render the profile page
	components.Profile().Render(r.Context(), w)
}
