package httpauth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
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
