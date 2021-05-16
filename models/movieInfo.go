package models

import (
	"errors"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var(
	db orm.Ormer
)

type MovieInfo struct{
	Id int64
	MovieId int64  				//电影ID
	MovieName string  			//名称
	MoviePic string				//电影图片url
	MovieDirector string		//导演
	MovieWriter string			//编剧
	MovieCountry string			//产地
	MovieLanguage string		//语言
	MovieMainCharacter string	//主演
	MovieType string			//类型
	MovieOnTime string			//上映时间
	MovieSpan string   			//时长
	MovieGrade string			//评分
	CreateTime string           //创建时间
	ModifyTime string           //更改时间
}

func init(){
	orm.Debug = true  //开启调试模式 调试模式下会打印出sql语句
	orm.RegisterDataBase("default", "mysql", "root:root@tcp(127.0.0.1:3306)/star?charset=utf8", 30)
	orm.RegisterModel(new(MovieInfo))
	db = orm.NewOrm()
}

func AddMovieInfo( movieInfo *MovieInfo)(int64, error){
	id, err := db.Insert(movieInfo)
	return id, err
}

//ID
func GetMovieId(movieHtml string)(int64, error){
	if movieHtml == ""{
		return 0, errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`SUBJECT_ID: '(.*?)',`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return 0, errors.New("get Id fail")
	}

	id, _ := strconv.ParseInt(result[0][1], 10, 64)
	return id, nil
}

//名字
func GetMovieName(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<span.*?property="v:itemreviewed">(.*)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Name fail")
	}

	return string(result[0][1]), nil
}
//封面
func GetMoviePic(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<img src="(.*?)".*? rel="v:image" />`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Pic fail")
	}

	return string(result[0][1]), nil
}

//导演
func GetMovieDirector(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<a.*?rel="v:directedBy">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Director fail")
	}

	var directors string
	for _, v := range result{
		directors += v[1] + "/"
	}

	return strings.Trim(directors, "/"), nil
}
//编剧
func GetMovieWriter(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<span class='pl'>编剧</span>: <span class='attrs'>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Writer fail")
	}

	var writers string
	for _, v := range result{
		reg := regexp.MustCompile(`<a.*?>(.*?)</a>`)
		result2 := reg.FindAllStringSubmatch(v[1], -1)
		if len(result2) == 0 {
			continue
		}

		for _, v := range result2{
			writers += v[1] + "/"
		}
	}
	return strings.Trim(writers, "/"), nil
}
//制片国家/地区
func GetMovieCountry(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<span class="pl">制片国家/地区:</span>(.*?)<br/>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Country fail")
	}

	return string(result[0][1]), nil
}
//语言
func GetMovieLanguage(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<span class="pl">语言:</span>(.*?)<br/>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Language fail")
	}

	return string(result[0][1]), nil
}
//类型
func GetMovieType(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<span property="v:genre">(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Type fail")
	}

	var types string
	for _, v := range result{
		types += v[1] + "/"
	}

	return strings.Trim(types, "/"), nil
}
//上映时间
func GetMovieOnTime(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<span property="v:initialReleaseDate" content="(.*?)">`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get OnTime fail")
	}

	timeArr := strings.Split(result[0][1], "(")
	if len(timeArr) == 0 {
		return "", errors.New("get OnTime fail")
	}

	var timeStamp time.Time
	if len(timeArr[0]) == 4 {
		timeStamp, _ = time.ParseInLocation("2006-01-02", timeArr[0] + "-01-01", time.Local)
	}else{
		timeStamp, _ = time.ParseInLocation("2006-01-02", timeArr[0], time.Local)
	}
	return time.Unix(timeStamp.Unix(), 0).Format("2006-01-02"), nil
}
//片长
func GetMovieSpan(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<span property="v:runtime".*?>(.*?)</span>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Span fail")
	}

	return string(result[0][1]), nil
}
//主演
func GetMovieMainCharacter(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<a.*?rel="v:starring">(.*?)</a>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get MainCharacter fail")
	}

	var mainCharacters string
	for _, v := range result{
		mainCharacters += v[1] + "/"
	}

	return strings.Trim(mainCharacters, "/"), nil
}

//评分
func GetMovieGrade(movieHtml string)(string, error){
	if movieHtml == ""{
		return "", errors.New("movieHtml is null")
	}

	reg := regexp.MustCompile(`<strong.*?property="v:average">(.*)</strong>`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return "", errors.New("get Grade fail")
	}

	return string(result[0][1]), nil
}


func GetMovieUrls(movieHtml string)[]string{
	if movieHtml == ""{
		return []string{}
	}

	reg := regexp.MustCompile(`<a.*?href="(https://movie.douban.com/subject.*?from=subject-page)" >`)
	result := reg.FindAllStringSubmatch(movieHtml, -1)
	if len(result) == 0 {
		return []string{}
	}

	movieUrlSet := make([]string, 0)
	for  _, v := range result{
		movieUrlSet = append(movieUrlSet, v[1])
	}
	return movieUrlSet
}

func GetMovieHtml(url string) string{
	res := httplib.Get(url)
	movieHtml, err := res.String()
	if err != nil {
		return ""
	}
	return movieHtml
}

func GetMovieInfo(movieHtml string)(*MovieInfo, error){
	movieId, err := GetMovieId(movieHtml)
	if err != nil{
		return nil, err
	}
	movieName, err := GetMovieName(movieHtml)
	if err != nil{
		return nil, err
	}
	moviePic, err := GetMoviePic(movieHtml)
/*	if err != nil{
		return nil, err
	}*/
	movieDirector, err := GetMovieDirector(movieHtml)
/*	if err != nil{
		return nil, err
	}*/
	movieWriter, err := GetMovieWriter(movieHtml)
/*	if err != nil{
		return nil, err
	}*/
	movieCountry, err := GetMovieCountry(movieHtml)
/*	if err != nil{
		return nil, err
	}*/
	movieLanguage, err := GetMovieLanguage(movieHtml)
/*	if err != nil{
		return nil, err
	}*/
	movieMainCharacter, err := GetMovieMainCharacter(movieHtml)
/*	if err != nil{
		return nil, err
	}*/
	movieType, err := GetMovieType(movieHtml)
/*	if err != nil{
		return nil, err
	}*/
	movieOnTime, err := GetMovieOnTime(movieHtml)
	if err != nil{
		movieOnTime = "2006-02-03"
	}
	movieSpan, err := GetMovieSpan(movieHtml)
/*	if err != nil{
		return nil, err
	}*/
	movieGrade, err := GetMovieGrade(movieHtml)
/*	if err != nil{
		return nil, err
	}*/

	return &MovieInfo{
		MovieId:movieId,
		MovieName: movieName,
		MoviePic: moviePic,
		MovieDirector: movieDirector,
		MovieWriter: movieWriter,
		MovieCountry: movieCountry,
		MovieLanguage: movieLanguage,
		MovieMainCharacter: movieMainCharacter,
		MovieType: movieType,
		MovieOnTime: movieOnTime,
		MovieSpan: movieSpan,
		MovieGrade: movieGrade,
		CreateTime: time.Now().Format("2006-01-02 15:04:05"),
		ModifyTime: time.Now().Format("2006-01-02 15:04:05"),
	}, nil
}

func GetMovieInfos(url string, num int)[]*MovieInfo{
	if num <= 0 || num >= 100 {
		num = 10   //默认值
	}

	movieInfos := make([]*MovieInfo, 0)       //保存结果
	allUrlsSet := make([]string, 0)           //保存访问链接
	allUrlsMap := make(map[string]struct{})   //用于去重
	allUrlsSet = append(allUrlsSet, url)
	count := 0
	for i := 0; i < len(allUrlsSet); i++ {
		movieHtml := GetMovieHtml(allUrlsSet[i])
		movieInfo, err := GetMovieInfo(movieHtml)
		if err == nil{
			movieInfos = append(movieInfos, movieInfo)
			count += 1
		}else{
			fmt.Println(err)
		}

		urls := GetMovieUrls(movieHtml)
		for _, url := range urls {
			if _, ok := allUrlsMap[url]; ok {
				continue
			}
			allUrlsMap[url] = struct{}{}
			allUrlsSet = append(allUrlsSet, url)
		}

		if count >= num{
			break
		}
	}

	return movieInfos
}

