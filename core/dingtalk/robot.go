package dingtalk

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tal-tech/go-zero/core/dingtalk/message"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	ResponseOk = 0
)

type Response struct {
	ErrCode int32  `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func GetMessageSign(timestamp int64, secret string) string {
	origin := fmt.Sprintf("%d\n%s", timestamp, secret)
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(origin))
	return url.QueryEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
}

func SendRobotMessage(url, secret string, msg message.Message) error {
	timestamp := time.Now().Unix() * 1000
	var reqUrl string
	if secret != "" {
		reqUrl = url + "&timestamp=" + strconv.FormatInt(timestamp, 10) + "&sign=" + GetMessageSign(timestamp, secret)
	} else {
		reqUrl = url + "&timestamp=" + strconv.FormatInt(timestamp, 10)
	}
	req, err := http.NewRequest(http.MethodPost, reqUrl, bytes.NewReader(msg.ToJson()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	if result.ErrCode != ResponseOk {
		return errors.New(result.ErrMsg)
	}
	return nil
}
