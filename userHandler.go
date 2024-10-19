package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abhi3832/jwt_auth/internal/database"
	"github.com/google/uuid"
)

type UserResponse struct {
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UserType  string `json:"user_type,omitempty"`
	ApiKey    string `json:"api_key,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}



func(apicfg * apiConfig) userLoginHandler(w http.ResponseWriter, r *http.Request){

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error Parsing JSON: %s", err))
		return
	}

	// check user already exits or not
	user,err := apicfg.DB.CheckUserAlreadyExits(r.Context(),database.CheckUserAlreadyExitsParams{
		Email: params.Email,
		Phone: "",
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error Fetching User or User Does not exist !!: %s", err))
		return
	}

	// if the passward matches then the user is found
	if CheckPassword(user.Passward, params.Password){

		session,err := apicfg.DB.GetSessionByUserId(r.Context(),user.UserID)

		if err == nil{
			errr := apicfg.DB.UpdateSessionByDelete(r.Context(),session.UserID)

			if errr != nil{
				respondWithError(w,403,"Could't Delete Your Session")
				return
			}
		}

		
		accessToken, err := GenerateAccessToken(user.UserID)
		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error Generating Access Token: %s", err))
			return
		}

		refreshToken, err := GenerateRefreshToken()

		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error Generating Refresh Token: %s", err))
			return
		}

		expiresAt := time.Now().Add(24 * time.Hour)

		// Create a new session record
		new_session, err := apicfg.DB.CreateSession(r.Context(), database.CreateSessionParams{
			UserID:       user.UserID,
			RefreshToken: refreshToken,
			CreatedAt:    sql.NullTime{Time :time.Now(), Valid: true},
			UpdatedAt:    sql.NullTime{Time :time.Now(), Valid: true},
			ExpiresAt:    expiresAt,
		})

		if err != nil {
			respondWithError(w, 500, fmt.Sprintf("Error Creating Session: %s", err))
			return
		}

		// Return the tokens and user information
		respondWithJSON(w, 200, map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": new_session.RefreshToken,
			"user": map[string]interface{}{
				"UserID":   user.UserID,
				"FirstName": user.FirstName,
				"LastName":  user.LastName.String,
				"Email":     user.Email,
				"Phone":     user.Phone,
				"UserType":  user.UserType.String,
			},
		})

	}else{
		respondWithError(w,401,"Invalid passward")
	}
	
}


func(apicfg *apiConfig) userSignupHandler(w http.ResponseWriter, r *http.Request){

	type parameters struct{
		First_name string `json:"first_name"`
		Last_name string `json:"last_name"`
		Email string `json:"email"`
		Phone string `json:"phone"`
		Passward string `json:"passward"`
		User_type string `json:"user_type"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	if err != nil{
		respondWithError(w,400, fmt.Sprintf("Error Parsing Json : %s", err))
		return
	}

	user,err := apicfg.DB.CheckUserAlreadyExits(r.Context(),database.CheckUserAlreadyExitsParams{
		Email: params.Email,
		Phone: params.Phone,
	})

	if err != nil{
		respondWithError(w,400, fmt.Sprintf("Unexpected Error Occurred : %s", err))
		return
	}


	if user.Email == "" && user.Phone == "" {

		hashedPassward, err := HashPassword(params.Passward)

		if err != nil{
			log.Println("Error Hashing Password:", err) 
			return 
		}
		
		new_user, err := apicfg.DB.CreateUser(r.Context(),database.CreateUserParams{
			UserID: 		uuid.New(),
			FirstName: 		params.First_name,
			LastName:  		sql.NullString{String: params.Last_name, Valid: params.Last_name != ""},
			Email:     		params.Email,
			Phone :  		params.Phone,
			Passward :    	hashedPassward,
			UserType:     	sql.NullString{String: params.User_type, Valid: params.User_type != ""},
			CreatedAt:    	sql.NullTime{Time :time.Now().UTC(), Valid: true},
			UpdatedAt:    	sql.NullTime{Time :time.Now().UTC(), Valid: true},
			
		})

		if err != nil{
			respondWithError(w,400, fmt.Sprintf("Error Creating user : %s", err))
			return
		}

		response := UserResponse{
		UserID:    new_user.UserID.String(),
		FirstName: new_user.FirstName,
		LastName:  new_user.LastName.String,
		Email:     new_user.Email,
		Phone:     new_user.Phone,
		UserType:  new_user.UserType.String,
		ApiKey:    new_user.ApiKey,
		CreatedAt: new_user.CreatedAt.Time,
		UpdatedAt: new_user.UpdatedAt.Time,
	}

		respondWithJSON(w,201,response)

	}else{
		respondWithError(w,409,"User Already Exists !")
	}
	
}

// only the admin can access the user details
func(apicfg *apiConfig) getUser(w http.ResponseWriter, r *http.Request, user database.User){

	if  user.UserType.Valid && user.UserType.String != "admin"{
		respondWithError(w,401,"Only admin can access this data !")
		return
	}

	type parameters struct{
		Email string `json:"email"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)

	
	if err != nil{
		respondWithError(w,400, fmt.Sprintf("Error Parsing Json12 : %s", err))
		return
	}

	if params.Email == "" {
		respondWithError(w, 400, "Email cannot be empty.")
		return
	}

	theUser,err := apicfg.DB.GetUserByEmail(r.Context(),params.Email)

	if err != nil {
		respondWithError(w, 404, "User not found.")
		return
	}

	response := UserResponse{
		UserID:    theUser.UserID.String(),
		FirstName: theUser.FirstName,
		LastName:  theUser.LastName.String,
		Email:     theUser.Email,
		Phone:     theUser.Phone,
		UserType:  theUser.UserType.String,
		ApiKey:    theUser.ApiKey,
		CreatedAt: theUser.CreatedAt.Time,
		UpdatedAt: theUser.UpdatedAt.Time,
	}

	respondWithJSON(w,200,response)


}