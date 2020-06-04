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

	// exec()
	exec2()
}

func exec() {
	var repos []model.Repo
	page := 10
	ch := make(chan model.Response, page)
	for i := 0; i < page; i++ {
		time.Sleep(400 * time.Millisecond)
		wg.Add(1)
		exit := func(){wg.Done()}
		go githubapi.GetGithubRepos("language:PHP", "starts", i, 2, ch, exit )
	}
	go func() {
		defer close(ch)
		wg.Wait()
	}()

	for response := range ch {
		if len(response.Items) > 0 {
			repos = append(repos, response.Items...)
		}
	}

	for _,v := range repos {
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

func exec2() {
	var repos []model.Repo
	page := 10
	mu := sync.Mutex{}
	wg.Add(page)
	for i := 0; i < page; i++ {
		i := i
		time.Sleep(400 * time.Millisecond)
		go func() {
			defer wg.Done()
			response, err := githubapi.GetGithubRepos2("language:JavaScript", "starts", i, 2)
			if err != nil {
				log.Fatal(err)
			}
			if len(response.Items) > 0 {
				mu.Lock()
				repos = append(repos, response.Items...)
				mu.Unlock()
			}
		}()
	}

	wg.Wait()

	fmt.Println(len(repos))
}




