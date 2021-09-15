/*
 * Copyright (c) 2020  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the \"License\");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an \"AS IS\" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * AdvantEDGE Application Mobility API
 *
 * Application Mobility Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC021 Application Mobility API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/021/02.01.01_60/gs_MEC021v020101p.pdf) <p>[Copyright (c) ETSI 2017](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-ams](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-ams) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about application mobility in the network <p>**Note**<br>AdvantEDGE supports all of Application Mobility API endpoints (see below).
 *
 * API version: 2.1.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

import (
	"fmt"
	"net/http"
	"strings"

	httpLog "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-http-logger"
	met "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-metrics"

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
	var handler http.Handler
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)
		handler = met.MetricsHandler(handler, sandboxName, serviceName)
		handler = httpLog.LogRx(handler, "")
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	// Path prefix router order is important
	// Service Api files
	handler = http.StripPrefix("/amsi/v1/api/", http.FileServer(http.Dir("./api/")))
	router.
		PathPrefix("/amsi/v1/api/").
		Name("Api").
		Handler(handler)
	// User supplied service API files
	handler = http.StripPrefix("/amsi/v1/user-api/", http.FileServer(http.Dir("./user-api/")))
	router.
		PathPrefix("/amsi/v1/user-api/").
		Name("UserApi").
		Handler(handler)

	return router
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/amsi/v1/",
		Index,
	},

	Route{
		"AdjAppInstGET",
		strings.ToUpper("Get"),
		"/amsi/v1/queries/adjacent_app_instances",
		AdjAppInstGET,
	},

	Route{
		"AppMobilityServiceByIdDELETE",
		strings.ToUpper("Delete"),
		"/amsi/v1/app_mobility_services/{appMobilityServiceId}",
		AppMobilityServiceByIdDELETE,
	},

	Route{
		"AppMobilityServiceByIdGET",
		strings.ToUpper("Get"),
		"/amsi/v1/app_mobility_services/{appMobilityServiceId}",
		AppMobilityServiceByIdGET,
	},

	Route{
		"AppMobilityServiceByIdPUT",
		strings.ToUpper("Put"),
		"/amsi/v1/app_mobility_services/{appMobilityServiceId}",
		AppMobilityServiceByIdPUT,
	},

	Route{
		"AppMobilityServiceGET",
		strings.ToUpper("Get"),
		"/amsi/v1/app_mobility_services",
		AppMobilityServiceGET,
	},

	Route{
		"AppMobilityServicePOST",
		strings.ToUpper("Post"),
		"/amsi/v1/app_mobility_services",
		AppMobilityServicePOST,
	},

	Route{
		"AppMobilityServiceDerPOST",
		strings.ToUpper("Post"),
		"/amsi/v1/app_mobility_services/{appMobilityServiceId}/deregister_task",
		AppMobilityServiceDerPOST,
	},

	Route{
		"Mec011AppTerminationPOST",
		strings.ToUpper("Post"),
		"/amsi/v1/notifications/mec011/appTermination",
		Mec011AppTerminationPOST,
	},

	Route{
		"SubByIdDELETE",
		strings.ToUpper("Delete"),
		"/amsi/v1/subscriptions/{subscriptionId}",
		SubByIdDELETE,
	},

	Route{
		"SubByIdGET",
		strings.ToUpper("Get"),
		"/amsi/v1/subscriptions/{subscriptionId}",
		SubByIdGET,
	},

	Route{
		"SubByIdPUT",
		strings.ToUpper("Put"),
		"/amsi/v1/subscriptions/{subscriptionId}",
		SubByIdPUT,
	},

	Route{
		"SubGET",
		strings.ToUpper("Get"),
		"/amsi/v1/subscriptions/",
		SubGET,
	},

	Route{
		"SubPOST",
		strings.ToUpper("Post"),
		"/amsi/v1/subscriptions/",
		SubPOST,
	},
}
