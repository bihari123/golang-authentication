package httpauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

jwt 	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

/*
create a struct to hold your Claims. This determines how your decoded tokens would be structured. The *jwt.RegisteredClaims (earlier we had used StandardClaims) include the basic claims such as exp, sub,jti,aud,etc imported from the JWT package. The UserInfo interface{} is whatever you include inside your token. This is usually a bunch of key:value pair.

- Never populate the struct with any sensitive information such as the user's password or secret keys.
*/
type MyJWTClaims struct {
	*jwt.RegisteredClaims
	UserInfo interface{}
}

// Generate your own secret key !

var secret = []byte("can-you-keep-a-secret?")

func CreateToken(sub string, userInfo interface{}) (string, error) {
	// Get the token instance with the signing method

	token := jwt.New(jwt.GetSigningMethod("HS256"))

	// choos an expiration time , shorter the better

	exp := time.Now().Add(time.Hour * 24)

	// Add your claims

	token.Claims = &MyJWTClaims{
		&jwt.RegisteredClaims{
			// set the exp and sub claims. sub is usually the userID
			ExpiresAt: jwt.NewNumericDate(exp),
			Subject:   sub,
		},
		userInfo,
	}

	// Sign the token with your secret key

	val, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	return val, nil
}

//The next snippet is for verifying the token and retrienving the claims from the JWT string
func GetClaimsFromtoken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

//Now lets make a middleware to verify the JWT on a protected resource.

func AuthenticationMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth_header := r.Header.Get("Authorization")

		if !strings.HasPrefix(auth_header, "Bearer") {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(auth_header, "Bearer ")

		claims, err := GetClaimsFromtoken(tokenString)

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
//This adds the decoded claims (sub, userId, username, full_name, etc) to the response object so that you can have access to it later in the request handler function where you are protecting your resources. To keep things clean we will create the SetJWTClaimsContext and JWTClaimsFromContext functions ourselves.
		r = r.WithContext(SetJWTClaimsContext(r.Context(), claims))

		next(w, r)
	}
}


type ClaimsKey int 
var claimsKey ClaimsKey

func SetJWTClaimsContext(ctx context.Context, claims jwt.MapClaims)context.Context{
	return context.WithValue(ctx,claimsKey,claims)
}

func JWTClaimsFromContext(ctx context.Context)(jwt.MapClaims, bool){
	claims,ok:= ctx.Value(claimsKey).(jwt.MapClaims)

	return claims,ok 
}


// Below is how you access the claims that was added to the http.request object in the middleware. The returned claim is a map of key:value pairs to help you identify the client/user.

/* 
  protectedHandler := func(w http.ResponseWriter, r *http.Request){
    claims, ok := JWTClaimsFromContext(r.Context())
// Rest of the code
...
}*/  



func main() {

	loginHandler := func(w http.ResponseWriter, r *http.Request){
    
    form := make(map[string]string)

    username, ok := form["username"]
    if !ok {
      http.Error(w, "No username field", http.StatusBadRequest)
      return
    }
    password, ok := form["password"]
    if !ok {
      http.Error(w, "No password field", http.StatusBadRequest)
      return
    }

    // Create user if your conditions match. Below, all username and passwords are accepted.
    user := &User{
      Id: "2332-abcd-acdf-ccd2",
      Name: "JWT Master",
      Username: username,
      Password: password,
    }

    tokenString, _ := CreateToken(user.Id, user)
    payload := make(map[string]string)
    payload["access_token"] = tokenString

    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    json.NewEncoder(w).Encode(payload)
  }

  protectedHandler := func(w http.ResponseWriter, r *http.Request){
    claims, ok := JWTClaimsFromContext(r.Context())
    
    //Do something with the UserInfo claims
    if val, ok := claims["UserInfo"]; ok {
      userinfo := val.(map[string]string)
      fmt.Print(userinfo)  
    }
    //Do something with the sub claim
    if val, ok := claims["sub"]; ok {
      userid := val.(string)
      fmt.Print(userid)  
    }

    if !ok {
      http.Error(w, "Something went wrong", http.StatusInternalServerError)
      return
    }
    json.NewEncoder(w).Encode(claims)
  }

  indexHandler := func(w http.ResponseWriter, r *http.Request){
    fmt.Fprint(w, "Status OK")
  }

  m := mux.NewRouter()
  m.HandleFunc("/", indexHandler).Methods("GET")
  m.HandleFunc("/login", loginHandler).Methods("POST")

	protected := m.PathPrefix("/").Subrouter()
  protected.Use(AuthenticationMW)
  protected.HandleFunc("/resource",protectedHandler).Methods("GET","POST")

  log.Fatal(http.ListenAndServe(":8080", m))
}
 
