package httpauth

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
)
var key= []byte{}

// generating a 64 bytes key. 
// why? see the function signMessage 
func GenerateKey(key []byte)[]byte{
  for i:=1;i<64;i++{
    key=append(key,byte(i))
  }
  return key 
}
func signMessage(msg []byte)([]byte,error){
  h:= hmac.New(sha512.New(),GenerateKey(key)) // key is something that you generate yourself and will use for all messages with this  hmac:- to create that cryptographic signature and validate that signature. So, HMAC is only good for sending messages to yourself bcoz you need the same key. Also, the key size should match the size of your hashing algorithm. A sha512 is of the size 64 bytes. Hence you need a key of size 64 bytes.
  _,err:=h.Write(msg)

  if err!=nil{
    return nil, fmt.Errorf("Error in signing message: %w",err)
  }

  signature:=h.Sum(nil)
  return signature,nil 
}

// you generate the message and send it to the user
// the user sends the message back to you
// then you compare the two messages usiong this func 
func checkSig(msg, sig []byte)(bool,error){
// you first sign the message 

newSig,err:=signMessage(msg)
if err!=nil{
    return false,fmt.Errorf("Error in checkSign while signing message: %w",err)
}
// then you comapre it using hmac.Equal
same:=hmac.Equal(newSig,sig) //  it will return a boolean
return same , nil 
}


