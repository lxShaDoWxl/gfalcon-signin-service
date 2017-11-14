package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/justinas/nosurf"
	"github.com/m0cchi/gfalcon/complex"
	"github.com/m0cchi/gfalcon/model"
	"github.com/m0cchi/gfalcon/util"
	"html/template"
	"log"
	"net/http"
	"os"
)

const MaxAge = 1 * 60 * 60 * 24

var db *sqlx.DB
var templates *template.Template
var allowedHost string

type AuthResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func init() {
	templates = template.Must(template.ParseGlob("resources/templates/*"))
}

func submitSigninForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	teamID := r.Form.Get("team-id")
	userID := r.Form.Get("user-id")
	password := r.Form.Get("password")

	authResponse := AuthResponse{Ok: false, Message: ""}

	team, err := model.GetTeam(db, teamID)
	if err != nil {
		authResponse.Message = fmt.Sprintf("%v", err)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(authResponse)
		return
	}

	user, err := model.GetUser(db, team.IID, userID)
	if err != nil {
		authResponse.Message = fmt.Sprintf("%v", err)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(authResponse)
		return
	}

	session, err := complex.AuthenticateWithPassword(db, user, password)
	if err != nil {
		authResponse.Message = fmt.Sprintf("%v", err)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(authResponse)
		return
	}

	if err = session.Validate(); err != nil {
		authResponse.Message = fmt.Sprintf("%v", err)
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(authResponse)
		return
	}

	expires := session.UpdateDate.AddDate(0, 0, 1)
	cgsession := http.Cookie{
		Name:     "gfalcon.session",
		Value:    session.SessionID,
		Expires:  expires,
		MaxAge:   MaxAge,
		HttpOnly: true,
	}

	cgiid := http.Cookie{
		Name:     "gfalcon.iid",
		Value:    fmt.Sprintf("%d", session.UserIID),
		Expires:  expires,
		MaxAge:   MaxAge,
		HttpOnly: true,
	}

	if allowedHost != "" {
		cgsession.Domain = allowedHost
		cgiid.Domain = allowedHost
	}
	http.SetCookie(w, &cgsession)
	http.SetCookie(w, &cgiid)

	authResponse.Ok = true
	json.NewEncoder(w).Encode(authResponse)
}

func signinForm(w http.ResponseWriter, r *http.Request) {
	args := map[string]interface{}{
		"csrf_token": nosurf.Token(r),
	}

	if err := templates.ExecuteTemplate(w, "signin.html.tmpl", args); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var secret string
	var port int
	var dbhost string
	var err error
	flag.StringVar(&secret, "secret", util.GenerateSessionID(32), "csrf's auth key")
	flag.IntVar(&port, "port", 8080, "service's port")
	flag.StringVar(&dbhost, "dbhost", "", "gfalcon's DB")
	flag.StringVar(&allowedHost, "allowed-host", "", "allowed host")
	flag.Parse()

	if dbhost == "" {
		fmt.Println("required --dbhost [host]")
		os.Exit(1)
	}

	db, err = sqlx.Connect("mysql", dbhost)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	r := mux.NewRouter()

	r.HandleFunc("/signin", submitSigninForm)
	r.HandleFunc("/", signinForm)

	fileServer := http.FileServer(http.Dir("resources/statics/"))
	r.PathPrefix("/statics/").Handler(http.StripPrefix("/statics/", fileServer))

	http.ListenAndServe(fmt.Sprintf(":%d", port), nosurf.New(r))
}
