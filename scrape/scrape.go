package scrape

import (
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type School struct {
	SchoolName string
	SchoolDist string
}

func SchoolList(c *colly.Collector) []School {
	// 訪問するURLが含まれたパス
	permission_url := "https://n.loilo.tv/district/document_group_permissions/schools"
	// ページネーション（見た目）
	page := 1
	// 情報取得してくるURL
	result := []School{}

	c.OnResponse(func(r *colly.Response) {
		// fmt.Printf("status: %d\n", r.StatusCode)
	})

	c.OnRequest(func(r *colly.Request) {
		// set user agent
		r.Headers.Set("User-Agent", Ua())
	})

	c.OnHTML("html body div#main div.container form table.table tbody", func(e *colly.HTMLElement) {
		rowCount := 0
		e.DOM.Find("tr").Each(func(rowIdx int, tr *goquery.Selection) {
			rowCount++
			school := new(School)
			tr.Find("td").Each(func(tdIdx int, td *goquery.Selection) {
				// 特にtdにclassやidを設定しているわけではないっぽいので、tdのindexで決め打ちする
				switch tdIdx {
				case 0: // 学校ID想定

				case 1: // 学校名想定
					school.SchoolName = td.Text()
				case 2:
					aTag := td.Find("a")
					if aTag.Length() > 0 {
						if link, exist := aTag.Attr("href"); exist {
							school.SchoolDist = link
						}
					}
				}
			})
			result = append(result, *school)
		})
		if rowCount != 0 {
			page++
			c.Visit(fmt.Sprintf("%s?page=%d", permission_url, page))
		}
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", Ua())
	})

	c.Visit(fmt.Sprintf("%s?page=%d", permission_url, page))
	c.Wait()
	fmt.Printf("length: %d\n", len(result))
	return result
}

func TeacherList(c *colly.Collector, schools []School) ([][]string, error) {
	// コールバック関数でエラー吐いても途中で返せないのでとりあえず保管する
	errMsg := []string{}

	result := [][]string{{"school", "name", "checked"}}
	// onHTMLからどの学校かを取得するにはURLくらいしか情報がないので辞書を作る
	urlBySchoolMap := map[string]string{}
	for _, val := range schools {
		urlBySchoolMap[val.SchoolDist] = val.SchoolName
	}

	baseUrl := "https://n.loilo.tv"
	asyncC := c.Clone()
	// 非同期設定
	asyncC.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       1 * time.Second,
		RandomDelay: 1 * time.Second,
		Parallelism: 5,
	})

	asyncC.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
		// fmt.Printf("status: %d\n", r.StatusCode)
	})

	// リファラを設定し、「前のページからリンクを辿ってきた」かのように振る舞わせる
	// あとUAも偽造
	asyncC.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Referer", r.URL.String())
		r.Headers.Set("User-Agent", Ua())
	})

	asyncC.OnHTML("div.row:nth-child(1)", func(e *colly.HTMLElement) {
		u, err := url.Parse(e.Request.URL.String())
		if err != nil {
			errMsg = append(errMsg, err.Error())
		}
		e.DOM.Find("div").Each(func(rowIdx int, div *goquery.Selection) {
			label := div.Find("label")
			div.Find("div>input").Each(func(inputIdx int, input *goquery.Selection) {
				if t, exist := input.Attr("type"); exist && t == "checkbox" {
					val, _ := input.Attr("checked")
					result = append(result, []string{urlBySchoolMap[u.Path], label.Text(), val})
				}
			})
		})
	})

	asyncC.OnError(func(r *colly.Response, err error) {
		errMsg = append(errMsg, err.Error())
	})

	for _, school := range schools {
		fmt.Printf("school name: %s\n", school.SchoolName)
		asyncC.Visit(fmt.Sprintf("%s%s", baseUrl, school.SchoolDist))
	}
	asyncC.Wait()
	if len(errMsg) != 0 {
		return nil, fmt.Errorf(strings.Join(errMsg, ","))
	}
	return result, nil
}
