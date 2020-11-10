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
 * ETSI GS MEC 013 - Location API
 *
 * The ETSI MEC ISG MEC013 WLAN Access Information API described using OpenAPI.
 *
 * API version: 2.1.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// Structure with attributes relating to the target entity’s velocity, as defined in [14].
type LocationInfoVelocity struct {
	// Bearing, expressed in the range 0° to 360°, as defined in [14].
	Bearing int32 `json:"bearing"`
	// Horizontal speed, expressed in km/h and defined in [14].
	HorizontalSpeed int32 `json:"horizontalSpeed"`
	// Horizontal uncertainty, as defined in [14]. Present only if \"velocityType\" equals 3 or 4
	Uncertainty int32 `json:"uncertainty,omitempty"`
	// Velocity information, as detailed in [14], associated with the reported location coordinate: <p>1 = HORIZONTAL <p>2 = HORIZONTAL_VERTICAL <p>3 = HORIZONTAL_UNCERT <p>4 = HORIZONTAL_VERTICAL_UNCERT
	VelocityType int32 `json:"velocityType"`
	// Vertical speed, expressed in km/h and defined in [14]. Present only if \"velocityType\" equals 2 or 4
	VerticalSpeed int32 `json:"verticalSpeed,omitempty"`
	// Vertical uncertainty, as defined in [14]. Present only if \"velocityType\" equals 4
	VerticalUncertainty int32 `json:"verticalUncertainty,omitempty"`
}
