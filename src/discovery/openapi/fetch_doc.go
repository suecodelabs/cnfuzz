// Copyright 2022 Sue B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/suecodelabs/cnfuzz/src/discovery"
	"github.com/suecodelabs/cnfuzz/src/log"
)

const UserAgent = "cnfuzz"

// GetCommonOpenApiLocations returns a list of locations commonly used for OpenAPI specifications in web API's
func GetCommonOpenApiLocations() []string {
	return []string{
		"/swagger/doc.json",
	}
}

// GetRemoteOpenApiDoc get OpenApi doc from a URL
// Use this method if the full URL to the OpenApi doc is known
func GetRemoteOpenApiDoc(url *url.URL) []byte {
	logger := log.L()
	client := http.Client{
		Timeout: time.Second * 10, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		logger.Errorf("error while creating Http request for the OpenAPI doc: %+v", err)
		return nil
	}

	req.Header.Set("User-Agent", UserAgent)

	res, getErr := client.Do(req)
	if getErr != nil {
		logger.Errorf("error while making a Http request for the OpenAPI doc: %+v", getErr)
		return nil
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		logger.Errorf("error while reading the body from Http response: %+v", readErr)
		return nil
	}
	return body
}

// TryGetOpenApiDoc try getting the OpenApi doc from a host without knowing the exact OpenApi doc location
func TryGetOpenApiDoc(ip string, ports []int32, locations []string) (webApiDescription *discovery.WebApiDescription, err error) {
	logger := log.L()
	if len(ports) == 0 {
		baseUri := "http://" + ip
		return tryGetOpenApiDoc(baseUri, locations)
	} else {
		// Try each port
		for _, port := range ports {
			proto := "http://"
			if port == 432 {
				proto = "https://"
			}
			baseUri := proto + ip + ":" + strconv.Itoa(int(port))
			logger.Debugf("trying to get OpenAPI doc from base uri %s ...", baseUri)

			result, err := tryGetOpenApiDoc(baseUri, locations)
			if err != nil {
				// Failed to get the OpenApi doc from this location
				// Check next location
				continue
			} else {
				return result, nil
			}
		}
	}
	return nil, fmt.Errorf("failed to get OpenAPI doc")
}

// tryGetOpenApiDoc attempts to retrieve the OpenAPI doc from the given locations
// continues trying locations until a location is successful or if every location has been tried
func tryGetOpenApiDoc(baseUri string, locations []string) (webApiDescription *discovery.WebApiDescription, err error) {
	logger := log.L()
	// TODO do Api versions

	// Search doc
	client := http.Client{
		Timeout: time.Second * 4,
	}

	for _, try := range locations {
		fullUri := baseUri + try
		logger.Debugf("trying to get OpenAPI doc from location %s ...", fullUri)
		req, err := http.NewRequest(http.MethodGet, fullUri, nil)
		if err != nil {
			logger.Errorf("error while attempting to get the OpenAPI doc: %+v", err)
			continue
		}

		req.Header.Set("User-Agent", UserAgent)

		res, getErr := client.Do(req)
		if getErr != nil {
			logger.Errorf("error while sending Http request to get OpenAPI doc: %+v", getErr)
			continue
		}
		if res.StatusCode == 200 {
			body, readErr := ioutil.ReadAll(res.Body)
			if readErr != nil {
				logger.Errorf("error while reading body from Http response while getting OpenAPI doc: %+v", readErr)
				continue
			}
			if res.Body != nil {
				defer res.Body.Close()
			} else {
				// Body is empty
				continue
			}

			result, err := UnMarshalOpenApiDoc(body, req.URL)
			if err != nil {
				logger.Errorf("error while unmarshalling OpenAPI doc request body: %+v", err)
				continue
			} else {
				// Got the OpenApi Doc :)
				return result, nil
			}
		}
	}

	return nil, fmt.Errorf("failed to get the OpenApi doc from %s", baseUri)
}
