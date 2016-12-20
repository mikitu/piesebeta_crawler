package utils

import (
	"net/http"
	"regexp"
	"io/ioutil"
	"github.com/mikitu/piesebeta_crawler/task"
	"golang.org/x/net/html"
)

func GetModels() {
	pool := NewPool(50)
	url := "http://partsfinder.onlinemicrofiche.com/americanbeta/showmodel.asp?type=17&make=betamc&modelid=Beta%20Dirt%20Bike"
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	queue := "crawler::get-models"
	var re = regexp.MustCompile(`(?iUm)<tr Class=Model_Line_[^>]+>[^<]+<TD Class=Model_LineYear_[^>]+>(\d+)</TD>[^<]+<TD Class=Model_LineModel_[^>]+><A href="([^"]+)" target="_top">([^>]+)</A></TD>[^<]+</tr>`)
	htmlData, _ := ioutil.ReadAll(resp.Body)
	for _, match := range re.FindAllStringSubmatch(string(htmlData), -1) {
		match = append(match[:0], match[1:]...)
		pool.Exec(task.NewMessagePublishTask(match, queue))
	}
	pool.Close()
	pool.Wait()
}

func ReadModels() {

}

func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}

	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}
//z := html.NewTokenizer(resp.Body)
//for {
//	tt := z.Next()
//	switch {
//	case tt == html.ErrorToken:
//		// End of the document, we're done
//		return
//	case tt == html.StartTagToken:
//		t := z.Token()
//
//		isAnchor := t.Data == "a"
//		if isAnchor {
//			fmt.Println("We found a link!")
//		}
//		for _, a := range t.Attr {
//			if a.Key == "href" {
//				fmt.Println("Found href:", a.Val)
//				break
//			}
//		}
//	}
//}
