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
 * AdvantEDGE WLAN Access Information API
 *
 * WLAN Access Information Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC028 WAI API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/028/02.02.01_60/gs_MEC028v020201p.pdf) <p>[Copyright (c) ETSI 2020](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-wais](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-wais) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about WLAN access information in the network <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address/api_ <p>AdvantEDGE supports a selected subset of WAI API subscription types. <p>Supported subscriptions: <p> - AssocStaSubscription <p> - StaDataRateSubscription
 *
 * API version: 2.2.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

type StaCounterTriggerCondition struct {
	// Configure and set threshold for dot11AckFailureCount trigger
	AckFailureCountThreshold int32 `json:"ackFailureCountThreshold,omitempty"`
	// Configure and set threshold for dot11FailedCount trigger
	FailedCountThreshold int32 `json:"failedCountThreshold,omitempty"`
	// Configure and set threshold for dot11FCSErrorCount trigger
	FcsErrorCountThreshold int32 `json:"fcsErrorCountThreshold,omitempty"`
	// Configure and set threshold for dot11FrameDuplicateCount trigger
	FrameDuplicateCountThreshold int32 `json:"frameDuplicateCountThreshold,omitempty"`
	// Configure and set threshold for dot11MultipleRetryCount trigger
	MultipleRetryCountThreshold int32 `json:"multipleRetryCountThreshold,omitempty"`
	// Configure and set threshold for dot11RetryCount trigger
	RetryCountThreshold int32 `json:"retryCountThreshold,omitempty"`
	// Configure and set threshold for dot11RTSFailureCount trigger
	RtsFailureCountThreshold int32 `json:"rtsFailureCountThreshold,omitempty"`
}
