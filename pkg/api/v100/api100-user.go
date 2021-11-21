package api100

import (
	"encoding/json"
	"net/http"

	"github.com/concertLabs/oaf-server/internal/helpers"
	db100 "github.com/concertLabs/oaf-server/pkg/db/v100"
	"github.com/gorilla/mux"
)

func createUser(w http.ResponseWriter, r *http.Request, overrideSuperUser bool) {
	decoder := json.NewDecoder(r.Body)
	var u db100.User
	err := decoder.Decode(&u)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusBadRequest, errorJSONError)
		return
	}

	if u.Username == "" {
		apierror(w, r, "Username not set", http.StatusBadRequest, errorInvalidParameter)
		return
	}

	if u.Password == "" {
		apierror(w, r, "Password not set", http.StatusBadRequest, errorInvalidParameter)
		return
	}

	s, err := helpers.GenerateSalt()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, errorNoHash)
		return
	}
	u.Salt = s

	pw, err := helpers.GeneratePasswordHash(u.Password, u.Salt)
	if err != nil {
		apierror(w, r, err.Error(), 500, errorNoHash)
		return
	}
	u.Password = pw

	if overrideSuperUser {
		u.SuperUser = false
	}

	err = u.Insert()
	if err != nil {
		apierror(w, r, err.Error(), 500, errorDBQueryFailed)
		return
	}

	writeStructtoResponse(w, r, u)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	token, err := getTokenfromRequest(r)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, errorMalformedAuth)
		return
	}

	ou, err := getUserfromToken(token)
	if err != nil {
		apierror(w, r, "Auth Request malformed", 401, errorMalformedAuth)
		return
	}

	vars := mux.Vars(r)
	n := vars["name"]
	u := db100.User{Username: n}
	err = u.GetDetailstoUsername()
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, errorDBQueryFailed)
		return
	}

	if (ou.UserID != u.UserID) && (!ou.SuperUser) {
		apierror(w, r, "User not permitted for this Action", http.StatusUnauthorized, errorUserNotAuthorized)
		return
	}

	j, err := json.Marshal(&u)
	if err != nil {
		apierror(w, r, err.Error(), http.StatusInternalServerError, errorJSONError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(j)
}
