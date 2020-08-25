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

package command

import (
	"encoding/json"
	"errors"
	"github.com/ExpediaGroup/flyte-bitbucket/domain"
	"reflect"
	"testing"
)

func TestCreateRepoAsExpected(t *testing.T) {
	mockCreate := create{func(_, _ string) (domain.Repo, error) {
		repo := domain.Repo{}
		repo.Links.Self = []domain.Self{{"href"}}
		return repo, nil
	}}
	input := toJson(createRepoInput{"foo", "bar"}, t)
	actualEvent := mockCreate.createHandler(input)
	expectedEvent := newCreateRepoEvent("href", "foo", "bar")

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Fatalf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}

func TestErrorCreatingRepo(t *testing.T) {
	mockCreate := create{func(_, _ string) (domain.Repo, error) {
		return domain.Repo{}, errors.New("Some error")
	}}
	input := toJson(createRepoInput{"foo", "bar"}, t)
	actualEvent := mockCreate.createHandler(input)
	expectedEvent := newCreateRepoErrorEvent("Fail: Could not create repo: Some error", "foo", "bar")

	if !reflect.DeepEqual(actualEvent, expectedEvent) {
		t.Fatalf("Expected: %v but got: %v", expectedEvent, actualEvent)
	}
}

func toJson(i interface{}, t *testing.T) []byte {

	b, err := json.Marshal(i)
	if err != nil {
		t.Errorf("error marshalling: %v", err)
	}
	return b
}

type createRepoInput struct {
	Project string `json:"project"`
	Name    string `json:"name"`
}
