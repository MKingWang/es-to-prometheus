package main

import (
	"fmt"
	"net/http"
	"time"
)

//启动http服务
func httpRun() {
	http.HandleFunc("/oom", oom)
	http.ListenAndServe(":8080", nil)
}

//oom页面的handle
func oom(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	oom := make(map[string]int)
	tmNow := time.Now().Unix()
	dayNow := time.Now().Format("2006.01.02")
	client := newClient("10.60.0.84", 9200, setAuth("wangleixin", "Compris2008"))
	client = client.EsIndex("logstash-" + dayNow).EsType("fluentd")
	msg := client.EsSearch("message", "Out of memory:")
	for _, v := range msg {
		if tmNow-v.Time.Unix() < 60 && tmNow-v.Time.Unix() > 0 {
			oom[v.Host]++

		}
	}

	oommsg := ""
	for k, v := range oom {
		oommsg = fmt.Sprintf("%soom{host=\"%s\"} %d\n", oommsg, k, v)
	}

	w.Write([]byte(oommsg))
}
