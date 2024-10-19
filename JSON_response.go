package main

import (
	"encoding/json"
	"log"
	"net/http"
)


func respondWithError(w http.ResponseWriter, code int, msg string){

	if code > 499{
		log.Println("Responding with 5XX ERROR !")
	}

	type ErrorResponse struct{
		Error string `json:"error"`
	}

	respondWithJSON(w,code,ErrorResponse{Error:msg})
}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}){

	data, err := json.Marshal(payload)

	if err != nil{
		log.Println("Failed to Marshal the Payload to Json : &v", payload)
		w.WriteHeader(500)
		return
	}else{
		w.Header().Add("Content-Type","application/json")
		w.WriteHeader(code)
		w.Write(data)
	}
}