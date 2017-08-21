package main

import (
	"github.com/labstack/gommon/log"
	"net/http"

	"github.com/nirasan/go-oauth-pkce-code-verifier"
	"encoding/json"
)

func main() {
	http.HandleFunc("/auth/start", handlerAuthStart)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type AuthStartResponse struct {
	CodeChallenge string
}

func handlerAuthStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	cv, e := go_oauth_pkce_code_verifier.CreateCodeVerifier()
	if e != nil {
		panic(e)
		return
	}

	//TODO save code_verifier
	log.Info(cv.String())

	res := AuthStartResponse{CodeChallenge: cv.CodeChallengeS256()}
	json.NewEncoder(w).Encode(&res)
}
