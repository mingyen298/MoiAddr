package moi

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
)

type FormDataFormat struct {
	TownID        string
	TaskName      string
	KW            string
	Road          string
	NO            string
	AddrSrchType  string
	CompeleteAddr string
	Start         string
	Limit         string
	LaneC         string
	LongC         string
}

// type RequestFormatIF interface {
// 	Init()
// 	Bind()
// }

type RequestBase struct {
	extractor *Extractor
	path      string
	body      string
	req       *http.Request
	header    http.Header
	session   string
}

func (m *RequestBase) bind(s *FormDataFormat) {

	param := reflect.ValueOf(s).Elem()
	temp := []string{}
	for i := 0; i < param.NumField(); i++ {
		if param.Field(i).Addr().Elem().String() != "" {
			if param.Field(i).Addr().Elem().String() == "!" {
				temp = append(temp, FormDataMap[param.Type().Field(i).Name]+`=`)
			} else {
				temp = append(temp, FormDataMap[param.Type().Field(i).Name]+`=`+param.Field(i).Addr().Elem().String())
			}
		}
	}
	m.body = strings.Join(temp, "&")
	fmt.Println(m.body)
}

func (m *RequestBase) Do() ([]byte, bool) {
	res, _ := http.DefaultClient.Do(m.req)
	if res.StatusCode > 200 {
		return nil, false
	}
	for _, val := range res.Cookies() {
		if val.Name == `JSESSIONID` {
			m.session = val.Name + `=` + val.Value
			break
		}
	}
	defer res.Body.Close()
	content, _ := io.ReadAll(res.Body)

	washed := m.extractor.GetJson(content)

	return washed, true
}

// func (m *RequestBase) Do2() ([]byte, bool) {
// 	res, _ := http.DefaultClient.Do(m.req)
// 	if res.StatusCode > 200 {
// 		return nil, false
// 	}
// 	for _, val := range res.Cookies() {
// 		if val.Name == `JSESSIONID` {
// 			m.session = val.Name + `=` + val.Value
// 			break
// 		}
// 	}
// 	defer res.Body.Close()
// 	content, _ := io.ReadAll(res.Body)

// 	washed := m.extractor.GetJson(content)

// 	return washed, true

// }

func (m *RequestBase) refreshSession() {
	req, err := http.NewRequest("GET", "https://addressrs.moi.gov.tw/address/index.cfm?city_id=10014", bytes.NewReader([]byte("")))
	if err != nil {
		log.Println(err.Error())
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	for _, val := range res.Cookies() {
		if val.Name == `JSESSIONID` {
			m.session = val.Name + `=` + val.Value
			break
		}
	}

}
