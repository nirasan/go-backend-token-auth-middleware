package middleware

import (
	"net/http"
	codeverifier "github.com/nirasan/go-oauth-pkce-code-verifier"
	"context"
)

type PrepareAuthenticationMiddlewareOption struct {
	Next               http.HandlerFunc
	CodeVerifierSetter func(codeVerifier, codeChallenge string)
	CodeChallengeContextKey string
}

func PrepareAuthenticationMiddleware(o PrepareAuthenticationMiddlewareOption) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cv, e := codeverifier.CreateCodeVerifier()
		if e != nil {
			panic(e)
			return
		}

		cc := cv.CodeChallengeS256()
		o.CodeVerifierSetter(cv.Value, cc)

		r = r.WithContext(context.WithValue(r.Context(), o.CodeChallengeContextKey, cc))

		o.Next(w, r)
	}
}
