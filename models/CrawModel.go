package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Result struct {
	Link     string
	Xpath    string
	ID       int
	Language string  `gorm:"column:language"`
	Rate     float32 `gorm:"column:rate"`
}

type NewsDetail struct {
	PageID       int       `gorm:"column:page_id"`
	TopicID      int       `gorm:"column:topic_id"`
	Slug         string    `gorm:"column:slug"`
	Title        string    `gorm:"column:title"`
	ShortContent string    `gorm:"column:short_content"`
	TargetLink   string    `gorm:"column:target_link"`
	Image        string    `gorm:"column:image"`
	Language     string    `gorm:"column:language"`
	CreatedAt    time.Time `gorm:"column:created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at"`
}

type NewsList []NewsDetail

func ListCrawl(db *gorm.DB) []Result {
	var result []Result
	db.Table("news_pages").Select("link", "xpath", "id", "language", "rate").Where("is_croned", "0").Limit(1).Scan(&result)

	return result
}

func InsertPost(post NewsList, db *gorm.DB) int64 {
	result := db.Table("tmp_posts").Select("page_id", "topic_id", "slug", "title", "short_content", "image", "target_link", "language", "created_at", "updated_at").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&post)
	return result.RowsAffected

}

func InsertPostFromTmp(post NewsList, db *gorm.DB) int64 {
	result := db.Table("posts").Select("page_id", "topic_id", "slug", "title", "short_content", "image", "target_link", "language", "created_at", "updated_at").Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&post)
	return result.RowsAffected

}

func UpdateCron(Cron Result, db *gorm.DB) int64 {
	result := db.Table("news_pages").Model(Cron).Where("is_croned", "0").Update("is_croned", "1")
	fmt.Println(result.RowsAffected)
	return result.RowsAffected

}

func UpdateNewCron(db *gorm.DB) int64 {
	result := db.Exec("UPDATE users SET is_croned = ? WHERE is_croned = 1", "0")

	fmt.Println(result.RowsAffected)
	return result.RowsAffected

}

func ListPostShuffle(db *gorm.DB) NewsList {
	var news NewsList
	db.Table("tmp_posts").Select("page_id", "topic_id", "slug", "title", "short_content", "target_link", "image", "language", "created_at", "updated_at").Where("title != ''").Order("RAND()").Scan(&news)

	return news
}

func DeletePostFromTmp(db *gorm.DB) int64 {
	result := db.Exec("DELETE FROM tmp_posts WHERE 1=1")
	return result.RowsAffected

}

func DeleteOldPost(db *gorm.DB) int64 {
	result := db.Exec("DELETE FROM posts WHERE HOUR(TIMEDIFF(NOW(), created_at)) > 24")
	return result.RowsAffected

}
