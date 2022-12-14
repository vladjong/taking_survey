package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

type Client struct {
	Client         http.Client
	NumberQuestion int
	Context        context.Context
	Cookies        []*http.Cookie
	Headers        map[string]string
	Ratelimiter    *rate.Limiter
}

func NewClinet(ctx context.Context) *Client {
	cookieJar, _ := cookiejar.New(nil)
	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return &Client{
		Client:         http.Client{Jar: cookieJar},
		NumberQuestion: 1,
		Context:        ctx,
		Headers:        headers,
		Ratelimiter:    rate.NewLimiter(rate.Limit(viper.GetInt("rps")), 1),
	}
}

func (c *Client) Run() error {
	_, cookies, err := c.getPage(http.MethodGet, viper.GetString("url"), nil, viper.GetInt("timeout"))
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	c.Cookies = cookies
	doc, _, err := c.getPage(http.MethodGet, c.getLink(), nil, viper.GetInt("timeout"))
	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}
	for {
		formDatas := c.parseForm(doc)
		doc, _, err = c.getPage(http.MethodPost, c.getLink(), formDatas, viper.GetInt("timeout"))
		if err != nil {
			return fmt.Errorf("%s", err.Error())
		}
		if strings.Contains(doc.Text(), "Test successfully passed") {
			return nil
		}
		c.NumberQuestion += 1
	}
}

func (c *Client) getPath() string {
	return viper.GetString("link") + strconv.Itoa(c.NumberQuestion)
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

func (c *Client) getPage(method, siteURL string, formDatas map[string]string, timeout int) (*goquery.Document, []*http.Cookie, error) {
	body := io.Reader(nil)
	if len(formDatas) > 0 {
		form := url.Values{}
		for key, val := range formDatas {
			form.Add(key, val)
		}
		body = strings.NewReader(form.Encode())
	}
	req, err := http.NewRequestWithContext(c.Context, method, siteURL, body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create http request context: %w", err)
	}
	if len(c.Headers) > 0 {
		for key, val := range c.Headers {
			req.Header.Add(key, val)
		}
	}
	if len(c.Cookies) > 0 {
		for _, val := range c.Cookies {
			req.AddCookie(val)
		}
	}
	reqTimeout := 10 * time.Second
	if timeout != 0 {
		reqTimeout = time.Duration(timeout) * time.Second
	}
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   reqTimeout,
	}
	err = c.Ratelimiter.Wait(c.Context)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to rps: %w", err)
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
	if resp.StatusCode != 200 {
		return nil, nil, fmt.Errorf("failed to connect: %w", err)
	}
	return doc, resp.Cookies(), nil
}
