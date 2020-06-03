package model

import "time"

type Repo struct {
	Id int `json:"id"`
	Name string `json:"name"`
	FullName string `json:"full_name"`
	Url string `json:"html_url"`
	Description string `json:"description"`
	Language string `json:"language"`
	StargazersCount int `json:"stargazers_count"`
	UpdateAt time.Time
}