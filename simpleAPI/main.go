package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/about", handleAbout)
	mux.HandleFunc("GET /comments", handleGetComments)
	mux.HandleFunc("GET /comments/{id}", handleSingleComment)
	mux.HandleFunc("POST /comments", handlePostComment)

	fmt.Println("Listening on port 8080! 📡")

	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err.Error())
		panic("Fatal error: Unable to serve and listen")
	}

}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}
func handleAbout(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About")
}
func handleGetComments(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "All comments here")
}
func handlePostComment(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Posted new comment!")
}
func handleSingleComment(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	fmt.Fprintf(w, "Got single comment for id: %s!", id)
}
