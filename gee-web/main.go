package main

import (
	"./gee"
	"fmt"
	"html/template"
	"net/http"
	"time"
)

//func onlyForV2() gee.HandlerFunc{
//	return func(c *gee.Context) {
//		t:=time.Now()
//		c.Fail(500,"Internal Server Error")
//		log.Printf("[%d] %s in %v for group v2",c.StatusCode,c.Req.RequestURI,time.Since(t))
//	}
//}
type student struct {
	Name string
	Age int8
}

func formatAsDate(t time.Time) string {
	year,month,day :=t.Date()
	return fmt.Sprintf("%d-%02d-%02d",year,month,day)
}
func main(){
	r := gee.New()
	r.Use(gee.Logger())
	r.SetFuncMap(template.FuncMap{
		"formatAsDate":formatAsDate,
	})
	r.LoadHTMLGlob("D:/go-workspace/godemo/gee-web/templates/*")
	r.Static("/assets","D:/go-workspace/godemo/gee-web/static")

	stu1:=&student{Name:"Geek",Age:20}
	stu2:=&student{Name:"Jack",Age:22}
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK,"css.tmpl",nil)
	})
	r.GET("/students",func(c *gee.Context){
		c.HTML(http.StatusOK,"arr.tmpl",gee.H{
			"title":"gee",
			"stuArr":[2]*student{stu1,stu2},
		})
	})

	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK,"custom_func.tmpl",gee.H{
			"title":"gee",
			"now":time.Date(2020,3,8,16,34,0,0,time.UTC),
		})
	})

	r.Run(":9999")
}
