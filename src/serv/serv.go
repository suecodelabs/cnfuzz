package serv

import (
	"fmt"
	"net/http"

	"github.com/suecodelabs/cnfuzz/src/log"
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
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.L().Fatal(err)
	}
}
