package main

import (
	moi "moi-addr/components"
)

func TestSearch() {
	moi.FormDataMapInit()
	// context := moi.NewContext()
	a := moi.ResultRequest{}

	// a.Run2("110", "■麻溝", "無", "無", "1000")
	a.Run2("010", "更生北路", "不限", "無", "1000")
	// a.Run2("010", "上海街", "不限", "無")
	// a.Run2("010", `%E4%B8%8A%E6%B5%B7%E8%A1%97`, `%E4%B8%8D%E9%99%90`, `%E7%84%A1`)

	// a.Run(context)

	// b := moi.RoadRequest{}
	// b.Run(context)
	// context.Print()

}

func TestLane() {
	moi.FormDataMapInit()
	// context := moi.NewContext()
	a := moi.LaneCRequest{}
	// a.Run2("110", "■麻溝")
	// a.Run2("010", "上海街")
	a.Run2("150", "台坂村８鄰啦里吧")
	// a.Run2("010", `%E4%B8%8A%E6%B5%B7%E8%A1%97`, `%E4%B8%8D%E9%99%90`, `%E7%84%A1`)

	// a.Run(context)

	// b := moi.RoadRequest{}
	// b.Run(context)
	// context.Print()

}

func TestLong() {
	moi.FormDataMapInit()
	// context := moi.NewContext()
	a := moi.LongCRequest{}
	// a.Run2("110", "■麻溝")
	// a.Run2("010", "上海街", "不限")
	a.Run2("150", "台坂村８鄰啦里吧", "無")
	// a.Run2("010", `%E4%B8%8A%E6%B5%B7%E8%A1%97`, `%E4%B8%8D%E9%99%90`, `%E7%84%A1`)

	// a.Run(context)

	// b := moi.RoadRequest{}
	// b.Run(context)
	// context.Print()

}
func TestResult() {
	moi.FormDataMapInit()

	// var a []byte = []byte(`{"results": 10,"rows": [{"ADDRESS_ALL":"池上鄉福文村002鄰鐵花路３５號","TM2X":"272542.37458506","TM2Y":"2558297.08495315"},{"ADDRESS_ALL":"池上鄉福文村002鄰鐵花路３７號","TM2X":"272628.72261714","TM2Y":"2558369.66479321"},{"ADDRESS_ALL":"池上鄉福文村002鄰鐵花路３７之１號","TM2X":"272654.36204117","TM2Y":"2558370.61813721"},{"ADDRESS_ALL":"池上鄉福文村003鄰鐵花路３號","TM2X":"272354.35205689","TM2Y":"2558151.59733701"},{"ADDRESS_ALL":"池上鄉福文村003鄰鐵花路１１號","TM2X":"272447.69068097","TM2Y":"2558204.25602506"},{"ADDRESS_ALL":"池上鄉福文村003鄰鐵花路１４號","TM2X":"272453.80","TM2Y":"2558201.11"},{"ADDRESS_ALL":"池上鄉福文村003鄰鐵花路１７號","TM2X":"272454.58911298","TM2Y":"2558229.19170508"},{"ADDRESS_ALL":"池上鄉福文村003鄰鐵花路２７號","TM2X":"272474.03256900","TM2Y":"2558200.91164106"}]}`)
	// total, _ := jsonparser.GetInt(a, "results")
	// fmt.Println(total)
	// jsonparser.ArrayEach(a, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 	jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
	// 		fmt.Printf("Key: '%s' ,Value: '%s', Type: %s\n", string(key), string(value), dataType)
	// 		return nil
	// 	})
	// 	// fmt.Println(string(value))
	// }, "rows")

	resultReq := moi.ResultRequest{}
	resultReq.Run()
	// resultReq.CrossValidation()
	resultReq.OutputFile()
	// context.Print2()
}

func TestAll() {
	moi.FormDataMapInit()
	context := moi.NewContext()
	townReq := moi.TownRequest{}
	townReq.Run(context)

	roadReq := moi.RoadRequest{}
	roadReq.Run(context)
	context.SaveTownAndRoad()

	laneReq := moi.LaneCRequest{}
	laneReq.Run(context)

	longReq := moi.LongCRequest{}
	longReq.Run(context)

	context.SaveLaneCMap()
}

func main() {

	// a.Prepare()
	// data, err := a.Do()
	// if err != nil {
	// 	log.Println(err)
	// }
	// fmt.Println(data)
	// data, _ := os.ReadFile("town.json")

	// fmt.Print(string(dat))
	// jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
	// 	fmt.Println("test")
	// 	jsonparser.ObjectEach(value, func(key []byte, value []byte, dataType jsonparser.ValueType, offset int) error {
	// 		fmt.Printf("Key: '%s'\n Value: '%s'\n Type: %s\n", string(key), string(value), dataType)
	// 		return nil
	// 	})

	// }, "root")
	// RTest()
	// AAAAA()
	// TestSearch()
	// TestLane()
	// TestLong()
	// TestLong()
	// TestAll()
	TestResult()
}
