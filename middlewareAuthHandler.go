package main

import (
	"fmt"
	"net/http"
	"time"

	auth "github.com/abhi3832/jwt_auth/Auth"
	"github.com/abhi3832/jwt_auth/internal/database"
	"github.com/google/uuid"
)


type authedHandler func(http.ResponseWriter, *http.Request, database.User)//This means any function that follows this 
//signature can be referred to as an authedHandler.

func(apicfg *apiConfig) middlewareAuthByToken(handler authedHandler) http.HandlerFunc{

	return func(w http.ResponseWriter, r* http.Request){

		token,err := auth.Get_token(r.Header)

		if err != nil{
			respondWithError(w,403,fmt.Sprintf("Error Getting token - %v",err))
			return
		}

		claim,err := validateJWTToken(token)

		if err != nil{
			respondWithError(w,401,fmt.Sprintf("Error Getting token - %v",err))
			return
		}

		// token expiry
		if exp, ok := claim["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if time.Now().After(expirationTime) {
				respondWithError(w,403,fmt.Sprintf("token expired!"))
				return
			}
		}

		userIDstr := claim["sub"].(string)
		
		user_id,err := uuid.Parse(userIDstr)

		if err != nil{
			respondWithError(w,401,"Error Extracting User Id from Token")
			return
		}


		user,err := apicfg.DB.GetUserByUserId(r.Context(),user_id)

		if err != nil{
			respondWithError(w,403,"Unauthorized Access !")
			return
		}

		handler(w,r,user)


	}

}

