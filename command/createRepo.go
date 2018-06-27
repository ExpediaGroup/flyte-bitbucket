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
	"fmt"
	"log"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"github.com/HotelsDotCom/flyte-bitbucket/domain"
	"github.com/HotelsDotCom/flyte-bitbucket/bitbucket"
)

var CreateRepoCommand = flyte.Command{
	Name:         "createRepo",
	OutputEvents: []flyte.EventDef{createRepoEventDef, createRepoFailEventDef},
	Handler:      NewCreateRepo().createHandler,
}

type createRepoFunc func(string, string) (domain.Repo, error)

type create struct {
	createRepoFunc createRepoFunc
}

func NewCreateRepo() create {
	return create{createRepo}
}

func (c create) createHandler(input json.RawMessage) flyte.Event {
	var handlerInput struct {
		Project string `json:"project"`
		Name    string `json:"name"`
	}
	// TODO change this to a fatal event once it is implemented
	if err := json.Unmarshal(input, &handlerInput); err != nil {
		err = fmt.Errorf("Could not marshal Project & name into json: %v", err)
		log.Println(err)
		return newCreateRepoErrorEvent(fmt.Sprintf("Fail: %s", err), "unknown", "unknown")
	}

	repo, err := c.createRepoFunc(handlerInput.Project, handlerInput.Name)
	if err != nil {
		err = fmt.Errorf("Could not create repo: %v", err)
		log.Println(err)
		return newCreateRepoErrorEvent(fmt.Sprintf("Fail: %s", err), handlerInput.Project, handlerInput.Name)
	}
	return newCreateRepoEvent(repo.Links.Self[0].Href, handlerInput.Project, handlerInput.Name)
}

func createRepo(project, repoName string) (domain.Repo, error) {
	client := bitbucket.NewClient()
	repo, err := client.CreateRepo(project, repoName)
	if err != nil {
		return domain.Repo{}, err
	}
	return repo, nil
}

var createRepoEventDef = flyte.EventDef{
	Name: "createRepo",
}

type createRepoSuccessPayload struct {
	Project string `json:"project"`
	Name    string `json:"name"`
	Url     string `json:"url"`
}

var createRepoFailEventDef = flyte.EventDef{
	Name: "createRepoFailure",
}

type createRepoFailurePayload struct {
	Project string `json:"project"`
	Name    string `json:"name"`
	Error   string `json:"error"`
}

func newCreateRepoEvent(url, project, name string) flyte.Event {
	return flyte.Event{
		EventDef: createRepoEventDef,
		Payload:  createRepoSuccessPayload{project, name, url},
	}
}

func newCreateRepoErrorEvent(error, project, name string) flyte.Event {
	return flyte.Event{
		EventDef: createRepoFailEventDef,
		Payload:  createRepoFailurePayload{project, name, error},
	}
}
