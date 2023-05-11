// Copyright 2023 Interlynk.io
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package leanix

import "encoding/json"

type GraphQLResponse struct {
	Data Data `json:"data"`
}

type Data struct {
	Products []Product `json:"products"`
}

type Product struct {
	Name                string                `json:"name"`
	ID                  string                `json:"id"`
	RelProductToService []RelProductToService `json:"relProductToService"`
}

type RelProductToService struct {
	Service Service `json:"service"`
}

type Service struct {
	Name                string                `json:"name"`
	ID                  string                `json:"id"`
	Description         string                `json:"description"`
	CreatedAt           string                `json:"createdAt"`
	Sources             []Source              `json:"sources"`
	RelServiceToLibrary []RelServiceToLibrary `json:"relServiceToLibrary"`
}

type Source struct {
	DiscoveryItemV2 DiscoveryItemV2 `json:"discoveryItemV2"`
}

type DiscoveryItemV2 struct {
	Contributors json.RawMessage `json:"contributors"`
}

type RelServiceToLibrary struct {
	CreatedAt string  `json:"createdAt"`
	Library   Library `json:"library"`
}

type Library struct {
	CreatedAt      string          `json:"createdAt"`
	Description    string          `json:"description"`
	Group          string          `json:"group"`
	ID             string          `json:"id"`
	Licenses       json.RawMessage `json:"licenses"`
	Name           string          `json:"name"`
	PackageManager string          `json:"packageManager"`
	Purl           string          `json:"purl"`
	UpdatedAt      string          `json:"updatedAt"`
	Version        string          `json:"version"`
	Sources        []LibrarySource `json:"sources"`
}

type LibrarySource struct {
	DiscoveryItem DiscoveryItem `json:"discoveryItem"`
}

type DiscoveryItem struct {
	Description    json.RawMessage `json:"description"`
	Version        json.RawMessage `json:"version"`
	Group          json.RawMessage `json:"group"`
	PackageManager json.RawMessage `json:"packageManager"`
	Contributors   json.RawMessage `json:"contributors"`
	Time           json.RawMessage `json:"time"`
	Licenses       json.RawMessage `json:"licenses"`
}
