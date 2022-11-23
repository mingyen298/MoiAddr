package moi

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

type RoadExtractor struct {
	f1 *regexp.Regexp
	f2 *regexp.Regexp
}

func NewRoadExtractor() *RoadExtractor {
	extractor := RoadExtractor{}
	extractor.f1 = regexp.MustCompile(`{ROADSEC:.{3,30},ROADSEC_NV:.{3,30}'}[,?]*`)
	extractor.f2 = regexp.MustCompile(`ROADSEC:'(.+)',ROADSEC_NV:'(.+)'`)
	return &extractor
}

type RoadRequest struct {
	RequestBase
	taskName string
	townID   string
	kW       string

	roadExtractor *RoadExtractor
}

func (m *RoadRequest) init() {

	if m.roadExtractor == nil {
		m.roadExtractor = NewRoadExtractor()
	}

	if m.extractor == nil {
		m.extractor = NewExtractor()
	}
	if m.path == "" {
		m.path = "https://addressrs.moi.gov.tw/address/cfm/json.cfm"
	}

	m.refreshSession()
}

func (m *RoadRequest) prepare() {
	data := FormDataFormat{TownID: m.townID, TaskName: m.taskName, KW: m.kW}
	m.bind(&data)

	req, err := http.NewRequest("POST", m.path, bytes.NewReader([]byte(m.body)))
	if err != nil {
		log.Println(err.Error())
		return
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Cookie", m.session)
	m.req = req
	log.Printf("body:%s\n", m.body)
	log.Printf("session:%s\n", m.session)
}

func (m *RoadRequest) fillRoads(context *Context, key string, data []byte) {
	temp := m.roadExtractor.f1.FindAllStringSubmatch(string(data), -1)

	for _, i := range temp {
		c := m.roadExtractor.f2.FindStringSubmatch(i[0])
		context.townAndRoad[key] = append(context.townAndRoad[key], c[1])

	}

}

func (m *RoadRequest) Run(context *Context) {
	m.init()
	m.taskName = "ROADSEC"
	m.kW = "0"

	for town, _ := range context.townAndRoad {

		m.townID = town
		var data []byte = nil
		var ok bool = false
		for {
			m.prepare()
			data, ok = m.Do()
			if !ok {
				log.Println("road req error")
			} else {
				break
			}
		}

		m.fillRoads(context, town, data)
		// time.Sleep(time.Millisecond * 100)
	}
}

func (m *RoadRequest) TestReq(townID string) {
	m.init()
	m.taskName = "ROADSEC"
	m.kW = "0"
	m.townID = townID
	var data []byte = nil
	var ok bool = false
	for {
		m.prepare()
		data, ok = m.Do()
		if !ok {
			log.Println("road req error")
		} else {
			break
		}
	}
	fmt.Println(string(data))

}
