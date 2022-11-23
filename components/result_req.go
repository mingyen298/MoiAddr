package moi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/buger/jsonparser"
)

func NewResultExtractor() *ContentExtractor {
	extractor := ContentExtractor{}
	extractor.f1 = regexp.MustCompile(`{ROADSEC:.{3,30},ROADSEC_NV:.{3,30}'}[,?]*`)
	extractor.f2 = regexp.MustCompile(`ROADSEC:'(.+)',ROADSEC_NV:'(.+)'`)
	return &extractor
}

type ResultRequest struct {
	RequestBase
	taskName        string
	townID          string
	road            string
	laneC           string
	longC           string
	no              string
	addrSrchType    string
	completeAddr    string
	start           string
	limit           string
	resultExtractor *ContentExtractor
	townAndRoadMap  map[string][]string
	laneAndLongMap  map[string][]string
	csvList         []string
}

func (m *ResultRequest) init() {

	if m.resultExtractor == nil {
		m.resultExtractor = NewResultExtractor()
	}

	if m.extractor == nil {
		m.extractor = NewExtractor()
	}
	if m.path == "" {
		m.path = "https://addressrs.moi.gov.tw/address/cfm/JSON.cfm"
	}

	if m.townAndRoadMap == nil {
		townroadData, err := os.ReadFile(path.Join("file", "townroad.json"))
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(townroadData, &m.townAndRoadMap)
		if err != nil {
			log.Fatal(err)
		}
	}
	if m.laneAndLongMap == nil {
		lanelongData, err := os.ReadFile(path.Join("file", "lanelong.json"))
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(lanelongData, &m.laneAndLongMap)
		if err != nil {
			log.Fatal(err)
		}
	}
	m.csvList = []string{}
	m.refreshSession()
}

func (m *ResultRequest) prepare() *http.Request {
	// data := FormDataFormat{TownID: m.townID, TaskName: m.taskName, Road: m.road}
	data := FormDataFormat{}
	data.TaskName = m.taskName
	data.TownID = m.townID
	data.Road = m.road
	data.LaneC = m.laneC
	data.LongC = m.longC
	data.NO = m.no
	data.AddrSrchType = m.addrSrchType
	data.CompeleteAddr = m.completeAddr
	data.Start = m.start
	data.Limit = m.limit

	req, err := m.makeRequest(&data)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return req
}

func (m *ResultRequest) CrossValidation() {
	m.init()
	for townID, roads := range m.townAndRoadMap {
		for _, road := range roads {
			_, isExist := m.laneAndLongMap[encodeURIComponent(townID+road)]
			if !isExist {
				log.Printf("err addr:%s , fix :%s", townID+road, encodeURIComponent(townID+road))
			}
		}
	}
}

func (m *ResultRequest) Run() {
	m.init()
	m.taskName = "SRCH_ADDR"
	m.no = "!"
	m.start = "0"
	m.completeAddr = "!"
	m.addrSrchType = "2"
	var data []byte
	var ok bool = false
	var req *http.Request
	var tempList []string = make([]string, 0)
	for town, roads := range m.townAndRoadMap {

		for _, road := range roads {
			data = nil
			ok = false
			m.townID = encodeURIComponent(town)
			m.road = encodeURIComponent(road)
			if len(m.laneAndLongMap[encodeURIComponent(town+road)]) < 1 {
				log.Fatalf("outofrange == %s", town+road)
			}
			m.limit = "1000"
			m.laneC = encodeURIComponent(m.laneAndLongMap[encodeURIComponent(town+road)][0])
			m.longC = encodeURIComponent(m.laneAndLongMap[encodeURIComponent(town+road)][1])

			for {
				req = m.prepare()
				data, ok = m.Do(req)
				if !ok {
					log.Println("result req error")
				} else {
					break
				}
				total, _ := jsonparser.GetInt(data, "results")
				limitNum, _ := strconv.Atoi(m.limit)
				if int(total) > limitNum {
					m.limit = strconv.Itoa(int(total))
					continue
				}
				fmt.Println(string(data))
				break

			}

			jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
					tempList = append(tempList, string(value))
					return nil
				})
				// tempList = append(tempList, "\n")
				out := strings.Join(tempList, ",") + "\n"
				tempList = tempList[:0]
				m.csvList = append(m.csvList, out)
			}, "rows")

			// m.fillLanes(context, town, data)
			// time.Sleep(time.Millisecond * 100)
		}

	}
}

func (m *ResultRequest) OutputFile() {
	f, err := os.Create(path.Join("file", "result.csv"))

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	f.WriteString(strings.Join(m.csvList, ""))
}

func (m *ResultRequest) TestReq(townID string, road string, lane string, long string, limit string) {
	m.init()
	m.taskName = "SRCH_ADDR"
	m.no = "!"
	m.limit = limit
	m.start = "0"
	m.completeAddr = "!"
	m.addrSrchType = "2"

	m.townID = townID
	m.road = encodeURIComponent(road)
	m.laneC = encodeURIComponent(lane)
	m.longC = encodeURIComponent(long)
	var data []byte = nil
	var ok bool = false
	var req *http.Request
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

}
