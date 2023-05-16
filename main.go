package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)


var baseURL string = "https://www.saramin.co.kr/zf_user/search/recruit?&searchword=python"

func main() {

	var jobs [] extractedJob
	totalPages := getPages()
	fmt.Println(totalPages)
	fmt.Println(totalPages)
	fmt.Println(totalPages)
	fmt.Println(totalPages)
	for i := 0; i < totalPages; i++ {
		extractedJobs := getPage(i)
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
}

type extractedJob struct{
	id string
	title string
	condition string
}

 func getPages() int {
	pages := 0
	res, err := http.Get(baseURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)



	doc.Find(".pagination").Each(func(i int, s *goquery.Selection){
		pages = s.Find("a").Length()
	})
	return pages
 }

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"ID", "Title", "Condition"}
	wErr := w.Write(headers)
	checkErr(wErr)
	for _, job := range jobs {
		jobSlice := []string{job.id, job.title, job.condition}
		jwErr := w.Write(jobSlice)
		checkErr(jwErr)
	}
}

func getPage(page int) []extractedJob{
	var jobs [] extractedJob
	pageURL := baseURL + "&recruitPage=" + strconv.Itoa(page) + "&recruitSort=relation&recruitPageCount=40&inner_com_type=&company_cd=0%2C1%2C2%2C3%2C4%2C5%2C6%2C7%2C9%2C10&show_applied=&quick_apply=&except_read=&ai_head_hunting=&mainSearch=n"
	fmt.Println("requesting From ...   ", pageURL)
	res, err := http.Get(pageURL)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	searchcards := doc.Find(".item_recruit")

	searchcards.Each(func(i int, card *goquery.Selection){

		job := extractJob(card)
		jobs = append(jobs, job)
	})
	return jobs
}

func extractJob(card *goquery.Selection) extractedJob{
		id, _ := card.Attr("value")
		title := cleanString(card.Find(".job_tit>a").Text())
		condition := cleanString(card.Find(".job_condition").Text())
		fmt.Println(id, title, condition)

		return extractedJob{id:id, title:title, condition: condition,}
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

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str))," ")
}
