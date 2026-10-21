package main

import "net/http"

func handlerErr(w http.ResponseWriter, r *http.Request) {
	//	This is a very specific function signature that is necessary to define an HTTP Handler.
	//	It always takes a ResponseWriter as the first parameter and a pointer to an HTTP request as the second parameter
	respondWithError(w, 400, "Something went wrong")

}
