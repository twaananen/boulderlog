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
		w.Write([]byte(`<button hx-post="/auth/logout" hx-swap="outerHTML" class="bg-blue-500 hover:bg-blue-400 px-4 py-2 rounded">Logout</button>`))
	} else {
		w.Write([]byte(`<a href="/login" class="bg-blue-500 hover:bg-blue-400 px-4 py-2 rounded">Login</a>`))
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
		if int64(exp) < time.Now().Unix() {
			return false
		}
		return true
	}

	return false
}

func LoginModal(w http.ResponseWriter, r *http.Request) {
	// This function will return the content of the login modal
	loginModalContent := `
			<div class="mt-3 text-center">
				<h3 class="text-lg leading-6 font-medium text-gray-900">Login</h3>
				<form class="mt-2 px-7 py-3" hx-post="/auth/login" hx-target="#login-modal">
					<input
						type="text"
						name="username"
						placeholder="Username"
						class="mt-2 px-3 py-2 bg-white border shadow-sm border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 block w-full rounded-md sm:text-sm focus:ring-1"
						required
					/>
					<input
						type="password"
						name="password"
						placeholder="Password"
						class="mt-2 px-3 py-2 bg-white border shadow-sm border-slate-300 placeholder-slate-400 focus:outline-none focus:border-sky-500 focus:ring-sky-500 block w-full rounded-md sm:text-sm focus:ring-1"
						required
					/>
					<button
						type="submit"
						class="mt-4 bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50"
					>
						Login
					</button>
				</form>
				<div class="items-center px-4 py-3">
					<button
						id="close-login"
						class="px-4 py-2 bg-gray-500 text-white text-base font-medium rounded-md w-full shadow-sm hover:bg-gray-600 focus:outline-none focus:ring-2 focus:ring-gray-300"
					>
						Close
					</button>
				</div>
			</div>
	`
	w.Write([]byte(loginModalContent))
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

	// Respond with a new login button
	w.Write([]byte(`<a href="/login" class="bg-blue-500 hover:bg-blue-400 px-4 py-2 rounded">Login</a>`))
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
