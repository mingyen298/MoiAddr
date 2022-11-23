package moi

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func NewLongCExtractor() *ContentExtractor {
	extractor := ContentExtractor{}
	extractor.f1 = regexp.MustCompile(`{LONG_C:'.{0,10}'}[,?]*`)
	extractor.f2 = regexp.MustCompile(`LONG_C:'(.+)'`)
	return &extractor
}

type LongCRequest struct {
	RequestBase
	taskName string
	townID   string
	road     string
	lane     string

	longCExtractor *ContentExtractor
}

func (m *LongCRequest) init() {

	if m.longCExtractor == nil {
		m.longCExtractor = NewLongCExtractor()
	}

	if m.extractor == nil {
		m.extractor = NewExtractor()
	}
	if m.path == "" {
		m.path = "https://addressrs.moi.gov.tw/address/cfm/JSON.cfm"
	}

	m.refreshSession()
}

func (m *LongCRequest) prepare() {
	data := FormDataFormat{TownID: m.townID, TaskName: m.taskName, Road: m.road, LaneC: m.lane}
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

func (m *LongCRequest) fillLongs(context *Context, key string, data []byte) {
	if data != nil {
		temp := m.longCExtractor.f1.FindAllStringSubmatch(string(data), -1)

		for i, v := range temp {
			if i > 0 {
				break
			}
			c := m.longCExtractor.f2.FindStringSubmatch(v[0])

			if len(c) > 0 {
				context.laneCMap[key] = append(context.laneCMap[key], c[1])

				// fmt.Println(string(c[1]))
			}

		}
	} else {
		context.laneCMap[key] = append(context.laneCMap[key], "無")
	}

}

func (m *LongCRequest) Run(context *Context) {
	m.init()
	m.taskName = "LONG_C"
	var data []byte = nil
	var ok bool = false
	for town, roads := range context.townAndRoad {

		for _, road := range roads {
			ok = false
			data = nil
			m.townID = town
			m.road = encodeURIComponent(road)
			m.lane = encodeURIComponent(context.laneCMap[encodeURIComponent(town+road)][0])

			for {
				m.prepare()
				data, ok = m.Do()
				if !ok {
					log.Println("longC req error")
				} else {
					break
				}
			}

			fmt.Println(string(data))
			m.fillLongs(context, encodeURIComponent(town+road), data)

		}

	}
}

func (m *LongCRequest) TestReq(townID string, road string, lane string) {
	m.init()
	m.taskName = "LONG_C"
	m.townID = townID
	m.road = encodeURIComponent(road)
	m.lane = encodeURIComponent(lane)

	var data []byte = nil
	var ok bool = false
	for {
		m.prepare()
		data, ok = m.Do()
		if !ok {
			log.Println("longC req error")
		} else {
			break
		}
	}

	temp := m.longCExtractor.f1.FindAllStringSubmatch(string(data), -1)
	for _, i := range temp {

		c := m.longCExtractor.f2.FindStringSubmatch(i[0])

		if len(c) > 0 {
			fmt.Println(c[1])
		}

	}

}
