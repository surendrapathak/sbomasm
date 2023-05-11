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

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const LEANIX_BASE_URL = "https://de-vsm.leanix.net"
const INTERLYNK_API_KEY = "Gz5efs6v6dGuzyvPYtmuajTPB4rEhyDLCcDyGbxn"
const LEANIX_BASE_GRAPHQL_URL = "https://de-vsm.leanix.net/services/vsm-compass/v1/graphql"

func Authenticate() (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", LEANIX_BASE_URL+"/services/mtm/v1/oauth2/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("apitoken", INTERLYNK_API_KEY)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("authentication failed")
	}
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return result["access_token"].(string), nil
}

func GetProduct(bearerToken string, name string) (*GraphQLResponse, error) {
	query := `
query Product {
  products
  {
    name
    id
    relProductToService {
      service {
        name
        id
        description
        createdAt
        sources {
          discoveryItemV2 {
            contributors: data(path: "contributors")
          }
        }
        relServiceToLibrary {
          createdAt
          library {
            createdAt
            description
            group
            id
            licenses
            name
            packageManager
            purl
            updatedAt
            version
            sources {
              discoveryItem: discoveryItemV2 {
                description: data(path: "description")
                version: data(path: "version")
                group: data(path: "group")
                packageManager: data(path: "packageManager")
                contributors: data(path: "contributors")
                licenses: data(path: "licenses")
                time
              }
            }
          }
        }
      }
    }
  }
}	
`

	client := &http.Client{}
	body := map[string]string{"query": query}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling query:", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", LEANIX_BASE_GRAPHQL_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return nil, err
		}

		var responseData GraphQLResponse
		err = json.Unmarshal(bodyBytes, &responseData)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling response data")
		}

		return &responseData, nil

	}
	fmt.Println("Graphql failed with response code:", resp.StatusCode)
	return nil, fmt.Errorf("query failed")
}
