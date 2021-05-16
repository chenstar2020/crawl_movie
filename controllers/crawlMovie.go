package controllers

import (
	"crawl_movie/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
)

type CrawlMovieController struct {
	beego.Controller
}

type CrawlMovieReq struct{
	Url string  `json:"url"`
	Num int		`json:"num"`
}

type CrawlMovieRes struct{
	Errno  int `json:"errno"`
	Errmsg string `json:"errmsg"`
}

func (c *CrawlMovieController) CrawlMovie() {
	var req CrawlMovieReq
	var res CrawlMovieRes

	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &req); err != nil{
		res.Errmsg = err.Error()
		res.Errno = 1
	}else{
		data := models.GetMovieInfos(req.Url, req.Num)
		for _, info := range data {
			_, err := models.AddMovieInfo(info)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}

		res.Errmsg = "success"
		res.Errno = 0
	}

	resData, _ := json.Marshal(res)
	c.Data["json"] = string(resData)
	c.ServeJSON()
	return
}
