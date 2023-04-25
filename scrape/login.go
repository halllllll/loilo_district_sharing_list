package scrape

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func Login(schoolId, userId, userPw string) (*colly.Collector, error) {
	url := "https://n.loilo.tv/users/sign_in"

	c := colly.NewCollector(
		colly.AllowedDomains("n.loilo.tv"),
		colly.UserAgent(Ua()),
		colly.Async(true), // あとでスクレイピングするときに使う
	)
	// random delay
	c.Limit(&colly.LimitRule{
		RandomDelay: 3 * time.Second,
	})

	success := false
	c.OnHTML("div.admin-menu-item:nth-child(2) > a:nth-child(1) > p:nth-child(3)", func(e *colly.HTMLElement) {
		if e.Text == "自治体に関する設定" {
			success = true
		}
	})

	// login
	err := c.Post(url, map[string]string{
		"user[school][code]": schoolId,
		"user[username]":     userId,
		"user[password]":     userPw,
		"commit":             "ログイン",
	})
	c.Wait()

	// ログイン確認
	c.Visit("https://n.loilo.tv/dashboard")

	c.Wait()

	if err != nil {
		return nil, err
	}
	if success != true {
		return nil, fmt.Errorf("login failed with the∂ `district account`")
	}
	return c, nil
}
