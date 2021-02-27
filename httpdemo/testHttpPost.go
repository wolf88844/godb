package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func main() {
	testHttpPost()
}

func testHttpPost() {
	data := url.Values{
		"theCityName": {"重庆"},
	}

	reader := strings.NewReader(data.Encode())

	response, err := http.Post("http://www.webxml.com.cn/WebServices/WeatherWebService.asmx/getWeatherbyCityName",
		"application/x-www-form-urlencoded", reader)

	CheckErr(err)

	fmt.Printf("响应状态: %v\n", response.StatusCode)

	if response.StatusCode == 200 {
		defer response.Body.Close()
		fmt.Println("网络请求成功")
		CheckErr(err)
	} else {
		fmt.Println("请求失败", response.Status)
	}
}

func CheckErr(err error) {
	defer func() {
		if ins, ok := recover().(error); ok {
			fmt.Println("程序出现异常：", ins.Error())
		}
	}()
	if err != nil {
		panic(err)
	}
}
