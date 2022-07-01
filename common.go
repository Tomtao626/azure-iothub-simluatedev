package main

import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func getHttpClient() (client *http.Client) {
	transport := http.Transport{
		Dial:              dialTimeout,
		DisableKeepAlives: true,
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: &transport,
	}
	return
}

func dialTimeout(network, addr string) (net.Conn, error) {
	c, err := net.DialTimeout(network, addr, time.Second*50) //设置建立连接超时
	if err != nil {
		return nil, err
	}
	err = c.SetDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		return nil, err
	} //设置发送接收数据超时
	return c, nil
}

func HttpRequest(method, url string, body []byte, header map[string]string) (resp *http.Response, err error) {
	newReq, _ := http.NewRequest(method, url, bytes.NewReader(body))

	for k, v := range header {
		newReq.Header.Set(k, v)
	}

	client := getHttpClient()
	resp, err = client.Do(newReq)
	if err != nil {
		Error("new http req err.", err.Error())
		return
	}
	return
}

func ParseResponseString(response *http.Response) (string, error) {
	//var result map[string]interface{}
	body, err := ioutil.ReadAll(response.Body) // response.Body 是一个数据流
	return string(body), err                   // 将 io数据流转换为string类型返回！
}
