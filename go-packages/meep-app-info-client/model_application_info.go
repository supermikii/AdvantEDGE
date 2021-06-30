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
 * AdvantEDGE Application Information API
 *
 * AdvantEDGE implementation to create an Application Instance information using OpenAPI. Developed as an extension to Application Enablement API.
 *
 * API version: 1.0.0
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// This type represents the general information of a MEC application.
type ApplicationInfo struct {
	// Application Instance Id
	AppInstanceId string `json:"appInstanceId"`
	// Application Name
	AppName string `json:"appName,omitempty"`
	// Application Version
	Version string            `json:"version,omitempty"`
	State   *ApplicationState `json:"state"`
}
