package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/getActivity", HandlerGetActivity)
	http.HandleFunc("/getKudos", HandlerGetKudos)

	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)

	// err = getKudos(token)
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
