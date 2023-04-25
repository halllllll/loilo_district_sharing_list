package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/github.com/halllllll/loilo_district_sharing_list/config"
	"github.com/github.com/halllllll/loilo_district_sharing_list/input"
	"github.com/github.com/halllllll/loilo_district_sharing_list/output"
	"github.com/github.com/halllllll/loilo_district_sharing_list/scrape"
)

var (
	schoolId string
	userId   string
	userPw   string
	err      error
	urls     []string
)

func hello() {
	// envファイルがあればそれを読む
	cfg, err := config.Load()
	if err != nil {
		if configErr, ok := err.(config.ConfigError); ok {
			switch configErr.Situation {
			case config.Unmarshalling:
				fmt.Println(err)
				bufio.NewScanner(os.Stdin).Scan()
				log.Fatal(err)
			case config.LoadingConfigFile:
				fmt.Println("The config file cannot be found. Please write directly in the lines below.")
				// ユーザーの入力から取得
				inputReader := input.NewDefaultInputReader()
				schoolId, userId, userPw, err = inputReader.PromptAndReadCredentials()
				if err != nil {
					fmt.Println(err)
					bufio.NewScanner(os.Stdin).Scan()
					log.Fatal(err)
				}
			default:
				fmt.Printf("other: %s\n", configErr.Situation)
				bufio.NewScanner(os.Stdin).Scan()
				log.Fatal(err)
			}
		}
	} else {
		// configのenvファイルから取得できなかった場合
		schoolId = cfg.Shcool_id
		userId = cfg.User_id
		userPw = cfg.User_pw
	}

}

func main() {
	hello()

	// ログインを試みる
	fmt.Println("login:")
	loginCollector, err := scrape.Login(schoolId, userId, userPw)

	if err != nil {
		fmt.Println(err)
		bufio.NewScanner(os.Stdin).Scan()
		log.Fatal(err)
	}
	fmt.Println("start scraping:")
	// データ取得
	scrapeCollector := loginCollector.Clone()
	schools := scrape.SchoolList(scrapeCollector)

	// collyのQueueより、collyのAsyncとか、デフォルトのWorkGroupを使ったほうがよさそう。Collector.Asyncを使っている
	result, err := scrape.TeacherList(scrapeCollector, schools)
	if err != nil {
		fmt.Println(err)
		bufio.NewScanner(os.Stdin).Scan()
		log.Fatal(err)
	}
	fmt.Println("generating excel:")
	// 作成するシート名
	ct := time.Now().Format("2006_01_02_150406")
	fileName := fmt.Sprintf("loilo_share_%s.xlsx", ct)
	workbook := output.NewExcel(fileName)

	if err := workbook.FillSheet("share", result); err != nil {
		fmt.Println(err)
		bufio.NewScanner(os.Stdin).Scan()
		log.Fatal(err)
	}
	// 保存先
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		bufio.NewScanner(os.Stdin).Scan()
		log.Fatal(err)
	}
	if err := workbook.Save(cwd); err != nil {
		fmt.Println(err)
		bufio.NewScanner(os.Stdin).Scan()
		log.Fatal(err)
	}
	// デフォルトのシート破壊（どこでやるべきかはわからない）
	workbook.Wb.DeleteSheet("Sheet1")
	fmt.Println("done!")
	bufio.NewScanner(os.Stdin).Scan()
}
