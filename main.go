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

package main

import (
	"github.com/ExpediaGroup/flyte-bitbucket/command"
	"github.com/HotelsDotCom/flyte-client/client"
	"github.com/HotelsDotCom/flyte-client/flyte"
	"log"
	"net/url"
	"os"
	"time"
)

func main() {
	flyteHost := os.Getenv("FLYTE_API_URL")
	if flyteHost == "" {
		log.Fatal("FLYTE_API_URL environment variable is not set")
	}
	hostUrl := getUrl(flyteHost)
	packDef := flyte.PackDef{
		Name:     "BitBucket",
		HelpURL:  getUrl("https://github.com/ExpediaGroup/flyte-bitbucket/blob/master/README.md"),
		Commands: []flyte.Command{command.CreateRepoCommand},
	}

	p := flyte.NewPack(packDef, client.NewClient(hostUrl, 10*time.Second))
	p.Start()
	// Waits indefinitely
	select {}
}

func getUrl(rawUrl string) *url.URL {
	url, err := url.Parse(rawUrl)
	if err != nil {
		log.Fatalf("%s is not a valid url", rawUrl)
	}
	return url
}
