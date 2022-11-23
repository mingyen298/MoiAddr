package moi

import (
	"net/url"
	"strings"
)

var FormDataMap map[string]string

type TownAndRoad map[string][]string

func FormDataMapInit() {
	FormDataMap = make(map[string]string)
	FormDataMap["TaskName"] = "task_name"
	FormDataMap["TownID"] = "TOWN_ID"
	FormDataMap["KW"] = "KW"
	FormDataMap["Limit"] = "limit"
	FormDataMap["AddrSrchType"] = "addr_srch_type"
	FormDataMap["Road"] = "ROAD"
	FormDataMap["LaneC"] = "LANE_C"
	FormDataMap["LongC"] = "LONG_C"
	FormDataMap["NO"] = "NO"
	FormDataMap["CompeleteAddr"] = "COMPLETE_ADDR"
	FormDataMap["Start"] = "start"

}

func encodeURIComponent(str string) string {
	r := url.QueryEscape(str)
	r = strings.Replace(r, "+", "%20", -1)
	return r
}

// type FormDataFormat struct {
// 	TownID        string
// 	TaskName      string
// 	KW            string
// 	Road          string
// 	NO            string
// 	AddrSrchType  string
// 	CompeleteAddr string
// 	Start         string
// 	Limit         string
// LaneC string
// LongC string
// }
