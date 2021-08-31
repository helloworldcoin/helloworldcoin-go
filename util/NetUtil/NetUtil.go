package NetUtil

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(requestUrl string, requestBody string) string {
	//TODO 异常返回nil
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(requestUrl, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		panic(err)
	}
	//TODO
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
