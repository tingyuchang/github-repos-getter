package main

import (
	"fmt"
	"github-repos-getter/api/gitgubapi"
	"github-repos-getter/model"
	"github-repos-getter/setting"
	"log"
	"strconv"
	"sync"
	"time"
)


var wg sync.WaitGroup

func init() {
	setting.Setup()
	model.Setup()
}

func main() {
	defer model.CloseDB()
}

func exec() {
	var datas []model.Repo
	page := 2

	wg.Add(page)
	for i := 0; i < page; i++ {
		time.Sleep(1000 * time.Millisecond)
		go func() {
			defer wg.Done()
			response, err := githubapi.GetGithubRepos("language:JavaScript", "starts", i, 100)
			if err != nil {
				log.Fatal(err)
			}
			if len(response.Items) > 0 {
				datas = append(datas, response.Items...)
			}
		}()
	}

	wg.Wait()

	for _,v := range datas {
		err := model.InsertRepo(v)
		if err != nil {
			if existRepo, err := model.GetRepo(strconv.Itoa(v.Id)); err == nil && existRepo.Id !=  0 {
				err = model.UpdateRepo(v)
				if err != nil {
					log.Fatal(err)
				}
			} else  {
				log.Fatal(err)
			}
		}
	}
}




