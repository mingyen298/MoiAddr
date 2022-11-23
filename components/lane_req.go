package moi

import (
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

func (m *LaneCRequest) prepare() *http.Request {
	data := FormDataFormat{TownID: m.townID, TaskName: m.taskName, Road: m.road}
	req, err := m.makeRequest(&data)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return req
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
				context.laneCMap.Append(key, string(c[1]))
				// context.laneCMap[key] = append(context.laneCMap[key], string(c[1]))
				// fmt.Println(string(c[1]))
			}

		}
	} else {
		// context.laneCMap[key] = append(context.laneCMap[key], "無")
		context.laneCMap.Append(key, "無")
	}

}

func (m *LaneCRequest) Run(context *Context) {
	m.init()
	m.taskName = "LANE_C"
	var data []byte
	var ok bool = false
	var req *http.Request
	for town, roads := range context.townAndRoad.All() {

		for _, road := range roads {
			ok = false
			data = nil
			m.townID = town
			m.road = encodeURIComponent(road)

			for {
				req = m.prepare()
				data, ok = m.Do(req)
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
	req := m.prepare()
	data, ok := m.Do(req)
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
