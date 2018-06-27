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
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"github.com/HotelsDotCom/flyte-bitbucket/domain"
	"testing"
)

func TestCreateRepoWorkingAsExpected(t *testing.T) {
	// Creates http client, returns 201 and marshals Name: bang into responseBody
	client := bitBucketClient{"host", "user", "pass", func(r *http.Request, responseBody interface{}) (*http.Response, error) {
		mockRepo, err := json.Marshal(domain.Repo{Name: "bang"})
		json.Unmarshal(mockRepo, &responseBody)
		return &http.Response{StatusCode: 201}, err
	}}

	actualRepo, err := client.CreateRepo("foo", "bar")
	expectedRepo := domain.Repo{Name: "bang"}
	if !reflect.DeepEqual(actualRepo, expectedRepo) {
		t.Fatalf("Expected: %v but got: %v", expectedRepo, actualRepo)
	}
	if err != nil {
		t.Fatalf("Expected: nil but got: %s", err)
	}
}

func TestIncorrectStatusCode(t *testing.T) {
	client := bitBucketClient{"host", "user", "pass", func(r *http.Request, responseBody interface{}) (*http.Response, error) {
		mockRepo, err := json.Marshal(domain.Repo{Name: "bang"})
		json.Unmarshal(mockRepo, &responseBody)
		return &http.Response{StatusCode: 404}, err
	}}

	_, actualErr := client.CreateRepo("foo", "bar")
	expectedErr := errors.New("POST failed: path=/rest/api/1.0/projects/foo/repos : repo=bar : statusCode=404")
	if !reflect.DeepEqual(actualErr, expectedErr) {
		t.Fatalf("Expected: %v but got: %v", expectedErr, actualErr)
	}
}
func TestJsonDecodingError(t *testing.T) {
	client := bitBucketClient{"host", "user", "pass", func(r *http.Request, responseBody interface{}) (*http.Response, error) {
		mockRepo, _ := json.Marshal(domain.Repo{Name: "bang"})
		json.Unmarshal(mockRepo, &responseBody)
		return &http.Response{StatusCode: 201}, errors.New("Decoding error")
	}}

	_, actualErr := client.CreateRepo("foo", "bar")
	expectedErr := errors.New("POST failed: path=/rest/api/1.0/projects/foo/repos : repo=bar : err=Decoding error")
	if !reflect.DeepEqual(actualErr, expectedErr) {
		t.Fatalf("Expected: %v but got: %v", expectedErr, actualErr)
	}
}
