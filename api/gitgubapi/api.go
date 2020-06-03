package githubapi

import (
	"encoding/json"
	"fmt"
	"github-repos-getter/model"
	"io/ioutil"
	"net/http"
	"github-repos-getter/setting"
)

const URL = "https://api.github.com"

var client *http.Client

func init() {
	client = &http.Client{}
}

func GetGithubRepos(q,sort string, page, perPage int) (model.Response, error) {
	urlPath := fmt.Sprintf("%v/search/repositories?q=%v&sort=%v&page=%v&per_page=%v", URL, q, sort, page, perPage)
	req, err := http.NewRequest("GET", urlPath, nil)

	if err != nil {
		return model.Response{}, err
	}

	addAuth(req)
	resp, err := client.Do(req)
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return model.Response{}, err
	}

	var response model.Response

	err = json.Unmarshal(robots, &response)
	if err != nil {
		return model.Response{}, err
	}

	return response, nil
}

func GetLimitation() error {
	urlPath := fmt.Sprintf("%v/rate_limit", URL)
	fmt.Println(urlPath)
	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return err
	}

	addAuth(req)
	resp, err := client.Do(req)
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	var data map[string]interface{}

	if err := json.Unmarshal(robots, &data); err != nil {
		return err
	}

	fmt.Println(data)

	return nil
}

func addAuth(req *http.Request) {
	req.Header.Add("Accept", `application/vnd.github.v3+json`)
	req.Header.Add("Authorization", fmt.Sprintf("token %v", setting.Config.APP.GitHubToken))
}