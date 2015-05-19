package handlers

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/jmoiron/sqlx"
	"github.com/tlehman/goga.me/dal"
	"github.com/tlehman/goga.me/libhttp"
	"html/template"
	"net/http"
)

// Creates a new match.
//
// By default, new matches have both users equal, then later,
// another user can join if they want
func PostMatches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	cookieStore := context.Get(r, "cookieStore").(*sessions.CookieStore)

	session, _ := cookieStore.Get(r, "goga.me-session")
	currentUser, ok := session.Values["user"].(*dal.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 301)
		return
	}

	db := context.Get(r, "db").(*sqlx.DB)

	// create new match with only one user
	m := dal.NewMatch(db)
	matchRow, err := m.BeginMatch(nil, currentUser, currentUser)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
	}

	http.Redirect(w, r, fmt.Sprintf("/matches/%d", matchRow.ID), 301)
}

func GetMatch(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	cookieStore := context.Get(r, "cookieStore").(*sessions.CookieStore)

	session, _ := cookieStore.Get(r, "goga.me-session")
	currentUser, ok := session.Values["user"].(*dal.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 301)
		return
	}

	// fetch match from database
	db := context.Get(r, "db").(*sqlx.DB)

	m := dal.NewMatch(db)
	matchId := parseMatchId(r.URL.Path)
	matchRow, err := m.GetById(nil, matchId)
	if err != nil {
		libhttp.HandleErrorJson(w, err)
	}

	// fetch other user from the database
	u := dal.NewUser(db)
	// we can ignore err here since white_user_id is a foreign key in Postgres
	otherUser, _ := u.GetById(nil, matchRow.WhiteUserID)

	data := struct {
		CurrentUser *dal.UserRow
		OtherUser   *dal.UserRow
	}{
		currentUser,
		otherUser,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/matches/show.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}
	tmpl.Execute(w, data)
}

func parseMatchId(path string) int64 {
	return int64(1)
}
