package districtloilo

import (
	"fmt"

	"github.com/gocolly/colly/v2"
)

var (
	baseUrl  = "n.loilo.tv"
	loginUrl = fmt.Sprintf("https://%s/users/sign_in", baseUrl)
	homeUrl  = fmt.Sprintf("https://%s/dashboard", baseUrl)
)

type LoiloClient struct{}

func (lc *LoiloClient) Touch(districtId, UserId, UserPw string) (c *colly.Collector, err error) {
	c = colly.NewCollector(
		colly.AllowedDomains(baseUrl),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100"),
		colly.Async(true),
	)
	// login
	err = c.Post(loginUrl, map[string]string{
		"user[school][code]": districtId,
		"user[username]":     UserId,
		"user[password]":     UserPw,
		"commit":             "ログイン",
	})
	if err != nil {
		return nil, err
	}

	c.Wait()
	return c, nil
}
