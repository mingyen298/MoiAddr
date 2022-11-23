package main

import (
	"fmt"
	moi "moi-addr/components"
)

func main() {
	a := moi.NewSafeMap()
	// a.Add("A", []string{"1", "2"})
	b := a.Get("A")
	fmt.Println(b)
	a.Append("A", "3")
	c := a.Get("A")
	fmt.Println(c)
	// b := a.All()
	// fmt.Println(b["A"])
	// b["A"] = append(b["A"], "3")
	// c := a.All()
	// fmt.Println(c["A"])
}
