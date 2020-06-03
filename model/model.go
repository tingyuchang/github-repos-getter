package model

import (
	"database/sql"
	"log"
	"time"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github-repos-getter/utils"
	"github-repos-getter/setting"
)

var db *sql.DB
var err error

func Setup() {
	driverName, dataSourceName := loadDBConfig()
	fmt.Println(dataSourceName)
	db, err = sql.Open(driverName, dataSourceName)

	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	defer db.Close()
}

func GetRepo(id string) (Repo, error) {
	rows, err := db.Query(`SELECT * FROM repo WHERE id=?;`, id)

	defer rows.Close()

	for rows.Next() {
		var repo Repo
		//var lastLoginAt sql.NullTime
		var url sql.NullString
		var description sql.NullString
		var language sql.NullString
		if err = rows.Scan(&repo.Id, &repo.Name, &repo.FullName, &url, &description, &language, &repo.StargazersCount, &repo.UpdateAt); err != nil {
			return Repo{}, err
		}

		if url.Valid {
			repo.Url = url.String
		}

		if description.Valid {
			repo.Description = description.String
		}

		if language.Valid {
			repo.Language = language.String
		}

		return repo, nil

	}

	return Repo{}, err
}

func InsertRepo(repo Repo) error {
	stmt, err := db.Prepare("INSERT INTO repo(id,name,full_name,url,description,language,stargazers_count,update_at) VALUES(?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(repo.Id, utils.DivideString(repo.Name, 100), utils.DivideString(repo.FullName, 200), repo.Url,utils.DivideString(repo.Description, 500), repo.Language, repo.StargazersCount, time.Now()); err != nil {
		return err
	}

	return nil
}

func UpdateRepo(repo Repo) error {
	stmt, err := db.Prepare("UPDATE repo SET id=?,name=?,full_name=?,url=?,description=?,language=?,stargazers_count=?,update_at=? WHERE id=?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err := stmt.Exec(repo.Id, utils.DivideString(repo.Name, 100), utils.DivideString(repo.FullName, 200), repo.Url,utils.DivideString(repo.Description, 500), repo.Language, repo.StargazersCount, time.Now(), repo.Id); err != nil {
		return err
	}

	return nil
}

func loadDBConfig() (string, string) {
	connection := setting.Config.DB.Driver
	host := setting.Config.DB.Host
	port := setting.Config.DB.Port
	database := setting.Config.DB.Database
	username := setting.Config.DB.Username
	password := setting.Config.DB.Password
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username,
		password,
		host,
		port,
		database)
	return connection, dataSourceName
}

