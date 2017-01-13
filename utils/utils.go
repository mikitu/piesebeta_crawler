package utils

import (
	"net/http"
	"regexp"
	"io/ioutil"
	"github.com/mikitu/piesebeta_crawler/task"
	"golang.org/x/net/html"
	"github.com/mikitu/piesebeta_crawler/qutils"
	"bytes"
	"encoding/gob"
	"log"
	"time"
)
const QueueModels = "crawler::get-models"
const QueueModelCategories = "crawler::model-categories"

func GetModels() {
	pool := NewPool(50)
	url := "http://partsfinder.onlinemicrofiche.com/americanbeta/showmodel.asp?type=17&make=betamc&modelid=Beta%20Dirt%20Bike"
	resp, _ := http.Get(url)
	defer resp.Body.Close()

	var re = regexp.MustCompile(`(?iUm)<tr Class=Model_Line_[^>]+>[^<]+<TD Class=Model_LineYear_[^>]+>(\d+)</TD>[^<]+<TD Class=Model_LineModel_[^>]+><A href="([^"]+)" target="_top">([^>]+)</A></TD>[^<]+</tr>`)
	htmlData, _ := ioutil.ReadAll(resp.Body)
	for _, match := range re.FindAllStringSubmatch(string(htmlData), -1) {
		match = append(match[:0], match[1:]...)
		pool.Exec(task.NewMessagePublishModelTask(match, QueueModels))
	}
	pool.Close()
	pool.Wait()
}

func ReadModels() {
	//pool := NewPool(50)
	conn, ch := qutils.GetChannel(qutils.Qurl)
	defer ch.Close()
	defer conn.Close()
	msgs, err := ch.Consume(
		QueueModels, //queue string,
		"",    //consumer string,
		false, //autoAck bool,
		true,  //exclusive bool,
		false, //noLocal bool,
		false, //noWait bool,
		nil)   //args amqp.Table)

	if err != nil {
		log.Fatalln("Failed to get access to messages")
	}
	type ti interface{}
	for msg := range msgs {
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)
		sd := make(map[string]string)
		err = dec.Decode(&sd)
		if err != nil {
			log.Fatal("decode:", err)
		}

		log.Printf("%+v\n", sd["url"])
		resp, _ := http.Get(sd["url"])
		var re = regexp.MustCompile(`(?iUs)<TABLE Class=sel_sec_list>(.*)</table>`)
		htmlData, _ := ioutil.ReadAll(resp.Body)
		cleanBody := CleanStr(string(htmlData))
		safe := RegexpReplace(cleanBody, "return overlib[^\"]+", "")

		for _, match := range re.FindAllStringSubmatch(string(safe), -1) {
			//match = append(match[:0], match[1:]...)
			log.Printf("%+v\n", match)
		}
		resp.Body.Close()
		msg.Ack(false)
		time.Sleep(5*time.Second)
	}



	//pool.Close()
	//pool.Wait()
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

func CleanStr(str string) string {
	safe := RegexpReplace(str, `(?mis)\n\r`, "")
	safe = RegexpReplace(safe, `(?mis)\n`, "")
	safe = RegexpReplace(safe, `(?mis)\t`, "")
	safe = RegexpReplace(safe, `(?mis)\s+`, " ")
	safe = RegexpReplace(safe, `(?mis)> <`, "><")
	return safe
}

func RegexpReplace(str, pattern, replace string) string {
	reg, err := regexp.Compile(pattern)
	if err != nil {
		log.Fatal(err)
	}
	safe := reg.ReplaceAllString(str, replace)
	return safe
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
