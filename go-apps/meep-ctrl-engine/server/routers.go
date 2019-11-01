/*
 * MEEP Controller REST API
 *
 * Copyright (c) 2019  InterDigital Communications, Inc Licensed under the Apache License, Version 2.0 (the \"License\"); you may not use this file except in compliance with the License. You may obtain a copy of the License at      http://www.apache.org/licenses/LICENSE-2.0  Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an \"AS IS\" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package server

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler = route.HandlerFunc
		//		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/"))))

	router.PathPrefix("/api").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static/api/"))))

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	{
		"Index",
		"GET",
		"/v1/",
		Index,
	},

	{
		"GetStates",
		strings.ToUpper("Get"),
		"/v1/states",
		GetStates,
	},

	{
		"CreateScenario",
		strings.ToUpper("Post"),
		"/v1/scenarios/{name}",
		CreateScenario,
	},

	{
		"DeleteScenario",
		strings.ToUpper("Delete"),
		"/v1/scenarios/{name}",
		DeleteScenario,
	},

	{
		"DeleteScenarioList",
		strings.ToUpper("Delete"),
		"/v1/scenarios",
		DeleteScenarioList,
	},

	{
		"GetScenario",
		strings.ToUpper("Get"),
		"/v1/scenarios/{name}",
		GetScenario,
	},

	{
		"GetScenarioList",
		strings.ToUpper("Get"),
		"/v1/scenarios",
		GetScenarioList,
	},

	{
		"SetScenario",
		strings.ToUpper("Put"),
		"/v1/scenarios/{name}",
		SetScenario,
	},

	{
		"ActivateScenario",
		strings.ToUpper("Post"),
		"/v1/active/{name}",
		ActivateScenario,
	},

	{
		"GetActiveNodeServiceMaps",
		strings.ToUpper("Get"),
		"/v1/active/serviceMaps",
		GetActiveNodeServiceMaps,
	},

	{
		"GetActiveScenario",
		strings.ToUpper("Get"),
		"/v1/active",
		GetActiveScenario,
	},

	{
		"SendEvent",
		strings.ToUpper("Post"),
		"/v1/events/{type}",
		SendEvent,
	},

	{
		"TerminateScenario",
		strings.ToUpper("Delete"),
		"/v1/active",
		TerminateScenario,
	},
}
