package fetcher

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io/ioutil"
	"log"
	"net/http"
)

func Fetch(url string) ([]byte, error) {

	log.Printf("Fetching url:%s", url)

	client := &http.Client{}

	//获取要请求的 url
	resq, err := http.NewRequest(http.MethodGet, url, nil)
	client.Do(resq)
	if err != nil {
		return nil, err
	}
	resq.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	resq.Header.Add("Referer", "http://www.zhenai.com/zhenghun/shanghai/nv")
	resq.Header.Add("Cookie", "sid=c55d8b19-ff31-4cf4-b1aa-6b5ca6456f86; ipCityCode=10127001; ipOfflineCityCode=10127001; Hm_lvt_2c8ad67df9e787ad29dbd54ee608f5d2=1555905496,1555940147,1555983971; Hm_lpvt_2c8ad67df9e787ad29dbd54ee608f5d2=1556005948")
	resp, err := client.Do(resq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	//防止 Reader 过
	bodyReader := bufio.NewReader(resp.Body)
	//获取url 的字符编码
	e := determineEncoding(bodyReader)

	utf8Reader := transform.NewReader(bodyReader, e.NewDecoder())
	//最近返回值
	return ioutil.ReadAll(utf8Reader)
}

/**
 * 设置1024
 * 自动转换成 对应的字符编码
 */
func determineEncoding(r *bufio.Reader) encoding.Encoding {
	//默认只读1024个字节,读了以后就不能再读
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
