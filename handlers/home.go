package handlers

import (
	"net/http"

	"github.com/twaananen/boulderlog/components"
)

func Home(w http.ResponseWriter, r *http.Request) {
	isLoggedIn := IsUserLoggedIn(r)
	if isLoggedIn {
		username, _ := GetUsernameFromSession(r)
		gradeCounts, toppedCounts, err := getTodayGradeCounts(username)
		if err != nil {
			http.Error(w, "Failed to get grade counts", http.StatusInternalServerError)
			return
		}
		components.Home(isLoggedIn, gradeCounts, toppedCounts, false).Render(r.Context(), w)
	} else {
		components.Home(isLoggedIn, nil, nil, false).Render(r.Context(), w)
	}
}
