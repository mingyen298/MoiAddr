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

type RequestBase struct {
	extractor *Extractor
	path      string
	body      string
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

func (m *RequestBase) makeRequest(s *FormDataFormat) (*http.Request, error) {
	m.bind(s)
	req, err := http.NewRequest("POST", m.path, bytes.NewReader([]byte(m.body)))
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", m.session)
	req.Header.Add("Origin", "https://addressrs.moi.gov.tw")
	req.Header.Add("Host", "addressrs.moi.gov.tw")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")
	req.Header.Add("Accept-Language", "zh-TW,zh-Hant;q=0.9")
	// req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("User-Agent", ":	Mozilla/5.0 (iPhone; CPU iPhone OS 15_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/107.0.5304.101 Mobile/15E148 Safari/604.1")
	req.Header.Add("Referer", "https://addressrs.moi.gov.tw/address/index.cfm?city_id=10014")

	log.Printf("body:%s\n", m.body)
	log.Printf("session:%s\n", m.session)
	return req, nil
}

func (m *RequestBase) Do(req *http.Request) ([]byte, bool) {
	res, _ := http.DefaultClient.Do(req)
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
