package main

import (
	"log"

	httpauth "github.com/bihari123/golang-authentication/http-auth"
)

func main() {
	// jsonencoding.JsonEncodingDecoding()
//	httpauth.ConvertToBase64()
pass:="123456789"
hashedPasss,err:=httpauth.HashPassword(pass)
if err!=nil{
	panic(err)
}

err=httpauth.ComparePassword(pass,hashedPasss)
if err!=nil{
	log.Fatalln("Not logged in")
}

log.Println("Logged In")
}
