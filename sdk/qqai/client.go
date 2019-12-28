package qqai

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const BaseUri = "https://api.ai.qq.com/fcgi-bin"
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type ClientInterface interface {
	Request(path string, query url.Values) ([]byte, error)
	Sign(values *url.Values)
}

type Client struct {
	AppId  string
	AppKey string
}

type Request interface {
	GetUrl() *url.URL
	Do(c *Client, result interface{}) error
}

func NewClient(appId string, appKey string) *Client {
	c := &Client{
		AppId:  appId,
		AppKey: appKey,
	}
	return c
}

func (c *Client) Request(reqUrl string, query url.Values) ([]byte, error) {
	resp, err := http.PostForm(reqUrl, query)
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}
	return body, nil
}

func (c *Client) Sign(values *url.Values) {
	values.Set("app_id", c.AppId)
	values.Set("time_stamp", strconv.FormatInt(time.Now().Unix(), 10))
	values.Set("nonce_str", RandStringBytes(8))
	str := fmt.Sprintf("%s&app_key=%s", values.Encode(), url.QueryEscape(c.AppKey))
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str)))
	values.Set("sign", strings.ToUpper(sign))
}

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}