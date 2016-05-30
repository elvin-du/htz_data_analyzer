package main

type rankItem struct {
	Rank  int
	Name  string
	Value int
}

//type jzavg struct {
//	Rank  int
//	Name  string
//	Value int
//}

//type zztotal struct {
//	Rank  int
//	Name  string
//	Value int
//}

//type zzavg struct {
//	Rank  int
//	Name  string
//	Value int
//}

type ParsedData struct {
	Name               string   `json:"姓名"`
	Total              int      `json:"共发言次数"` //总共发言次数
	MorningCount       int      `json:"早安次数"`
	MorningTime        []string `json:"早安时间"` //几点早安
	MorningTimeAvg     int      `json:"平均早安时间"`
	MorningTimeAvgStr  string   `json:"平均早安时间-字符串"`
	ReportCount        int      `json:"回报次数"`
	ReportTime         []string `json:"回报时间"`
	ReportTimeAvg      int      `json:"平均回报时间"`
	ReportTimeAvgStr   string   `json:"平均回报时间-字符串"`
	ZhanZhuangCount    int      `json:"站桩次数"`
	ZhanZhuangDuration int      `json:"站桩总时长"`
	ZhanZhuangAvg      int      `json:"站桩平均时间"`
	JingZuoCount       int      `json:"静坐次数"`
	JingZuoDuration    int      `json:"静坐总时长"`
	JingZuoAvg         int      `json:"静坐平均时间"`
//	KLMCount           int      `json:"宽两秒次数"`
}

type ZZTotalSlice []*ParsedData

func (p ZZTotalSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ZZTotalSlice) Less(i, j int) bool {
	return p[i].ZhanZhuangDuration > p[j].ZhanZhuangDuration
}

func (p ZZTotalSlice) Len() int {
	return len(p)
}

type ZZAvgSlice []*ParsedData

func (p ZZAvgSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ZZAvgSlice) Less(i, j int) bool {
	return p[i].ZhanZhuangAvg > p[j].ZhanZhuangAvg
}

func (p ZZAvgSlice) Len() int {
	return len(p)
}

type JZAvgSlice []*ParsedData

func (p JZAvgSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p JZAvgSlice) Less(i, j int) bool {
	return p[i].JingZuoAvg > p[j].JingZuoAvg
}

func (p JZAvgSlice) Len() int {
	return len(p)
}

type JZTotalSlice []*ParsedData

func (p JZTotalSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p JZTotalSlice) Less(i, j int) bool {
	return p[i].JingZuoDuration > p[j].JingZuoDuration
}

func (p JZTotalSlice) Len() int {
	return len(p)
}

type TimeMorningSlice []*ParsedData

func (p TimeMorningSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p TimeMorningSlice) Less(i, j int) bool {
	return p[i].MorningTimeAvg < p[j].MorningTimeAvg
}

func (p TimeMorningSlice) Len() int {
	return len(p)
}

type TimeNightSlice []*ParsedData

func (p TimeNightSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p TimeNightSlice) Less(i, j int) bool {
	return p[i].ReportTimeAvg < p[j].ReportTimeAvg
}

func (p TimeNightSlice) Len() int {
	return len(p)
}
