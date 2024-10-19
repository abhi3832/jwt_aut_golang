package main

import (
	"net/http"

)

func serverCheckHandler(w http.ResponseWriter, r *http.Request){

	response := map[string]string{
        "status":  "success",
        "message": "Server is running !",
    }

	respondWithJSON(w,200,response)

}

