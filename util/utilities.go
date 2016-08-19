/*
* Copyright 2016 Comcast Cable Communications Management, LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package util

import (
	"encoding/json"
	"strings"

	"github.com/cloudfoundry/cli/plugin"
)

//Utilities - interface
type Utilities interface {
	GetDefaultDomain(plugin.CliConnection) string
}

//Utility - struct
type Utility struct{}

//GetDefaultDomain - function to get the default domain from cloud foundry
func (u *Utility) GetDefaultDomain(conn plugin.CliConnection) (domain string) {
	infoArgs := []string{"curl", "/v2/info"}
	infoOutput, err := conn.CliCommandWithoutTerminalOutput(infoArgs...)

	if err != nil {
		domain = "unknown"
		return
	}
	infoBytes := []byte(infoOutput[0])

	var info map[string]interface{}
	err = json.Unmarshal(infoBytes, &info)

	if err == nil {
		apiurl := info["authorization_endpoint"].(string)
		domain = apiurl[strings.IndexRune(apiurl, '.')+1:]
	}
	return
}
