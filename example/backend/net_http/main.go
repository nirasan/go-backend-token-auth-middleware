package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/labstack/gommon/log"
	"net/http"
	"net/url"
	"os"
	"net/http/httputil"
	"github.com/nirasan/go-backend-token-auth-middleware/token"
	"io/ioutil"
	"github.com/dgrijalva/jwt-go"
	"github.com/nirasan/go-backend-token-auth-middleware/middleware"
	"github.com/rs/cors"
)

const (
	CodeChallengeContextKey = "CodeChallenge"
)

func main() {
	prepareAuthenticationMiddleware := middleware.PrepareAuthenticationMiddleware(middleware.PrepareAuthenticationMiddlewareOption{
		Next: handlerAuthStart,
		CodeVerifierSetter: func(cv, cc string) { codeVerifierMap[cc] = cv },
		CodeChallengeContextKey: CodeChallengeContextKey,
	})
	http.HandleFunc("/auth/start", prepareAuthenticationMiddleware)
	http.HandleFunc("/auth", handlerAuth)
	http.HandleFunc("/userinfo", handlerUserinfo)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
		AllowCredentials: true,
		AllowedHeaders: []string{"Authorization"},
	})

	log.Fatal(http.ListenAndServe(":8080", c.Handler(http.DefaultServeMux)))
}

//TODO use any data store
var codeVerifierMap map[string]string = make(map[string]string)

func handlerAuthStart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	cc := r.Context().Value(CodeChallengeContextKey).(string)

	res := struct {
		CodeChallenge string `json:"code_challenge"`
	}{
		CodeChallenge: cc,
	}
	json.NewEncoder(w).Encode(&res)
}

//TODO use any data store
var userMap map[string]*User = make(map[string]*User)

type User struct {
	ID string
	Name string
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
		fmt.Printf("code_verifier not found: %+v, %+v", req, codeVerifierMap)
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
	v.Add("code_verifier", cv)

	resp1, err := http.PostForm("https://www.googleapis.com/oauth2/v4/token", v)
	if err != nil {
		fmt.Printf("get token error: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	defer resp1.Body.Close()

	dumpResponse(resp1)

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

	dumpResponse(resp2)

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

	/**
	 * check user existence
	 */
	isNew := false
	u, ok := userMap[body2.Sub]
	if !ok {
		isNew = true
		u = &User{ID: body2.Sub, Name: body2.Name}
		userMap[u.ID] = u
	}

	/**
	 * create token
	 */
	m, err := createTokenManager()
	if err != nil {
		fmt.Printf("create token manager error: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}
	token, err := m.CreateSignedToken(m.CreateToken(u.ID))
	if err != nil {
		fmt.Printf("create token manager error: %v", err)
		http.Error(w, err.Error(), 400)
		return
	}

	res := struct {
		Name string `json:"name"`
		AccessToken string `json:"access_token"`
		IsNew bool `json:"is_new"`
	}{
		Name: u.Name,
		AccessToken: token,
		IsNew: isNew,
	}
	json.NewEncoder(w).Encode(&res)
}

func handlerUserinfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Printf("%+v\n", r)

	m, err := createTokenManager()
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	t, err := m.ParseTokenFromRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if claims, ok := t.Claims.(jwt.MapClaims); !ok || !t.Valid {
		http.Error(w, "invalid token", 403)
		return
	} else if id, ok := claims["sub"].(string); !ok {
		http.Error(w, "invalid token", 403)
		return
	} else if u, ok := userMap[id]; !ok {
		http.Error(w, "user not found", 403)
		return
	} else {
		// success
		res := struct {
			Name string `json:"name"`
		}{
			Name: u.Name,
		}
		json.NewEncoder(w).Encode(&res)
	}
}

func dumpResponse(r *http.Response) {
	dump2, err := httputil.DumpResponse(r, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n\n", dump2)
}

func createTokenManager() (*token.TokenManager, error) {
	return token.CreateTokenManager(token.CreateTokenManagerOption{
		SigningAlgorithm: "ES256",
		PrivateKeyLoader: func() interface{} {
			keyData, e := ioutil.ReadFile("./ec256-private.pem")
			if e != nil {
				panic(e.Error())
			}
			key, e := jwt.ParseECPrivateKeyFromPEM(keyData)
			if e != nil {
				panic(e.Error())
			}
			return key
		},
		PublicKeyLoader: func() interface{} {
			keyData, e := ioutil.ReadFile("./ec256-public.pem")
			if e != nil {
				panic(e.Error())
			}
			key, e := jwt.ParseECPublicKeyFromPEM(keyData)
			if e != nil {
				panic(e.Error())
			}
			return key
		},
	})
}
