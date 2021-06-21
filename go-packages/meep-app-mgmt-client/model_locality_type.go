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
 * AdvantEDGE MEC Application Management API
 *
 * AdvantEDGE implementation to create a Application Instance information using OpenAPI. Developed as an extension to Application Enablement API.
 *
 * API version: 1.0.0
 * Contact: cti_support@etsi.org
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package client

// LocalityType : The scope of locality as expressed by \"consumedLocalOnly\" and \"isLocal\". If absent, defaults to MEC_HOST
type LocalityType string

// List of LocalityType
const (
	MEC_SYSTEM_LocalityType LocalityType = "MEC_SYSTEM"
	MEC_HOST_LocalityType   LocalityType = "MEC_HOST"
	NFVI_POP_LocalityType   LocalityType = "NFVI_POP"
	ZONE_LocalityType       LocalityType = "ZONE"
	ZONE_GROUP_LocalityType LocalityType = "ZONE_GROUP"
	NFVI_NODE_LocalityType  LocalityType = "NFVI_NODE"
)
