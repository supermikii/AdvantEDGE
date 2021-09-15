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
 * WLAN Access Information Service is AdvantEDGE's implementation of [ETSI MEC ISG MEC028 WAI API](http://www.etsi.org/deliver/etsi_gs/MEC/001_099/028/02.02.01_60/gs_MEC028v020201p.pdf) <p>[Copyright (c) ETSI 2020](https://forge.etsi.org/etsi-forge-copyright-notice.txt) <p>**Micro-service**<br>[meep-wais](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-wais) <p>**Type & Usage**<br>Edge Service used by edge applications that want to get information about WLAN access information in the network <p>**Note**<br>AdvantEDGE supports a selected subset of WAI API subscription types. <p>Supported subscriptions: <p> - AssocStaSubscription <p> - StaDataRateSubscription
 *
 * API version: 2.2.1
 * Contact: AdvantEDGE@InterDigital.com
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

type ChannelLoadConfig struct {
	// Channel for which the channel load report is requested.
	Channel int32 `json:"channel"`
	// Operating Class field indicates an operating class value as defined in Annex E within IEEE 802.11-2016 [8].
	OperatingClass int32 `json:"operatingClass"`
	// Reporting condition for the Beacon Report as per Table 9-153 of IEEE 802.11-2016 [8]: 0 = Report to be issued after each measurement. 1 = Report to be issued when Channel Load is greater than or equal to the threshold. 2 = Report to be issued when Channel Load is less than or equal to the threshold.  If this optional field is not provided, channel load report should be issued after each measurement (reportingCondition = 0).
	ReportingCondition int32 `json:"reportingCondition,omitempty"`
	// Channel Load reference value for threshold reporting. This field shall be provided for reportingCondition values 1 and 2.
	Threshold int32 `json:"threshold,omitempty"`
}
