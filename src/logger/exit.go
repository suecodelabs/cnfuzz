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

package logger

// Inspired by
// https://stackoverflow.com/a/30690532

import "os"

// Func takes a code as exit status
type Func func(int)

// Exit has an exit func, and will memorize the exit status code
type Exit struct {
	exit   Func
	status int
}

// Exit calls the exiter, and then returns code as status.
// If e was declared, but never set (since only a test would set e),
// simply calls os.Exit()
func (e *Exit) Exit(code int) {
	if e != nil {
		e.status = code
		e.exit(code)
	} else {
		os.Exit(code)
	}
}

// Status get the exit status code as memorized
// after the call to the exit func.
func (e *Exit) Status() int {
	return e.status
}

// Default returns an Exit with default os.Exit() call.
// That means the status will never be visible,
// since os.Exit() stops everything.
func Default() *Exit {
	return &Exit{exit: os.Exit}
}

// CreateExiter returns an exiter with a custom function
func CreateExiter(exit Func) *Exit {
	return &Exit{exit: exit}
}
