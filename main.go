package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func LogRequest(w http.ResponseWriter, r *http.Request){
	fmt.Println(JSONEncode(r))
	w.WriteHeader(http.StatusOK)
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/log",LogRequest)
	http.Handle("/",r)

	if err := http.ListenAndServe(":8080", nil); err !=nil {
		fmt.Println(err)
	}
}
