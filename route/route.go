package route

import (
	"my_vocab/controller"
	"my_vocab/middleware"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func InitRoute() {
	mux := mux.NewRouter()
	var imgServer = http.FileServer(http.Dir("./public/"))
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", imgServer))
	mux.HandleFunc("/register", controller.Register).Methods("POST")
	mux.HandleFunc("/login", controller.Login).Methods("POST")

	// Vocabulary
	eJwt := mux.PathPrefix("/api").Subrouter()
	eJwt.HandleFunc("/vocabulary", controller.PostVocab).Methods("POST")
	eJwt.HandleFunc("/vocabulary", controller.PatchVocab).Methods("PATCH")
	eJwt.HandleFunc("/vocabulary/{id_user}/{type_request}/{keyword}", controller.GetListVocabulary).Methods("GET")
	eJwt.Use(middleware.JWTMiddleware)
	http.ListenAndServe(":8080", middleware.CustomLogger(os.Stderr, mux))
}
