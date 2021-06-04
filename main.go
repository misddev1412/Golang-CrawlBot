package main

import (
	"math"
	"time"

	"aibot.com/news/models"
	"aibot.com/news/packages"
	"aibot.com/news/rds"
	"gopkg.in/robfig/cron.v2"
)

func main() {

	loc, _ := time.LoadLocation("Asia/Ho_Chi_Minh")
	time.Local = loc // -> this is setting the global timezone

	c := cron.New()
	c.AddFunc("@every 0h0m5s", executeCraw)
	c.AddFunc("@every 0h10m0s", shuffleTmpToPost)
	c.AddFunc("@every 1h0m0s", updateNewCron)
	c.AddFunc("@every 1h0m0s", deleteOldPost)
	c.Start()

	// Added time to see output
	time.Sleep(math.MaxInt64)
}

func executeCraw() {
	db := rds.MysqlConnect()
	dataCraw := models.ListCrawl(db)
	sqlDB, _ := db.DB()
	sqlDB.Close()
	for _, element := range dataCraw {
		packages.CrawlInit(element, element.Link, element.Xpath, element.ID, element.Language)

		db := rds.MysqlConnect()
		sqlDB, _ := db.DB()
		models.UpdateCron(element, db)
		sqlDB.Close()
	}

}

func shuffleTmpToPost() {
	db := rds.MysqlConnect()
	dataNews := models.ListPostShuffle(db)
	models.InsertPostFromTmp(dataNews, db)
	models.DeletePostFromTmp(db)
	sqlDB, _ := db.DB()
	sqlDB.Close()

}

func updateNewCron() {
	db := rds.MysqlConnect()
	models.UpdateNewCron(db)
	sqlDB, _ := db.DB()
	sqlDB.Close()

}

func deleteOldPost() {
	db := rds.MysqlConnect()
	models.DeleteOldPost(db)
	sqlDB, _ := db.DB()
	sqlDB.Close()

}
