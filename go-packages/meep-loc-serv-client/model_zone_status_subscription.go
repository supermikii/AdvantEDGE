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

// A type containing zone status subscription.
type ZoneStatusSubscription struct {
	CallbackReference *CallbackReference `json:"callbackReference"`
	// A correlator that the client can use to tag this particular resource representation during a request to create a resource on the server.
	ClientCorrelator string `json:"clientCorrelator,omitempty"`
	// Threshold number of users in an access point which if crossed shall cause a notification
	NumberOfUsersAPThreshold int32 `json:"numberOfUsersAPThreshold,omitempty"`
	// Threshold number of users in a zone which if crossed shall cause a notification
	NumberOfUsersZoneThreshold int32 `json:"numberOfUsersZoneThreshold,omitempty"`
	// List of operation status values to generate notifications for (these apply to all access points within a zone).
	OperationStatus []OperationStatus `json:"operationStatus,omitempty"`
	// Self referring URL
	ResourceURL string `json:"resourceURL,omitempty"`
	// Identifier of zone
	ZoneId string `json:"zoneId"`
}
