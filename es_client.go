package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type EsMsg struct {
	Took int   `json:"took"`
	Hits EsHit `json:"hits"`
}
type EsHit struct {
	Total    int      `json:"total"`
	MaxScore float64  `json:"max_score"`
	Hits     []ResHit `json:"hits"`
}

type ResHit struct {
	Index  string   `json:"_index"`
	Type   string   `json:"_type"`
	ID     string   `json:"_id"`
	Score  float64  `json:"_score"`
	Source EsSource `json:"_source"`
}

type EsSource struct {
	Host     string    `json:"host"`
	Ident    string    `json:"ident"`
	Message  string    `json:"message"`
	Env      string    `json:"env"`
	Priority string    `json:prrority`
	Time     time.Time `json:"@timestamp"`
}

type EsClient struct {
	Host   string
	Port   int
	Auth   string
	Url    string
	Index  string
	Type   string
	client *http.Client
}

//设置帐号密码
func setAuth(username, password string) string {
	return fmt.Sprintf("%s:%s@", username, password)
}

// 创建es客户端
func newClient(host string, port int, auth ...string) *EsClient {
	return &EsClient{
		Host:   host,
		Port:   port,
		Auth:   auth[0],
		client: &http.Client{},
	}
}

//设置index
func (c *EsClient) EsIndex(index string) *EsClient {
	c.Index = index
	return c
}

//设置type
func (c *EsClient) EsType(Type string) *EsClient {
	c.Type = Type

	return c
}

/*
查找数据功能待完善
目前只支持查询单个字段
*/

//查找数据
func (c *EsClient) EsSearch(k, v string) []EsSource {
	// http://user:password@es_host:9200/logstash-2018.08.31/fluentd/_search?q=message:Out
	url := "http://" + c.Auth + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.Index + "/" + c.Type + "/_search"
	var esmsg EsMsg
	var esSre []EsSource
	query := MakeQuery(k, v)
	esrequest, err := http.NewRequest("GET", url, query)
	Err(err)
	esrequest.Header.Set("Content-Type", "application/json")

	esRes, err := c.client.Do(esrequest)
	Err(err)

	buf, err := ioutil.ReadAll(esRes.Body)
	Err(err)

	json.Unmarshal(buf, &esmsg)
	for _, v := range esmsg.Hits.Hits {
		esSre = append(esSre, v.Source)
	}

	return esSre

}

func MakeQuery(k, v string) *bytes.Buffer {
	type querysearch struct {
		Query struct {
			Match interface{} `json:"match"`
		} `json:"query"`
	}
	match := make(map[string]string, 0)
	match[k] = v
	search := querysearch{}
	search.Query.Match = match
	encode, err := json.Marshal(search)
	Err(err)

	return bytes.NewBuffer(encode)
}

//分解出message信息
func (e *EsSource) GetMsg() {
	//
}

//错误处理
func Err(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

//嗯 。。。
/*
func searchDoc() {
	var esmsg EsMsg
	esrequest, err := http.NewRequest("GET", "http://user:password@es_host:9200/logstash-2018.08.31/fluentd/_search?q=message:Out", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

	c := http.Client{}
	esRes, err := c.Do(esrequest)
	if err != nil {
		fmt.Println(err.Error())
	}

	buf, err := ioutil.ReadAll(esRes.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	json.Unmarshal(buf, &esmsg)
	for _, v := range esmsg.Hits.Hits {
		fmt.Println(v.Source.Message, v.Source.Time.Unix())
	}
}
*/
