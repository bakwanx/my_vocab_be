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
	mux.HandleFunc("/refresh-token", controller.RefreshToken).Methods("POST")

	// Vocabulary
	eJwt := mux.PathPrefix("/api").Subrouter()
	eJwt.HandleFunc("/vocabulary", controller.PostVocab).Methods("POST")
	eJwt.HandleFunc("/vocabulary", controller.PatchVocab).Methods("PATCH")
	eJwt.HandleFunc("/vocabulary/{id_user}", controller.GetVocabularyByOrder).Methods("GET")
	eJwt.HandleFunc("/vocabulary/date/{id_user}", controller.GetVocabularyByDate).Methods("GET")
	eJwt.HandleFunc("/vocabulary/search", controller.GetVocabularyBySearch).Methods("POST")
	eJwt.HandleFunc("/type/vocabulary", controller.PostTypeVocab).Methods("POST")
	eJwt.Use(middleware.MiddlewareJWTAuthorization)
	http.ListenAndServe(":8080", middleware.CustomLogger(os.Stderr, mux))
}
