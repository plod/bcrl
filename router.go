package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

var r = mux.NewRouter()

func init() {
	//generic pages
	r.HandleFunc("/", home)
	r.HandleFunc("/about", about)
	r.HandleFunc("/contact", contact)
	r.HandleFunc("/rules", rules)
	r.HandleFunc("/pictures", pictures)

	//login flow
	r.HandleFunc("/login", login)
	r.HandleFunc("/logout", logout)
	r.HandleFunc("/signup", signup)

	//dynamic content
	r.HandleFunc("/results/{year}/{event}", results)
	r.HandleFunc("/events/{year}/{event}", events)

	//static assets
	r.PathPrefix("/assets").Handler(fs)

	//logged in sub router
	au := r.PathPrefix("/authd").Subrouter()

	//for every user
	// * need to add middleWare to check user is authd at this point
	au.HandleFunc("/", authdHome)
	au.HandleFunc("/edit", authEdit)
	au.HandleFunc("/myresults", authMyResults)

	au.HandleFunc("/enter/{year}/{event}", authEnter)

	//for admins lots of middleware needed as described
	ad := r.PathPrefix("/admin").Subrouter()

	ad.HandleFunc("/club/add", adminClubAdd)          //needs middleware to make sure only admin can do
	ad.HandleFunc("/club/edit/{club}", adminClubEdit) //needs middleware only club primary or above can do

	ad.HandleFunc("/event/add", adminEventAdd)                  //needs middleware club admin can only add events for their own club
	ad.HandleFunc("/event/edit/{year}/{event}", adminEventEdit) //can only edit events in future unless super user

	ad.HandleFunc("/result/add", adminResultAdd)                  //needs middleware super user only
	ad.HandleFunc("/result/edit/{year}/{event}", adminResultEdit) //needs middleware super user only

	router.NotFoundHandler = http.HandlerFunc(notFound)
}
