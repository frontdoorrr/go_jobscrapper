package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)


var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"

func main() {
	totalPages := getPages()
	fmt.Println(totalPages)



}


func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()


	doc, err := goquery.NewDocumentFromReader(res.Body)
	doc.Find(".pagination").Each(func(i int, s*goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status Code :", res.StatusCode)

	}
}
