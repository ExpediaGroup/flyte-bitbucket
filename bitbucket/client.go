/*
Copyright (C) 2018 Expedia Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ExpediaGroup/flyte-bitbucket/domain"
	"log"
	"net/http"
	"os"
)

type httpSendRequest func(r *http.Request, responseBody interface{}) (*http.Response, error)

type bitBucketClient struct {
	Host       string
	User       string
	Pass       string
	httpClient httpSendRequest
}

type Client interface {
	CreateRepo(string, string) (domain.Repo, error)
}

func NewClient() Client {
	host := getEnv("BITBUCKET_HOST")
	user := getEnv("BITBUCKET_USER")
	pass := getEnv("BITBUCKET_PASS")

	return &bitBucketClient{host, user, pass, sendRequest}
}

func (c bitBucketClient) CreateRepo(project, newRepo string) (domain.Repo, error) {
	var repo domain.Repo
	var name = struct {
		Name string `json:"name"`
	}{Name: newRepo}

	repoName, err := json.Marshal(name)
	if err != nil {
		return repo, err
	}

	path := fmt.Sprintf("/rest/api/1.0/projects/%s/repos", project)
	request, err := c.constructPostRequest(path, string(repoName))
	if err != nil {
		return repo, err
	}

	response, err := c.httpClient(request, &repo)
	if err != nil {
		err = fmt.Errorf("POST failed: path=%s : repo=%s : err=%v", path, name.Name, err)
		return repo, err
	}

	if response.StatusCode != http.StatusCreated {
		err = fmt.Errorf("POST failed: path=%s : repo=%s : statusCode=%v", path, name.Name, response.StatusCode)
		return repo, err
	}
	return repo, err
}

func (c bitBucketClient) constructPostRequest(path, data string) (*http.Request, error) {
	urlStr := fmt.Sprintf("%s%s", c.Host, path)
	r, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return r, err
	}

	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.SetBasicAuth(c.User, c.Pass)

	return r, err
}

func sendRequest(r *http.Request, responseBody interface{}) (*http.Response, error) {
	response, err := http.DefaultClient.Do(r)
	if err != nil {
		return response, err
	}

	defer response.Body.Close()
	return response, json.NewDecoder(response.Body).Decode(&responseBody)
}

func getEnv(env string) string {
	value := os.Getenv(env)
	if value == "" {
		log.Fatalf("%s environment variable is not set", env)
	}
	return value
}
