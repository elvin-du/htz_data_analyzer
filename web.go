package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	TPL *template.Template
)

func init() {
	http.HandleFunc("/jz/total", JZTotalHandler)
	http.HandleFunc("/jz/avg", JZAvgHandler)
	http.HandleFunc("/zz/total", ZZTotalHandler)
	http.HandleFunc("/zz/avg", ZZAvgHandler)
	http.HandleFunc("/time/morning", TimeMorningHandler)
	http.HandleFunc("/time/night", TimeNightHandler)
	http.HandleFunc("/data/upload", UploadDataHandler)
	//	http.HandleFunc("/raw", RawHandler)
	http.Handle("/", http.StripPrefix("/static", http.FileServer(http.Dir("./pub"))))

	//	var err error = nil
	//	TPL = template.New("htz_data_analyzer")
	//	TPL = TPL.Funcs(template.FuncMap{"Plus": plus})
	//	TPL, err = TPL.ParseGlob("./pub/*.html")
	//	if nil != err {
	//		log.Fatal(err)
	//	}
}

func Tmp() {
	var err error = nil
	TPL = template.New("htz_data_analyzer")
	TPL = TPL.Funcs(template.FuncMap{"Plus": plus})
	TPL, err = TPL.ParseGlob("./pub/*.html")
	if nil != err {
		log.Fatal(err)
	}
}

func plus(args ...interface{}) string {
	if len(args) <= 0 {
		return "args length <= 0"
	}

	index, ok := args[0].(int)
	if !ok {
		return "args type is not int"
	}

	return strconv.Itoa(index + 1)
}

func JZTotalHandler(w http.ResponseWriter, r *http.Request) {
	Tmp()
	data := struct {
		Title string
		Items []*ParsedData
	}{
		Title: "静坐总时间排序结果",
		Items: jztotal,
	}

	err := TPL.ExecuteTemplate(w, "jztotal.html", data)
	if nil != err {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func JZAvgHandler(w http.ResponseWriter, r *http.Request) {
	Tmp()
	data := struct {
		Title string
		Items []*ParsedData
	}{
		Title: "静坐平均时间排序结果",
		Items: jzavg,
	}

	err := TPL.ExecuteTemplate(w, "jzavg.html", data)
	if nil != err {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func ZZTotalHandler(w http.ResponseWriter, r *http.Request) {
	Tmp()
	data := struct {
		Title string
		Items []*ParsedData
	}{
		Title: "站桩总时间排序结果",
		Items: zztotal,
	}

	err := TPL.ExecuteTemplate(w, "zztotal.html", data)
	if nil != err {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func ZZAvgHandler(w http.ResponseWriter, r *http.Request) {
	Tmp()
	data := struct {
		Title string
		Items []*ParsedData
	}{
		Title: "站桩平均时间排序结果",
		Items: zzavg,
	}

	err := TPL.ExecuteTemplate(w, "zzavg.html", data)
	if nil != err {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func TimeNightHandler(w http.ResponseWriter, r *http.Request) {
	Tmp()
	data := struct {
		Title string
		Items []*ParsedData
	}{
		Title: "早睡时间排序结果",
		Items: timeNight,
	}

	err := TPL.ExecuteTemplate(w, "time_night.html", data)
	if nil != err {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func TimeMorningHandler(w http.ResponseWriter, r *http.Request) {
	Tmp()
	data := struct {
		Title string
		Items []*ParsedData
	}{
		Title: "早起时间排序结果",
		Items: timeMorning,
	}

	err := TPL.ExecuteTemplate(w, "time_morning.html", data)
	if nil != err {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func RawHandler(w http.ResponseWriter, r *http.Request) {
	Tmp()
	data := struct {
		Title string
		Items []*ParsedData
	}{
		Title: "原始数据",
		Items: timeMorning,
	}

	err := TPL.ExecuteTemplate(w, "raw.html", data)
	if nil != err {
		log.Println(err)
		w.Write([]byte(err.Error()))
	}
}

func UploadDataHandler(w http.ResponseWriter, r *http.Request) {
	Tmp()
	if r.Method == "GET" {
		data := struct {
			Title string
		}{
			Title: "上传文件",
		}

		err := TPL.ExecuteTemplate(w, "upload.html", data)
		if nil != err {
			log.Println(err)
			w.Write([]byte(err.Error()))
		}
		return
	}

	err := r.ParseMultipartForm(32 << 20)
	if nil != err {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	f, fh, err := r.FormFile("data")
	if nil != err {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	defer f.Close()

	log.Println("upload file mime:", fh.Header.Get("Content-Type"))
	log.Println("filename:", fh.Filename)
	bin, err := ioutil.ReadAll(f)
	if nil != err {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	err = ioutil.WriteFile("data.txt", bin, os.ModePerm)
	if nil != err {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	err = Parse()
	if nil != err {
		log.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}

	http.Redirect(w, r, "/jz/total", http.StatusSeeOther)
}

func Run(addr string) {
	if err := http.ListenAndServe(addr, nil); nil != err {
		log.Fatal(err)
	}
}
