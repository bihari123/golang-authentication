package main

import "encoding/json"

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
}
