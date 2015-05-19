package handlers

import (
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"github.com/tlehman/goga.me/dal"
	"github.com/tlehman/goga.me/libhttp"
	"html/template"
	"net/http"
)

func GetHome(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	cookieStore := context.Get(r, "cookieStore").(*sessions.CookieStore)

	session, _ := cookieStore.Get(r, "goga.me-session")
	currentUser, ok := session.Values["user"].(*dal.UserRow)
	if !ok {
		http.Redirect(w, r, "/logout", 301)
		return
	}

	data := struct {
		CurrentUser *dal.UserRow
	}{
		currentUser,
	}

	tmpl, err := template.ParseFiles("templates/dashboard.html.tmpl", "templates/home.html.tmpl")
	if err != nil {
		libhttp.HandleErrorJson(w, err)
		return
	}

	tmpl.Execute(w, data)
}
