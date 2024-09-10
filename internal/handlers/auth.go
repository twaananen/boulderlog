package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/twaananen/boulderlog/internal/templates"
	"golang.org/x/crypto/bcrypt"
)

const (
	userFile    = "users.csv"
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
	isLoggedIn := isUserLoggedIn(r)

	if isLoggedIn {
		templates.AuthStatusLoggedIn().Render(r.Context(), w)
	} else {
		templates.AuthStatusLoggedOut().Render(r.Context(), w)
	}
}

func isUserLoggedIn(r *http.Request) bool {
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

func LoginModal(w http.ResponseWriter, r *http.Request) {
	templates.LoginModal().Render(r.Context(), w)
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	templates.Login("", "").Render(r.Context(), w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := findUser(username)
	if err != nil {
		// User not found, create a new one
		if err := createUser(username, password); err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}
	} else {
		// User found, check password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
			templates.LoginForm("Username and password do not match", username).Render(r.Context(), w)
			return
		}
	}

	// Generate JWT token
	token, err := generateToken(username)
	if err != nil {
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
	file, err := os.Open(userFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	for _, record := range records {
		if record[0] == username {
			return &User{Username: record[0], Password: record[1]}, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func createUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(userFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

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
	if !isUserLoggedIn(r) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Render the profile page
	templates.Profile().Render(r.Context(), w)
}
