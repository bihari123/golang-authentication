package jsonencoding

import (
	"encoding/json"
	"log"
	"net/http"
)

type person struct {
	First string
}

func JsonEncodingDecoding() {
	// p1 := person{
	// 	First: "jenny",
	// }
	// p2 := person{
	// 	First: "james",
	// }
	//
	// xp := []person{p1, p2}
	//
	// bs, err := json.Marshal(xp)
	//
	// if err != nil {
	// 	log.Panic(err)
	// }
	// xp2 := []person{}
	// err = json.Unmarshal(bs, &xp2)
	//
	// if err != nil {
	// 	fmt.Println(" back to golang structure", xp2)
	// }
	//
	http.HandleFunc("/encode", foo)
	http.HandleFunc("/decode", bar)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Panic("server did not set up")
	}
}

func foo(w http.ResponseWriter, r *http.Request) {

	p1 := person{
		First: "jenny",
	}

	err:=json.NewEncoder(w).Encode(p1)

	if err !=nil{
		log.Println("Encoded bad data",err)
	}
}
func bar(w http.ResponseWriter, r *http.Request) {
 var p1 person 

 err:=json.NewDecoder(r.Body).Decode(&p1)
if err!=nil{
	log.Println("Decode bad data", err)
}
log.Println("Person: ",p1) 
}
