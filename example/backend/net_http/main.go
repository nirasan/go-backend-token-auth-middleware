package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	codeverifier "github.com/nirasan/go-oauth-pkce-code-verifier"
	"net/http"
	"net/url"
	"os"
	"net/http/httputil"
)

func main() {
	http.HandleFunc("/auth/start", handlerAuthStart)
	http.HandleFunc("/auth", handlerAuth)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

//TODO データベースなどに記録する
var codeVerifierMap map[string]*codeverifier.CodeVerifier = make(map[string]*codeverifier.CodeVerifier)

func handlerAuthStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	cv, e := codeverifier.CreateCodeVerifier()
	if e != nil {
		panic(e)
		return
	}

	cc := cv.CodeChallengeS256()
	codeVerifierMap[cc] = cv

	res := struct {
		CodeChallenge string `json:"code_challenge"`
	}{
		CodeChallenge: cc,
	}
	json.NewEncoder(w).Encode(&res)
}

func handlerAuth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	/**
	 * get params
	 */
	req := struct {
		Code          string `json:"code"`
		CodeChallenge string `json:"code_challenge"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Printf("decode error 1: \n%v, \n%+v, \nbody: %+v", err, r, r.Body)
		http.Error(w, err.Error(), 400)
		return
	}

	/**
	 * get code verifier
	 */
	cv, ok := codeVerifierMap[req.CodeChallenge]
	if !ok {
		fmt.Printf("code_verifier not found: %v", req)
		http.Error(w, "code verifier not found", 400)
		return
	}

	/**
	 * get access token
	 */
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		http.Error(w, err.Error(), 400)
		return
	}
	v := url.Values{}
	v.Add("code", req.Code)
	v.Add("client_id", os.Getenv("CLIENT_ID"))
	v.Add("client_secret", os.Getenv("CLIENT_SECRET"))
	v.Add("redirect_uri", "http://localhost:4200/callback")
	v.Add("grant_type", "authorization_code")
	//v.Add("access_type", "offline")
	v.Add("code_verifier", cv.Value)

	resp1, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", v)
	if err != nil {
		fmt.Printf("get token error: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer resp1.Body.Close()

	dump1, err := httputil.DumpResponse(resp1, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n\n", dump1)

	body := struct {
		AccessToken string `json:"access_token"`
	}{}
	err = json.NewDecoder(resp1.Body).Decode(&body)
	if err != nil {
		fmt.Printf("decode error 2: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Printf("%q\n\n", body)

	/**
	 * get user info
	 */
	v = url.Values{}
	v.Add("access_token", body.AccessToken)
	resp2, err := http.Get("https://www.googleapis.com/oauth2/v3/userinfo" + "?" + v.Encode())
	if err != nil {
		fmt.Printf("get userinfo error: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer resp2.Body.Close()

	dump2, err := httputil.DumpResponse(resp2, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n\n", dump2)

	body2 := struct {
		Sub string `json:"sub"`
		Name string `json:"name"`
	}{}
	err = json.NewDecoder(resp2.Body).Decode(&body2)
	if err != nil {
		fmt.Printf("decode error: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Printf("userinfo: %+v", body2)

	//TODO サーバー発行のトークンを返す
	res := struct {
		Success bool `json:"success"`
	}{
		Success: true,
	}
	json.NewEncoder(w).Encode(&res)
}
