package httpauth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	
)

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	var cred Credentials

	err := json.NewDecoder(r.Body).Decode(&cred)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword, ok := users[cred.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != cred.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time

	claims := &Claims{
		Username: cred.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// create the jwt string

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})
	//If a user logs in with the correct credentials, this handler will then set a cookie on the client side with the JWT value. Once a cookie
	// is set on a client, it is sent along with every request henceforth. Now we can write our welcome handler to handle user specific
	// information.

}

// Handling Post Authentication Routes

// Now that all logged in clients have session information stored on their end as cookies, we can use it to:

// • Authenticate subsequent user requests
// • Get information about the user making the request

// Let’s write our Welcome handler to do just that:

func Welcome(w http.ResponseWriter, r *http.Request) {
	// we can obtain the session token from the requests cookies, which come with every request

	c, err := r.Cookie("token")

	if err != nil {
		if err == http.ErrNoCookie {
			// if the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value

	// Initialize a new instance of "Claims"

	claims := &Claims{}

	// Parse the JWT string and store the result in the "claims",
	// NOTE that we are passing the key in this method as well. This method will return
	// an error if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusBadRequest)
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user, along with their username  given in the token

	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

//  Renewing Your Token

// In this example, we have set a short expiry time of five minutes. We should not expect the user to login every five minutes if their
// token expires. To solve this, we will create another /refresh route that takes the previous token (which is still valid), and returns a
// new token with a renewed expiry time.

// To minimize misuse of a JWT, the expiry time is usually kept in the order of a few minutes. Typically the client application would
// refresh the token in the background

func Refresh(w http.ResponseWriter, r *http.Request) {
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
// We ensure that a new token is not issued until enough time has elapsed
        // In this case, a new token will only be issued if the old token is within
        // 30 seconds of expiry. Otherwise, return a bad request status

				if time.Unix(claims.ExpiresAt,0).Sub(time.Now()) > 30*time.Second{
            w.WriteHeader(http.StatusBadRequest)
            return
				}

				// Now, create a new token for the current use, with a renewed expiration expirationTime
				expirationTime:= time.Now().Add(5*time.Minute)

				claims.ExpiresAt= expirationTime.Unix()

				token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

				tokenString,err:=token.SignedString(jwtKey)

				if err!=nil{
					w.WriteHeader(http.StatusInternalServerError)
					return 
				}

				// Set the new token as the users 'token' cookie

				http.SetCookie(w, &http.Cookie{
					Name: "token",
					Value: tokenString,
					Expires: expirationTime,
				})
}
