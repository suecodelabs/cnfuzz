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

package serv

import (
	"fmt"
	"net/http"

	"github.com/suecodelabs/cnfuzz/src/log"
)

func livez(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Alive.\n")
}

func readyz(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Ready.\n")
}

// Serv start http server that contains ready and live endpoints
// warning: this function is blocking
func Serv() {

	http.HandleFunc("/livez", livez)
	http.HandleFunc("/readyz", readyz)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.L().Fatal(err)
	}
}
