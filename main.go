package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"strconv"

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
	client := NewClinet()
	if err := client.getFirstLink(); err != nil {
		log.Fatalf("error get first link: %s", err.Error())
	}
	for {
		client.takingSurvey()
		client.NumberQuestion += 1
		return
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

func (c *Client) getFirstLink() error {
	resp, err := c.Client.Get(viper.GetString("url"))
	if err != nil {
		log.Fatalf("error get url: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("error initializing goquery response body: %s", err.Error())
	}
	link, ok := doc.Find("a").Attr("href")
	if !ok {
		log.Fatal(err)
	}
	if link == c.getPath() {
		return nil
	}
	return fmt.Errorf("error this is incorect page")
}

func (c *Client) getPath() string {
	return config.Link + strconv.Itoa(c.NumberQuestion)
}

func (c *Client) getLink() string {
	return viper.GetString("url") + c.getPath()
}

func (c *Client) takingSurvey() {
	resp, err := c.Client.Get(c.getLink())
	if err != nil {
		log.Fatalf("error get url: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Fatalf("failed to fetch data: %d %s", resp.StatusCode, resp.Status)
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatalf("error initializing goquery response body: %s", err.Error())
	}
	fmt.Println(doc.Html())
	fmt.Println("\\n")
	form := doc.Find("form").First()
	fmt.Println(form.Text())
	form.Find("p").Each(func(i int, s *goquery.Selection) {
		s.Find("select").Each(func(i int, s *goquery.Selection) {
			if name, _ := s.Attr("name"); name == "" {
				return
			}
			s.Find("option").Each(func(i int, s *goquery.Selection) {
				value, ok := s.Attr("value")
				if !ok {
					return
				}
				fmt.Println(value)
				// name, _ := s.Attr("name")
				// if name == "" {
				// 	return
				// }
				// fmt.Println(s.Text())
			})
		})
	})
}
