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
 * AdvantEDGE MEC Service Management API
 *
 * MEC Service Management API describe is AdvantEDGE's implementation of [ETSI MEC ISG MEC011 MEC Service Management API](https://www.etsi.org/deliver/etsi_gs/MEC/001_099/011/02.01.01_60/gs_MEC011v020101p.pdf)
 *
 * API version: 2.1.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

type ProblemDetails struct {
	// A URI reference according to IETF RFC 3986 that identifies the problem type
	Type_ string `json:"type,omitempty"`
	// A short, human-readable summary of the problem type
	Title string `json:"title,omitempty"`
	// The HTTP status code for this occurrence of the problem
	Status int32 `json:"status,omitempty"`
	// A human-readable explanation specific to this occurrence of the problem
	Detail string `json:"detail,omitempty"`
	// A URI reference that identifies the specific occurrence of the problem
	Instance string `json:"instance,omitempty"`
}
