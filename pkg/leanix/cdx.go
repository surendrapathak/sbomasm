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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	cydx "github.com/CycloneDX/cyclonedx-go"
	"github.com/google/uuid"
	"sigs.k8s.io/release-utils/version"
)

func WriteTempCDXFiles(productData *GraphQLResponse) ([]string, error) {
	var cdxFiles []string
	for _, product := range productData.Data.Products {
		fmt.Println("Writing product: ", product.Name)
		for _, service := range product.RelProductToService {
			fmt.Println("Writing service: ", service.Service.Name)
			f, _ := writeService(service.Service)
			cdxFiles = append(cdxFiles, f)
		}
	}
	return cdxFiles, nil

}

func writeService(service Service) (string, error) {
	bom := cydx.NewBOM()
	bom.SerialNumber = newSerialNumber()
	bom.Version = 1

	meta := &cydx.Metadata{}
	meta.Timestamp = utcNowTime()

	meta.Component = &cydx.Component{
		Type:        "application",
		Name:        service.Name,
		Description: service.Description,
	}
	meta.Tools = &[]cydx.Tool{
		{Vendor: "Interlynk.io", Name: "sbomasm", Version: version.GetVersionInfo().GitVersion},
	}

	meta.Licenses = &cydx.Licenses{
		{License: &cydx.License{ID: "CC-BY-1.0"}},
	}

	meta.Authors = &[]cydx.OrganizationalContact{}

	for _, source := range service.Sources {
		if source.DiscoveryItemV2.Contributors != nil {
			*meta.Authors = append(*meta.Authors, cydx.OrganizationalContact{
				Email: string(source.DiscoveryItemV2.Contributors)})
		}
	}

	bom.Metadata = meta
	bom.Components = &[]cydx.Component{}

	for _, rel := range service.RelServiceToLibrary {
		lib := rel.Library
		comp := cydx.Component{}
		comp.Type = "library"
		comp.Name = lib.Name
		comp.Version = lib.Version
		comp.Description = lib.Description
		comp.PackageURL = lib.Purl
		for _, source := range lib.Sources {
			item := source.DiscoveryItem
			if comp.Description == "" {
				var desc string
				json.Unmarshal([]byte(item.Description), &desc)
				comp.Description = desc
			}
			if comp.Version == "" {
				var ver string
				json.Unmarshal([]byte(item.Version), &ver)
				comp.Version = ver
			}
			if comp.Licenses == nil {
				comp.Licenses = &cydx.Licenses{}
				inLics := []string{}
				json.Unmarshal([]byte(item.Licenses), &inLics)
				for _, lic := range inLics {
					l := &cydx.License{ID: string(lic)}
					*comp.Licenses = append(*comp.Licenses, cydx.LicenseChoice{License: l})
				}
			}

		}
		*bom.Components = append(*bom.Components, comp)
	}

	file, err := ioutil.TempFile("./tmp/", "leanix-cdx-*.json")
	if err != nil {
		return "", err
	}

	encoder := cydx.NewBOMEncoder(file, cydx.BOMFileFormatJSON)
	encoder.SetPretty(true)

	encoder.SetPretty(true)
	encoder.SetEscapeHTML(true)
	if err := encoder.Encode(bom); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func utcNowTime() string {
	location, _ := time.LoadLocation("UTC")
	locationTime := time.Now().In(location)
	return locationTime.Format(time.RFC3339)
}

func newSerialNumber() string {
	u := uuid.New().String()

	return fmt.Sprintf("urn:uuid:%s", u)
}
