package main

import (
	"database/sql"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type GroupName struct {
	ID          int64
	VtuberGroup string
}

func Conn() *sql.DB {
	log.Info("Open DB")
	user := os.Getenv("SQLUSER")
	pass := os.Getenv("SQLPASS")
	host := os.Getenv("DBHOST")
	db, err := sql.Open("mysql", ""+user+":"+pass+"@tcp("+host+":3306)/Vtuber?parseTime=true")
	if err != nil {
		log.Error("Something worng with database,make sure you create Vtuber database first")
		log.Error(err)
		os.Exit(1)
	}
	//make sure can access database
	_, err = db.Exec(`SELECT NOW()`)
	if err != nil {
		log.Error("Something worng with database,make sure you create Vtuber database first")
		log.Error(err)
		os.Exit(1)
	}
	return db
}

func GetGroup() []GroupName {
	rows, err := db.Query(`SELECT id,VtuberGroupName FROM VtuberGroup`)
	if err != nil {
		log.Error(err)
	}
	var Data []GroupName
	for rows.Next() {
		var list GroupName
		err = rows.Scan(&list.ID, &list.VtuberGroup)
		if err != nil {
			log.Error(err)
		}
		Data = append(Data, list)
	}
	rows.Close()
	return Data
}

func GetTweetID(limit int64, GroupID int64) []int64 {
	rows, err := db.Query(`SELECT TweetID FROM Vtuber.Twitter inner join VtuberMember on VtuberMember.id = Twitter.VtuberMember_id inner join VtuberGroup on VtuberMember.VtuberGroup_id = VtuberGroup.id where VtuberGroup.id=? order by Twitter.id DESC limit ?`, GroupID, limit)
	if err != nil {
		log.Error(err)
	}
	var Data []int64
	for rows.Next() {
		var tmp string
		err = rows.Scan(&tmp)
		if err != nil {
			log.Error(err)
		}
		ID64, err := strconv.ParseInt(tmp, 10, 64)
		if err != nil {
			log.Error(err)
		}
		Data = append(Data, ID64)
	}
	rows.Close()
	return Data
}
