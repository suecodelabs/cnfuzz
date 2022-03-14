package serv

import (
	"fmt"
	"net/http"
)

func livez(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Alive.\n")
}

func readyz(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Ready.\n")
}

func Serv() {

	http.HandleFunc("/livez", livez)
	http.HandleFunc("/readyz", readyz)
	http.ListenAndServe(":80", nil)

}
