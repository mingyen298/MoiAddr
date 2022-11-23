package moi

import (
	"encoding/json"
	"log"
	"os"
	"path"
	"sync"
)

type SafeMap struct {
	dict map[string][]string
	lock sync.RWMutex
}

func NewSafeMap() *SafeMap {
	return &SafeMap{dict: map[string][]string{}, lock: sync.RWMutex{}}
}

func (m *SafeMap) Add(key string, value []string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.dict[key] = value
}
func (m *SafeMap) Append(key string, value string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.dict[key] = append(m.dict[key], value)
}

func (m *SafeMap) Get(key string) []string {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return m.dict[key]
}

func (m *SafeMap) All() map[string][]string {
	return m.dict
}

type Context struct {
	// townAndRoad map[string][]string
	// laneCMap    map[string][]string
	townAndRoad *SafeMap
	laneCMap    *SafeMap
}

func NewContext() *Context {
	return &Context{townAndRoad: NewSafeMap(), laneCMap: NewSafeMap()}
}

func (m *Context) SaveTownAndRoad() {
	data, err := json.Marshal(m.townAndRoad.All())
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
	data, err := json.Marshal(m.laneCMap.All())
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
