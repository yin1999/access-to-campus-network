package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	ErrNeedLogin = errors.New("login needed")

	commonHeader = http.Header{
		"Accept": []string{"*/*"},
	}
)

type service string

const (
	outCampus service = "校园外网服务(out-campus NET)"
	cmcc      service = "中国移动(CMCC NET)"
)

type client struct {
	http.Client
	href string
}

func newClient() *client {
	return &client{
		Client: http.Client{
			Timeout: 15 * time.Second,
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				IdleConnTimeout: 10 * time.Second,
			},
		},
	}
}

func (c *client) check() error {
	req, err := newRequestWithCommonHeader(http.MethodGet, "http://www.google.cn/generate_204", nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer drainBody(resp.Body)
	if resp.StatusCode == http.StatusNoContent { // internet access
		return nil
	}
	data := make([]byte, 1000)
	var length int
	length, err = resp.Body.Read(data)
	if err != nil && err != io.EOF {
		return err
	}
	data = data[:length]
	re := regexp.MustCompile(`^<script>top.self.location.href='(.+)'</script>`)
	matches := re.FindSubmatch(data)
	if len(matches) <= 1 {
		return fmt.Errorf("未知错误，响应体:\n%s", string(data))
	}
	c.href = string(matches[1])
	return ErrNeedLogin
}

func (c *client) login(username, password string, netAccessType service) error {
	// 准备数据
	URL, err := url.Parse(c.href)
	if err != nil {
		return fmt.Errorf("parse URL error, err: %w", err)
	}

	query := doubleEncodeURIComponent(URL.RawQuery)
	username = doubleEncodeURIComponent(username)
	password = doubleEncodeURIComponent(password)
	service := doubleEncodeURIComponent(string(netAccessType))
	URL.Path = "/eportal/InterFace.do"
	URL.RawQuery = "method=login"
	// login
	content := "userId=" + username +
		"&password=" + password +
		"&service=" + service +
		"&queryString=" + query +
		"&operatorPwd=" + "" +
		"&operatorUserId=" + "" +
		"&validcode=" + "" +
		"&passwordEncrypt=" + "false"

	var req *http.Request
	req, err = newRequestWithCommonHeader(http.MethodPost, URL.String(), strings.NewReader(content))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	var resp *http.Response
	resp, err = c.Do(req)
	if err != nil {
		return err
	}
	defer drainBody(resp.Body)
	if resp.StatusCode != http.StatusOK {
		log.Printf("[Info] 登录接口返回状态码：%d %s\n", resp.StatusCode, resp.Status)
	}
	data, _ := io.ReadAll(resp.Body)
	res := &response{}
	err = json.Unmarshal(data, res)
	if err != nil {
		return fmt.Errorf("登录失败，响应正文：%s", string(data))
	}
	if res.Result != "success" {
		err = fmt.Errorf("login failed：%s", res.Message)
	}
	return err
}

func newRequestWithCommonHeader(method string, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err != nil {
		return
	}
	req.Header = commonHeader.Clone()
	return
}

func doubleEncodeURIComponent(raw string) string {
	return encodeURIComponent(encodeURIComponent(raw))
}

type response struct {
	Result  string
	Message string
}

func drainBody(reader io.ReadCloser) error {
	io.Copy(io.Discard, reader)
	return reader.Close()
}

func shouldEscape(c byte) bool {
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}
	switch c {
	case '-', '_', '.', '!', '~', '*', '\'', '(', ')':
		return false
	}

	return true
}

const upperhex = "0123456789ABCDEF"

func encodeURIComponent(s string) string {
	hexCount := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if shouldEscape(c) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return s
	}

	var buf [64]byte
	var t []byte

	required := len(s) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	j := 0
	for i := 0; i < len(s); i++ {
		switch c := s[i]; {
		case shouldEscape(c):
			t[j] = '%'
			t[j+1] = upperhex[c>>4]
			t[j+2] = upperhex[c&15]
			j += 3
		default:
			t[j] = s[i]
			j++
		}
	}
	return string(t)
}
