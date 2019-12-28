package vision

import (
	"encoding/json"
	"fmt"
	"github.com/shybily/go-utils/sdk/qqai"
	"net/url"
)

type TextRequest struct {
	Image string
}

type TextResponse struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		Text string `json:"text"`
	} `json:"data"`
}

const ToTextUri = "/vision/vision_imgtotext"

func ToText(image string) *TextRequest {
	t := &TextRequest{
		Image: image,
	}
	return t
}

func (t *TextRequest) Do(c *qqai.Client, result interface{}) error {
	values := url.Values{}
	values.Set("vision", t.Image)
	values.Set("session_id", qqai.RandStringBytes(16))
	c.Sign(&values)
	response, err := c.Request(t.GetUrl().String(), values)
	if err != nil {
		return err
	}
	err = json.Unmarshal(response, result)
	return err
}

func (t *TextRequest) GetUrl() *url.URL {
	reqUrl, _ := url.Parse(fmt.Sprintf("%s%s", qqai.BaseUri, ToTextUri))
	return reqUrl
}
