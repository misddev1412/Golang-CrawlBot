package packages

import (
	"fmt"
	"net/url"
	"time"

	"aibot.com/news/models"
	"aibot.com/news/rds"
	"github.com/gocolly/colly/v2"
)

func CrawlInit(Cron models.Result, link, path string, id int, language string) {
	TotalLimit := float32(40)
	c := colly.NewCollector()

	// Find and visit all links
	c.OnXML("/html", func(e *colly.XMLElement) {
		dirTarget := e.ChildAttrs(path, "href")
		i := float32(0)
		LimitCron := Cron.Rate / TotalLimit * float32(100)
		var batchData models.NewsList
		for _, element := range dirTarget {
			if i >= LimitCron {
				break
			}
			linkDetail := element
			if !IsUrl(element) {
				linkParse, _ := url.Parse(link)
				linkDetail = linkParse.Scheme + "://" + linkParse.Host + element
			}

			newsData := CrawlDetail(linkDetail, id, language)
			batchData = append(batchData, newsData)
			i++
		}
		e.Request.Visit(e.Attr("href"))

		db := rds.MysqlConnect()
		models.InsertPost(batchData, db)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(link)
}

func CrawlDetail(link string, id int, language string) models.NewsDetail {

	var newsDetail models.NewsDetail
	c := colly.NewCollector()

	// Find and visit all links
	c.OnXML("/html", func(e *colly.XMLElement) {
		newsDetail.Title = e.ChildText(`//meta[@property="og:title"]/@content`)
		newsDetail.Image = e.ChildText(`//meta[@property="og:image"]/@content`)
		newsDetail.ShortContent = e.ChildText(`//meta[@property="og:description"]/@content`)
		newsDetail.TopicID = 1
		newsDetail.PageID = id
		newsDetail.Language = language
		newsDetail.Slug = GenerateSlug(newsDetail.Title)
		newsDetail.TargetLink = link
		newsDetail.CreatedAt = time.Now()
		newsDetail.UpdatedAt = time.Now()
	})

	c.Visit(link)

	return newsDetail
}
