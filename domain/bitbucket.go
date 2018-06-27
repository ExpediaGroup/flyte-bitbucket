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

package domain

type Repo struct {
	Slug  string `json:"slug"`
	Name  string `json:"name"`
	Links Links  `json:"links"`
}

type Links struct {
	Self []Self `json:"self"`
}

type Self struct {
	Href string `json:"href"`
}
