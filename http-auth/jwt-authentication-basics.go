package httpauth

import (
	"fmt"
	"time"

  jwt	 "github.com/dgrijalva/jwt-go"
)

type UserClaims struct{
  jwt.StandardClaims
  SessionID int64 `json:"session_id"`
 }

func (u *UserClaims)Valid()error{
  if !u.VerifyExpiresAt(time.Now().Unix(),true){
    return fmt.Errorf("Token has expired")
  }
  if u.SessionID == 0{
    return fmt.Errorf("invalid session ID")
  }
  return nil 
}

func createToken(c *UserClaims)(string , error){
  t:= jwt.NewWithClaims(jwt.SigningMethodHS512,c)
    signedToken,err:=t.SignedString(key)

    if err!=nil{
    	return "",fmt.Errorf("Error in create Token when signing token: %w",err)
    }

    return signedToken,nil

}
func parseToken( signedToken string)(*UserClaims , error){
	claims:= &UserClaims{}
	t,err:= jwt.ParseWithClaims(signedToken,claims,func (t *jwt.Token)(interface{},error)  {

		if t.Method.Alg()!=jwt.SigningMethodHS512.Alg(){
			return nil, fmt.Errorf("Invalid signing algorithm")
		}
     return GenerateKey(key),nil 
	})

	if err!=nil{
		return nil,fmt.Errorf("Error in parse token while parsing token: %w",err)
	}
	
	if !t.Valid{
		return nil,fmt.Errorf("Error in parsing tokien, token invalid")

	}

	return t.Claims.(*UserClaims),nil  
}
