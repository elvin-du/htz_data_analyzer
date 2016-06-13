package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"htz_data_analyzer/log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	All = map[string]*ParsedData{}
	Raw = []byte{}
)

var (
	jzavg       JZAvgSlice
	jztotal     JZTotalSlice
	zzavg       ZZAvgSlice
	zztotal     ZZTotalSlice
	timeMorning TimeMorningSlice
	timeNight   TimeNightSlice
)

var (
	regNameTime   = regexp.MustCompile(`(.*)(\s\s)(\d\d:\d\d)`)
	regMorning    = regexp.MustCompile(`早安`)
	regZhanZhuang = regexp.MustCompile(`站桩(\d+)\+(\d+)`)
	regJingZuo    = regexp.MustCompile(`静坐(\d+)\+(\d+)`)
	klmCount      = regexp.MustCompile(`宽两秒(\d+)`)
)

func Parse() error {
	f, err := os.Open("data.txt")
	if nil != err {
		log.Debugln(err)
		return err
	}
	defer f.Close()

	name, t := "", ""
	scanner := bufio.NewScanner(f)
	flag := false
	for scanner.Scan() {
		text := scanner.Text()
		data := regNameTime.FindStringSubmatch(text)
		if nil != data {
			var d *ParsedData
			ok := false
			name = data[1]
			if d, ok = All[name]; ok {
				d.Total += 1
			} else {
				d = &ParsedData{}
				d.Name = name
				d.Total = 1
				All[name] = d
			}
			t = data[3]
			flag = false
		} else {
			//早安打卡
			data = regMorning.FindStringSubmatch(text)
			if nil != data {
				All[name].MorningCount += 1
				All[name].MorningTime = append(All[name].MorningTime, t)
				continue
			}

			//站桩
			data = regZhanZhuang.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}

				All[name].ZhanZhuangCount += 2
				duration, _ := strconv.Atoi(data[1])
				duration2, _ := strconv.Atoi(data[2])
				All[name].ZhanZhuangDuration += duration + duration2
				All[name].ZhanZhuangAvg = All[name].ZhanZhuangDuration / All[name].ZhanZhuangCount
				goto JINGZUO
			}
		}
	JINGZUO:
		//静坐
		data = regJingZuo.FindStringSubmatch(text)
		if nil != data {
			if !flag {
				flag = true
				All[name].ReportTime = append(All[name].ReportTime, t)
			}
			All[name].JingZuoCount += 2
			duration, _ := strconv.Atoi(data[1])
			duration2, _ := strconv.Atoi(data[2])
			All[name].JingZuoDuration += duration + duration2
			All[name].JingZuoAvg = All[name].JingZuoDuration / All[name].JingZuoCount
			continue
		}

		All[name].ReportCount = len(All[name].ReportTime)
	}

	for _, v := range All {
		total := 0
		for _, v2 := range v.MorningTime {
			v3 := strings.Split(v2, ":")
			hour, err := strconv.Atoi(v3[0])
			if nil != err {
				log.Debugln(err)
				return err
			}

			sec, err := strconv.Atoi(v3[1])
			if nil != err {
				log.Debugln(err)
				return err
			}

			total += hour*60 + sec
		}
		if 0 == v.MorningCount {
			v.MorningTimeAvgStr = "无"
			v.MorningTimeAvg = 10000 * 10000
		} else {
			total = total / v.MorningCount
			v.MorningTimeAvg = total

			hour := total / 60
			sec := total % 60
			v.MorningTimeAvgStr = fmt.Sprintf("%02d:%02d", hour, sec)
		}

		//汇报时间
		for _, v2 := range v.ReportTime {
			v3 := strings.Split(v2, ":")
			hour, err := strconv.Atoi(v3[0])
			if nil != err {
				log.Debugln(err)
				return err
			}

			sec, err := strconv.Atoi(v3[1])
			if nil != err {
				log.Debugln(err)
				return err
			}

			total += hour*60 + sec
		}
		if 0 == v.ReportCount {
			v.ReportTimeAvgStr = "无"
			v.ReportTimeAvg = 10000 * 10000
		} else {
			total = total / v.ReportCount
			v.ReportTimeAvg = total

			hour := total / 60
			sec := total % 60

			v.ReportTimeAvgStr = fmt.Sprintf("%02d:%02d", hour, sec)
		}
	}

	for _, v := range All {
		jzavg = append(jzavg, v)
		jztotal = append(jztotal, v)
		zzavg = append(zzavg, v)
		zztotal = append(zztotal, v)
		timeMorning = append(timeMorning, v)
		timeNight = append(timeNight, v)
	}

	sort.Stable(jzavg)
	sort.Stable(jztotal)
	sort.Stable(zzavg)
	sort.Stable(zztotal)
	sort.Stable(timeMorning)
	sort.Stable(timeNight)

	f2, err := os.Create("sorted.txt")
	if nil != err {
		log.Debugln(f2)
	}
	defer f2.Close()

	for _, v := range timeMorning {
		bin, err := json.MarshalIndent(v, "", "    ")
		if nil != err {
			log.Debugln(err)
			return err
		}

		Raw = append(Raw, bin...)

		_, err = f2.Write(bin)
		if nil != err {
			log.Debugln(err)
			return err
		}
	}

	return nil
}
