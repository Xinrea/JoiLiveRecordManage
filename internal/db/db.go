package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

func Init() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s?parseTime=true", viper.GetString("database.user"), viper.GetString("database.password"), "tcp", viper.GetString("database.host"), 3306, viper.GetString("database.dbname"))
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db = d
}

type DanmuDB struct {
	CreateTime time.Time      `json:"create_time"`
	Room       int            `json:"room"`
	MedalName  sql.NullString `json:"medal_name"`
	UserName   string         `json:"user_name"`
	Content    string         `json:"content"`
}

type Danmu struct {
	CreateTime time.Time `json:"create_time"`
	Room       int       `json:"room"`
	MedalName  string    `json:"medal_name"`
	UserName   string    `json:"user_name"`
	Content    string    `json:"content"`
}

func GetDanmu(room int, text string) []Danmu {
	rows, err := db.Query("select create_time, medal_name, user_name, content from danmurecord where room = ? and content like ? order by create_time desc", room, "%"+text+"%")
	if err != nil {
		log.Print(err)
		return nil
	}
	var danmuList []Danmu
	for rows.Next() {
		var d DanmuDB
		d.Room = room
		_ = rows.Scan(&d.CreateTime, &d.MedalName, &d.UserName, &d.Content)
		var d2 Danmu
		d2.CreateTime = d.CreateTime
		d2.Room = d.Room
		d2.MedalName = d.MedalName.String
		d2.UserName = d.UserName
		d2.Content = d.Content
		danmuList = append(danmuList, d2)
	}
	return danmuList
}
