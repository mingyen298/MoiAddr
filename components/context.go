package moi

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

type Context struct {
	townAndRoad map[string][]string
	laneCMap    map[string][]string
}

func NewContext() *Context {
	return &Context{townAndRoad: make(map[string][]string), laneCMap: make(map[string][]string)}
}

func (m *Context) SaveTownAndRoad() {
	data, err := json.Marshal(m.townAndRoad)
	if err != nil {
		log.Println("covert to json err")
		return
	}
	f, err := os.Create(path.Join("file", "townroad.json"))

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	f.WriteString(string(data))

}

func (m *Context) SaveLaneCMap() {
	// fmt.Println(m.townAndRoad)
	data, err := json.Marshal(m.laneCMap)
	if err != nil {
		log.Println("covert to json err")
		return
	}

	f, err := os.Create(path.Join("file", "lanelong.json"))

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	f.WriteString(string(data))

}
