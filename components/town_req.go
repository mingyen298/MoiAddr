package moi

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func NewTownExtractor() *ContentExtractor {
	extractor := ContentExtractor{}
	extractor.f1 = regexp.MustCompile(`{.{13},.{14,15}'}*`)
	extractor.f2 = regexp.MustCompile(`TOWN_ID:'(.+)',TOWN_NAME:'(.+)'`)
	return &extractor
}

type TownRequest struct {
	RequestBase
	taskName      string
	townExtractor *ContentExtractor
}

func (m *TownRequest) init() {

	if m.townExtractor == nil {
		m.townExtractor = NewTownExtractor()
	}
	// if m.header == nil {
	// 	m.header = make(http.Header)
	//
	// }

	if m.extractor == nil {
		m.extractor = NewExtractor()
	}
	if m.path == "" {
		m.path = "https://addressrs.moi.gov.tw/address/cfm/json.cfm"
	}

	m.refreshSession()
}
func (m *TownRequest) prepare() *http.Request {
	data := FormDataFormat{TaskName: m.taskName}
	req, err := m.makeRequest(&data)
	// req, err := http.NewRequest("POST", m.path, bytes.NewReader([]byte(m.body)))
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return req
}

func (m *TownRequest) fillTowns(context *Context, data []byte) {
	temp := m.townExtractor.f1.FindAllStringSubmatch(string(data), -1)
	for _, i := range temp {
		c := m.townExtractor.f2.FindStringSubmatch(i[0])
		// context.townAndRoad[c[1]] = []string{}
		context.townAndRoad.Add(c[1], []string{})
	}

}

func (m *TownRequest) Run(context *Context) {
	m.init()
	m.taskName = "TOWN_ID"
	var data []byte = nil
	var ok bool = false
	var req *http.Request
	for {
		req = m.prepare()
		data, ok = m.Do(req)
		if !ok {
			log.Println("road req error")
		} else {
			break
		}
	}

	m.fillTowns(context, data)
}

func (m *TownRequest) TestReq() {
	m.init()
	m.taskName = "TOWN_ID"
	var data []byte = nil
	var ok bool = false
	var req *http.Request
	for {
		req = m.prepare()
		data, ok = m.Do(req)
		if !ok {
			log.Println("road req error")
		} else {
			break
		}
	}
	fmt.Println(string(data))

}
