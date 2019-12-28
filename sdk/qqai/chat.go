package qqai

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const ChatUri = "/nlp/nlp_textchat"

type ChatRequest struct {
	Session  string
	Question string
}

type ChatResponse struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		Session string `json:"session"`
		Answer  string `json:"answer"`
	} `json:"data"`
}

func Chat(session string, question string) *ChatRequest {
	c := &ChatRequest{
		Session:  session,
		Question: question,
	}
	return c
}

func (q *ChatRequest) GetUrl() *url.URL {
	reqUrl, _ := url.Parse(fmt.Sprintf("%s%s", BaseUri, ChatUri))
	return reqUrl
}

func (q *ChatRequest) Do(c *Client, result interface{}) error {
	params := url.Values{}
	params.Set("session", q.Session)
	params.Set("question", q.Question)
	c.Sign(&params)
	response, err := c.Request(q.GetUrl().String(), params)
	if err != nil {
		return err
	}
	_ = json.Unmarshal(response, result)
	return nil
}
