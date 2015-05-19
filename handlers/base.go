// Package handlers provides request handlers.
package handlers

import (
	"errors"
	"github.com/tlehman/goga.me/dal"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"net/http"
	"strconv"
)

func getCurrentUser(w http.ResponseWriter, r *http.Request) *dal.UserRow {
	cookieStore := context.Get(r, "cookieStore").(*sessions.CookieStore)
	session, _ := cookieStore.Get(r, "goga.me-session")
	return session.Values["user"].(*dal.UserRow)
}

func getIdFromPath(w http.ResponseWriter, r *http.Request) (int64, error) {
	userIdString := mux.Vars(r)["id"]
	if userIdString == "" {
		return -1, errors.New("user id cannot be empty.")
	}

	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return -1, err
	}

	return userId, nil
}
