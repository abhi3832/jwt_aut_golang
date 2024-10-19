package auth

import (
	"errors"
	"net/http"
	"strings"
)

// AUTHRIZATION : token {INSERT YOUR token HERE}
func Get_token(headers http.Header) (string, error){

	val := headers.Get("Authorization")

	if val == ""{
		return "", errors.New("no Auth info Found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2{
		return "", errors.New("malformed Auth Header")
	}

	if vals[0] != "token"{
		return "", errors.New("malformed first part of Auth Header")
	}

	return vals[1],nil

}

// AUTHRIZATION : userid {INSERT YOUR id HERE}
func Get_Id(headers http.Header) (string, error){

	val := headers.Get("Authorization")

	if val == ""{
		return "", errors.New("no Auth info Found")
	}

	vals := strings.Split(val, " ")

	if len(vals) != 2{
		return "", errors.New("malformed Auth Header")
	}

	if vals[0] != "userid"{
		return "", errors.New("malformed first part of Auth Header")
	}

	return vals[1],nil

}