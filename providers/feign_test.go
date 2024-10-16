package providers

import (
	"log"
	"os"
	"testing"
)

func TestHandleFeignRequest(t *testing.T) {
	params, headers := map[string]string{}, map[string]string{
		// "Host":       "https://translate.alibaba.com",
		// "User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36",
	}

	response, error := HandleFeignRequest("https://translate.alibaba.com/api/translate/csrftoken", "GET", params, headers)

	if error != nil {
		log.Fatalln(error.Error())
		return
	}

	log.Println(response)
}

func TestHandleFeignPutFileRequest(t *testing.T) {
	file, error := os.ReadFile("/Users/wdy/Desktop/app.txt")

	if error != nil {
		log.Fatalln(error.Error())
		return
	}

	if _, error := HandleFeignPutFileRequest("http://127.0.0.1:5000/data/hello world/b.txt", file, map[string]string{}); error != nil {
		log.Fatalln(error.Error())
		return
	}

	log.Println("put success")
}
