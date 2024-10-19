package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"time"

	auth "github.com/abhi3832/jwt_auth/Auth"
	"github.com/abhi3832/jwt_auth/internal/database"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var secretKey = []byte("i_am_batman")

func GenerateAccessToken(userId uuid.UUID) (string,error){

	claims := jwt.MapClaims{
		"sub": userId.String(),
		"exp": time.Now().Add(15*time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secretKey)
}

func GenerateRefreshToken() (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(24 * time.Hour).Unix(), // Set expiration time for the refresh token
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func validateJWTToken(tokenString string) (jwt.MapClaims, error) {
	
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("couldn't parse token to claim")
	}

	return claims, nil
}

func(apicfg *apiConfig) refreshTokens(w http.ResponseWriter, r *http.Request){

	userIDstr, err := auth.Get_Id(r.Header)

	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Error getting UserId %v", err))
		return
	}

	userId,err:= uuid.Parse(userIDstr)

	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Error Parsing UserId %v", err))
		return
	}

	user, err := apicfg.DB.GetUserByUserId(r.Context(),userId)

	if err != nil{
		respondWithError(w,400,fmt.Sprintf("User Not Found %v", err))
		return
	}

	session,err := apicfg.DB.GetSessionByUserId(r.Context(),userId)

	if err != nil{
		respondWithError(w,403,"Unauthorized, No session Found. Login to Create a Session!")
		return
	}

	// if the refresh token is also expired, the user needs to login again

	expireTime := session.ExpiresAt

	currentTime := time.Now()

	// Compare current time with expiresAt
	if currentTime.After(expireTime) {
		respondWithError(w,403,"Unauthorized, Token Expired!! You Have been logged out!!")
		return
	} 

	errr := apicfg.DB.UpdateSessionByDelete(r.Context(),session.UserID)

	if errr != nil{
		respondWithError(w,403,"Could't Delete Your Session")
		return
	}

	accessToken, err := GenerateAccessToken(userId)
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

		// Create a session record
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
}

