package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/about", handleAbout)

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
