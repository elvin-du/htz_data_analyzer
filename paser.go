package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	//	"io/ioutil"
	"log"
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
	regNameTime    = regexp.MustCompile(`(.*)(\s\s)(\d\d:\d\d)`)
	regMorning     = regexp.MustCompile(`早安`)
	regZhanZhuang  = regexp.MustCompile(`站桩(\d+)\+(\d+)`)
	regZhanZhuang2 = regexp.MustCompile(`站椿(\d+)\+(\d+)`)
	regZhanZhuang3 = regexp.MustCompile(`站桩(\d+)`)
	regZhanZhuang4 = regexp.MustCompile(`站椿(\d+)`)
	regZhanZhuang5 = regexp.MustCompile(`站桩:(\d+)`)
	regZhanZhuang6 = regexp.MustCompile(`站桩 (\d+)`)

	regJingZuo  = regexp.MustCompile(`静坐(\d+)\+(\d+)`)
	regJingZuo2 = regexp.MustCompile(`靜坐(\d+)\+(\d+)`)
	regJingZuo3 = regexp.MustCompile(`打坐(\d+)\+(\d+)`)
	regJingZuo4 = regexp.MustCompile(`静坐(\d+)`)
	regJingZuo5 = regexp.MustCompile(`靜坐(\d+)`)
	regJingZuo6 = regexp.MustCompile(`打坐(\d+)`)
	regJingZuo7 = regexp.MustCompile(`禅坐:(\d+)`)
	regJingZuo8 = regexp.MustCompile(`靜坐 (\d+)`)
	regJingZuo9 = regexp.MustCompile(`静坐 (\d+)`)

	klmCount  = regexp.MustCompile(`宽(\d+)`)
	klmCount2 = regexp.MustCompile(`宽两秒(\d+)`)
	klmCount3 = regexp.MustCompile(`寛兩秒(\d+)`)
)

func Parse() error {
	f, err := os.Open("data.txt")
	if nil != err {
		log.Println(err)
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

			data = regZhanZhuang2.FindStringSubmatch(text)
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

			data = regZhanZhuang3.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].ZhanZhuangCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].ZhanZhuangDuration += duration
				All[name].ZhanZhuangAvg = All[name].ZhanZhuangDuration / All[name].ZhanZhuangCount
				goto JINGZUO
			}

			data = regZhanZhuang4.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].ZhanZhuangCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].ZhanZhuangDuration += duration
				All[name].ZhanZhuangAvg = All[name].ZhanZhuangDuration / All[name].ZhanZhuangCount
				goto JINGZUO
			}

			data = regZhanZhuang5.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].ZhanZhuangCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].ZhanZhuangDuration += duration
				All[name].ZhanZhuangAvg = All[name].ZhanZhuangDuration / All[name].ZhanZhuangCount
				goto JINGZUO
			}

			data = regZhanZhuang6.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].ZhanZhuangCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].ZhanZhuangDuration += duration
				All[name].ZhanZhuangAvg = All[name].ZhanZhuangDuration / All[name].ZhanZhuangCount
				goto JINGZUO
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

			data = regJingZuo2.FindStringSubmatch(text)
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

			data = regJingZuo3.FindStringSubmatch(text)
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

			data = regJingZuo4.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].JingZuoCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].JingZuoDuration += duration
				All[name].JingZuoAvg = All[name].JingZuoDuration / All[name].JingZuoCount
				continue
			}

			data = regJingZuo5.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].JingZuoCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].JingZuoDuration += duration
				All[name].JingZuoAvg = All[name].JingZuoDuration / All[name].JingZuoCount
				continue
			}

			data = regJingZuo6.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].JingZuoCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].JingZuoDuration += duration
				All[name].JingZuoAvg = All[name].JingZuoDuration / All[name].JingZuoCount
				continue
			}

			data = regJingZuo7.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].JingZuoCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].JingZuoDuration += duration
				All[name].JingZuoAvg = All[name].JingZuoDuration / All[name].JingZuoCount
				continue
			}

			data = regJingZuo8.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].JingZuoCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].JingZuoDuration += duration
				All[name].JingZuoAvg = All[name].JingZuoDuration / All[name].JingZuoCount
				continue
			}

			data = regJingZuo9.FindStringSubmatch(text)
			if nil != data {
				if !flag {
					flag = true
					All[name].ReportTime = append(All[name].ReportTime, t)
				}
				All[name].JingZuoCount += 1
				duration, _ := strconv.Atoi(data[1])
				All[name].JingZuoDuration += duration
				All[name].JingZuoAvg = All[name].JingZuoDuration / All[name].JingZuoCount
				continue
			}

			//			data = klmCount.FindStringSubmatch(text)
			//			if nil != data {
			//				if !flag {
			//					flag = true
			//					All[name].ReportTime = append(All[name].ReportTime, t)
			//				}
			//				All[name].JingZuoCount += 1
			//				duration, _ := strconv.Atoi(data[1])
			//				All[name].JingZuoDuration += duration
			//				All[name].JingZuoAvg = All[name].JingZuoDuration / All[name].JingZuoCount
			//				continue
			//			}

			All[name].ReportCount = len(All[name].ReportTime)
		}
	}

	for _, v := range All {
		total := 0
		for _, v2 := range v.MorningTime {
			v3 := strings.Split(v2, ":")
			hour, err := strconv.Atoi(v3[0])
			if nil != err {
				log.Println(err)
				return err
			}

			sec, err := strconv.Atoi(v3[1])
			if nil != err {
				log.Println(err)
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
				log.Println(err)
				return err
			}

			sec, err := strconv.Atoi(v3[1])
			if nil != err {
				log.Println(err)
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
		log.Println(f2)
	}
	defer f2.Close()

	for _, v := range timeMorning {
		bin, err := json.MarshalIndent(v, "", "    ")
		if nil != err {
			log.Println(err)
			return err
		}

		Raw = append(Raw, bin...)

		_, err = f2.Write(bin)
		if nil != err {
			log.Println(err)
			return err
		}
	}

	return nil
}
