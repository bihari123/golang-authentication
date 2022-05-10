package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type person struct {
	First string
}

func main() {
	p1 := person{
		First: "jenny",
	}
	p2 := person{
		First: "james",
	}

	xp := []person{p1, p2}

	bs,err:=json.Marshal(xp)

	if err!=nil{
		log.Panic(err)
	}
	xp2:=[]person{}
	err= json.Unmarshal(bs,&xp2)

	if err!=nil{
		fmt.Println(" back to golang structure",xp2) 
	}

	http.HandleFunc("/encode",foo)
	http.HandleFunc("/decode",bar)
	err = http.ListenAndServe(":8080",nil)

	if err != nil{
log.Panic("server did not set up")
	}
}

func foo(w http.ResponseWriter,r *http.Request){

}
func bar (w http.ResponseWriter,r *http.Request){

}
