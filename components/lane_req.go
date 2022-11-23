package moi

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func NewLaneCExtractor() *ContentExtractor {
	extractor := ContentExtractor{}
	extractor.f1 = regexp.MustCompile(`{LANE_C:'.{0,10}'}[,?]*`)
	extractor.f2 = regexp.MustCompile(`LANE_C:'(.+)'`)
	return &extractor
}

type LaneCRequest struct {
	RequestBase
	taskName string
	townID   string
	road     string

	laneCExtractor *ContentExtractor
}

func (m *LaneCRequest) init() {

	if m.laneCExtractor == nil {
		m.laneCExtractor = NewLaneCExtractor()
	}

	if m.extractor == nil {
		m.extractor = NewExtractor()
	}
	if m.path == "" {
		m.path = "https://addressrs.moi.gov.tw/address/cfm/json.cfm"
	}

	m.refreshSession()
}

func (m *LaneCRequest) prepare() {
	data := FormDataFormat{TownID: m.townID, TaskName: m.taskName, Road: m.road}
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

func (m *LaneCRequest) fillLanes(context *Context, key string, data []byte) {
	if data != nil {
		temp := m.laneCExtractor.f1.FindAllStringSubmatch(string(data), -1)

		for i, v := range temp {
			if i > 0 {
				break
			}
			c := m.laneCExtractor.f2.FindStringSubmatch(v[0])

			if len(c) > 0 {
				context.laneCMap[key] = append(context.laneCMap[key], string(c[1]))
				// fmt.Println(string(c[1]))
			}

		}
	} else {
		context.laneCMap[key] = append(context.laneCMap[key], "無")
	}

}

func (m *LaneCRequest) Run(context *Context) {
	m.init()
	m.taskName = "LANE_C"
	var data []byte
	var ok bool = false
	for town, roads := range context.townAndRoad {

		for _, road := range roads {
			ok = false
			data = nil
			m.townID = town
			m.road = encodeURIComponent(road)

			for {
				m.prepare()
				data, ok = m.Do()
				if !ok {
					log.Println("laneC req error")
				} else {
					break
				}
			}

			fmt.Println(string(data))
			m.fillLanes(context, encodeURIComponent(town+road), data)
		}

	}
}

func (m *LaneCRequest) TestReq(townID string, road string) {
	m.init()
	m.taskName = "LANE_C"
	m.townID = townID
	m.road = encodeURIComponent(road)
	m.prepare()
	data, ok := m.Do()
	if !ok {
		log.Println("laneC req error")
		return
	}

	fmt.Println(string(data))
	temp := m.laneCExtractor.f1.FindAllStringSubmatch(string(data), -1)

	for i, v := range temp {
		if i > 0 {
			break
		}
		c := m.laneCExtractor.f2.FindStringSubmatch(v[0])

		if len(c) > 0 {
			fmt.Println(string(c[1]))
		}

	}

}
