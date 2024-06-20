package routes

import (
	"fmt"
	"net/http"
)

// handler for the root of the site
func Root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "welcome to dumbbin! send a POST request to this page with the bin content as the request body to create a new page")
	}
}

// handler for ratelimiting
func RateLimited() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintln(w, "you have been rate limited! try again later")
	}
}
