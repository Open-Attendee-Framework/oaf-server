package api100

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/concertLabs/oaf-server/internal/helpers"
	db100 "github.com/concertLabs/oaf-server/pkg/db/v100"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
)

const apiVersion = "1.0.0"

func getSubrouter(prefix string) http.Handler {

	a100 := mux.NewRouter().PathPrefix(prefix).Subrouter()
	a100 = a100.StrictSlash(true)
	a100.HandleFunc("/auth", authHandler).Methods("GET")
	a100.HandleFunc("/auth-refesh", authRefreshHandler).Methods("GET")
	a100.HandleFunc("/register", registerHandler).Methods("POST")

	chain := alice.New().Then(a100)

	return chain
}

func buildAPIResponse(r *http.Request, i interface{}) ([]byte, error) {
	ar := apiResponse{Version: apiVersion, Path: r.RequestURI, Data: i}
	j, err := json.Marshal(&ar)
	if err != nil {
		return j, errors.New("Error marshalling JSON:" + err.Error())
	}
	return j, nil
}

func writeJSONtoResponse(w http.ResponseWriter, j []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}

func writeStructtoResponse(w http.ResponseWriter, r *http.Request, i interface{}) {
	j, err := buildAPIResponse(r, i)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, errorJSONError)
		return
	}
	writeJSONtoResponse(w, j)
}

func apierror(w http.ResponseWriter, r *http.Request, err string, httpcode int, ecode apiErrorcode) {
	//Erzeugt einen json error Response und gibt ihn über http.Error zurück
	log.Println(err)
	er := errorResponse{strconv.Itoa(httpcode), ecode, ecode.String() + ":" + err}
	j, _ := buildAPIResponse(r, er)
	http.Error(w, string(j), httpcode)
}

func generateNewToken(un db100.User) (string, error) {
	mySigningKey := ""
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":       "oafserver",
		"exp":       time.Now().Add(time.Hour * 72).Unix(),
		"user":      un.UserID,
		"superuser": un.SuperUser,
	})
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(mySigningKey))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(tokenString)), err
}

func getTokenfromRequest(r *http.Request) (*jwt.Token, error) {
	mySigningKey := ""
	m, t, err := getAuthorization(r)
	if err != nil {
		return nil, err
	}

	if m != "Bearer" {
		return nil, errors.New("Bearer head missing")
	}

	data, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(string(data), func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil

	})

	return token, err
}

func getUserfromToken(token *jwt.Token) (db100.User, error) {
	un := db100.User{}

	claims := token.Claims.(jwt.MapClaims)
	ui, ok := claims["user"].(float64)
	if !ok {
		return un, errors.New("No id")
	}

	uid := int(ui)
	un.UserID = uid
	err := un.GetDetails()
	if err != nil {
		return un, errors.New("Error getting user details:" + err.Error())
	}

	return un, nil
}

func getAuthorization(r *http.Request) (string, string, error) {
	auth := r.Header.Get("Authorization")
	s := strings.Split(auth, " ")
	if len(s) < 2 {
		return "", "", errors.New("Authorization header malformed. Expected \"Bearer <token>\" got " + auth)
	}
	return s[0], s[1], nil
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	m, t, err := getAuthorization(r)
	if err != nil {
		apierror(w, r, err.Error(), 401, errorMalformedAuth)
		return
	}

	if strings.ToLower(m) != "bearer" {
		apierror(w, r, "Auth Request malformed", 401, errorMalformedAuth)
		return
	}

	data, err := base64.StdEncoding.DecodeString(t)
	if err != nil {
		apierror(w, r, err.Error(), 401, errorMalformedAuth)
		return
	}

	s := strings.Split(string(data), ":")
	if len(s) < 2 {
		apierror(w, r, "Auth Request malformed", 401, errorMalformedAuth)
		return
	}

	u, p := s[0], s[1]

	b, err := db100.DoesUserExist(u)
	if err != nil {
		apierror(w, r, err.Error(), 500, errorDBQueryFailed)
		return
	}
	if !b {
		apierror(w, r, "Wrong Username or Password", 401, errorWrongCredentials)
		return
	}

	un := db100.User{Username: u}
	err = un.GetDetailstoUsername()
	if err != nil {
		apierror(w, r, err.Error(), 500, errorDBQueryFailed)
		return
	}

	pw, err := helpers.GeneratePasswordHash(p, un.Salt)
	if err != nil {
		apierror(w, r, err.Error(), 500, errorNoHash)
		return
	}
	if pw != un.Password {
		apierror(w, r, "Wrong Username or Password", 401, errorWrongCredentials)
		return
	}

	tokenString, err := generateNewToken(un)
	if err != nil {
		apierror(w, r, err.Error(), 500, errorNoToken)
		return
	}
	ar := authResponse{tokenString}
	writeStructtoResponse(w, r, ar)
}

func authRefreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenfromRequest(r)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, errorMalformedAuth)
		return
	}

	un, err := getUserfromToken(token)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, errorMalformedAuth)
		return
	}

	tokenString, err := generateNewToken(un)
	if err != nil {
		apierror(w, r, err.Error(), 500, errorNoToken)
		return
	}

	ar := authResponse{tokenString}
	writeStructtoResponse(w, r, ar)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	createUser(w, r, true)
}
