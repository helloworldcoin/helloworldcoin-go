package NetUtil

/*
 @author king 409060350@qq.com
*/

import (
	"bytes"
	"helloworld-blockchain-go/util/LogUtil"
	"io/ioutil"
	"net/http"
	"time"
)

func Get(requestUrl string, requestBody string) string {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Post(requestUrl, "application/json", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		LogUtil.Debug(err.Error())
		return ""
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	return string(result)
}
