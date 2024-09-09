package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/internal/templates"
)

func AuthStatus(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual authentication check
	isLoggedIn := false

	if isLoggedIn {
		w.Write([]byte(`<button hx-post="/auth/logout" hx-swap="outerHTML" class="bg-blue-500 hover:bg-blue-400 px-4 py-2 rounded">Logout</button>`))
	} else {
		w.Write([]byte(`<button hx-get="/auth/login-modal" hx-target="#login-modal" hx-swap="innerHTML" class="bg-blue-500 hover:bg-blue-400 px-4 py-2 rounded">Login</button>`))
	}
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
	templates.Login().Render(r.Context(), w)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual login logic
	username := r.FormValue("username")
	password := r.FormValue("password")

	// For demonstration purposes, we'll just check if both fields are non-empty
	if username != "" && password != "" {
		// Redirect to home page after successful login
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		// If login fails, redirect back to login page
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement actual logout logic
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
