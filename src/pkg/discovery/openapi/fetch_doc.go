/*
 * Copyright 2022 Sue B.V.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openapi

import (
	"fmt"
	"github.com/suecodelabs/cnfuzz/src/pkg/discovery"
	"github.com/suecodelabs/cnfuzz/src/pkg/logger"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const UserAgent = "cnfuzz"
const timeout = time.Second * 4

// GetCommonOpenApiLocations returns a list of locations commonly used for OpenAPI specifications in web API's
func GetCommonOpenApiLocations() []string {
	return []string{
		"/swagger/doc.json",
	}
}

// GetRemoteOpenApiDoc get OpenApi doc from a URL
// Use this method if the full URL to the OpenApi doc is known
func GetRemoteOpenApiDoc(l logger.Logger, url *url.URL) ([]byte, error) {
	client := http.Client{
		Timeout: timeout, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		l.V(logger.ImportantLevel).Error(err, "error while creating Http request for the OpenAPI doc")
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)

	res, getErr := client.Do(req)
	if getErr != nil {
		l.V(logger.ImportantLevel).Error(getErr, "error while making a Http request for the OpenAPI doc")
		return nil, getErr
	}

	if res.Body != nil {
		defer res.Body.Close()
	} else {
		return nil, fmt.Errorf("no OpenAPI doc found on the given location")
	}
	if res.StatusCode == 200 {

		body, readErr := io.ReadAll(res.Body)
		if readErr != nil {
			l.V(logger.ImportantLevel).Error(readErr, "error while reading the body from Http response")
			return nil, readErr
		}
		return body, nil
	} else {
		return nil, fmt.Errorf("target returned %d status code when getting the OpenAPI doc", res.StatusCode)
	}
}

// TryGetOpenApiDoc try getting the OpenApi doc from a host without knowing the exact OpenApi doc location
func TryGetOpenApiDoc(l logger.Logger, ip string, ports []int32, locations []string) (webApiDescription *discovery.WebApiDescription, err error) {
	if len(ports) == 0 {
		baseUri := "http://" + ip
		return tryGetOpenApiDoc(l, baseUri, locations)
	} else {
		// Try each port
		for _, port := range ports {
			proto := "http://"
			if port == 432 {
				proto = "https://"
			}
			baseUri := proto + ip + ":" + strconv.Itoa(int(port))
			l.V(logger.DebugLevel).Info("trying to get OpenAPI doc from base uri ...", "docUri", baseUri)

			result, err := tryGetOpenApiDoc(l, baseUri, locations)
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
func tryGetOpenApiDoc(l logger.Logger, baseUri string, locations []string) (webApiDescription *discovery.WebApiDescription, err error) {
	// TODO do Api versions

	for _, try := range locations {
		fullUri, err := url.Parse(baseUri + try)
		if err != nil {
			l.V(logger.InfoLevel).Error(err, "generated URI while attempting to find the OpenAPI doc is invalid")
			continue
		}
		l.V(logger.DebugLevel).Info("trying to get OpenAPI doc from guessed uri ...", "docUri", fullUri)

		body, err := GetRemoteOpenApiDoc(l, fullUri)

		result, err := UnMarshalOpenApiDoc(l, body, fullUri)
		if err != nil {
			l.V(logger.ImportantLevel).Error(err, "error while unmarshalling OpenAPI doc request body")
			continue
		} else {
			// Found the OpenApi Doc :)
			return result, nil
		}
	}

	return nil, fmt.Errorf("failed to get the OpenApi doc from %s", baseUri)
}
