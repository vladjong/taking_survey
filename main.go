package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"github.com/vladjong/taking_survey/config"
)

type Client struct {
	Client         http.Client
	NumberQuestion int
}

func main() {
	if err := InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	ctx := context.Background()
	client := NewClinet()
	_, cookies, err := client.GetPage(ctx, http.MethodGet, viper.GetString("url"), nil, nil, nil, 10)
	if err != nil {
		log.Fatalf("error get first link: %s", err.Error())
	}
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	doc, _, err := client.GetPage(ctx, http.MethodGet, client.getLink(), cookies, headers, nil, 10)
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}
	for {
		formDatas := client.parseForm(doc)
		doc, _, err = client.GetPage(ctx, http.MethodPost, client.getLink(), cookies, headers, formDatas, 10)
		if err != nil {
			log.Fatalf("error: %s", err.Error())
		}
		if strings.Contains(doc.Text(), "Test successfully passed") {
			fmt.Println("ПРОШЕЛ ТЕСТ")
			return
		}
		client.NumberQuestion += 1
	}
}

func InitConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func NewClinet() *Client {
	cookieJar, _ := cookiejar.New(nil)
	return &Client{
		Client:         http.Client{Jar: cookieJar},
		NumberQuestion: 1,
	}
}

func (c *Client) getPath() string {
	return config.Link + strconv.Itoa(c.NumberQuestion)
}

func (c *Client) getLink() string {
	return viper.GetString("url") + c.getPath()
}

func (c *Client) parseForm(doc *goquery.Document) map[string]string {
	formDatas := make(map[string]string)
	form := doc.Find("form")
	form.Find("p").Each(func(i int, s *goquery.Selection) {
		c.parseSelect(s, formDatas)
		c.parseType(s, formDatas)
	})
	return formDatas
}

func (c *Client) parseSelect(s *goquery.Selection, formDatas map[string]string) {
	s.Find("select").Each(func(i int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if name == "" || !ok {
			return
		}
		s.Find("option").Each(func(i int, s *goquery.Selection) {
			value, ok := s.Attr("value")
			if !ok {
				return
			}
			if len(value) > len(formDatas[name]) {
				formDatas[name] = value
			}
		})
	})
}

func (c *Client) parseType(s *goquery.Selection, formDatas map[string]string) {
	s.Find("input").Each(func(i int, s *goquery.Selection) {
		name, ok := s.Attr("name")
		if name == "" || !ok {
			return
		}
		typ, ok := s.Attr("type")
		if !ok {
			return
		}
		if typ == "text" {
			formDatas[name] = "test"
		} else if typ == "radio" {
			value, ok := s.Attr("value")
			if !ok {
				return
			}
			if len(value) > len(formDatas[name]) {
				formDatas[name] = value
			}
		}
	})
}

func maxString(value string, temp *string) {
	if len(value) > len(*temp) {
		*temp = value
	}
}

func (c *Client) GetPage(ctx context.Context, method, siteURL string, cookies []*http.Cookie, headers, formDatas map[string]string, timeout int) (*goquery.Document, []*http.Cookie, error) {
	body := io.Reader(nil)
	if len(formDatas) > 0 {
		form := url.Values{}
		for k, v := range formDatas {
			form.Add(k, v)
		}
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequestWithContext(ctx, method, siteURL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create http request context: %w", err)
	}
	if len(headers) > 0 {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	if len(cookies) > 0 {
		for _, c := range cookies {
			req.AddCookie(c)
		}
	}
	reqTimeout := 10 * time.Second
	if timeout != 0 {
		reqTimeout = time.Duration(timeout) * time.Second
	}
	httpClient := &http.Client{
		Transport:     http.DefaultTransport,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       reqTimeout,
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute http request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse html: %w", err)
	}
	return doc, resp.Cookies(), nil
}
