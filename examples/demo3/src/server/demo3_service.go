/*
 * Copyright (c) 2021  InterDigital Communications, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/InterDigitalInc/AdvantEDGE/example/demo3/src/util"
	ams "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-ams-client"
	asc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-app-support-client"
	log "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-logger"
	smc "github.com/InterDigitalInc/AdvantEDGE/go-packages/meep-service-mgmt-client"
)

var mutex sync.Mutex

// App-enablement client
var srvMgmtClient *smc.APIClient
var srvMgmtClientPath string
var appSupportClient *asc.APIClient
var appSupportClientPath string

// Ams client & payload
var amsClient *ams.APIClient
var amsServiceId string
var amsTargetId string
var contextState ContextState

// Demo 3 edge-case handling
var subscriptionSent bool
var confirmReady bool
var registeredService bool
var amsSubscriptionSent bool
var amsServiceCreated bool

// Config attributes
var instanceName string
var mecUrl string
var localPort string
var local string
var mep string

// Demo 3 discovered services & create service
var mecServices = make(map[string]string)
var serviceName string = "user-app"
var scopeOfLocality string = defaultScopeOfLocality
var consumedLocalOnly bool = defaultConsumedLocalOnly

const serviceAppVersion = "2.1.1"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true

// Demo 3 termination handling
var amsSubscriptionId string
var appEnablementServiceId string
var serviceSubscriptionId string
var terminationSubscriptionId string
var terminationSubscription bool = false
var terminated bool = false
var terminateNotification bool = false

type ContextState struct {
	Counter int    `json:"counter"`
	AppId   string `json:"appId,omitempty"`
	Mep     string `json:"mep,omitempty"`
}

func IncrementCounter() {
	contextState.Counter++
}

func Init(envPath string, envName string) (port string, err error) {
	// Start counter & initalize context state for ams
	contextState = ContextState{
		Counter: 0,
	}

	var config util.Config
	var configErr error

	log.Info("Using config values from ", envPath, "/", envName)
	config, configErr = util.LoadConfig(envPath, envName)

	if configErr != nil {
		log.Fatal(configErr)
	}

	// Retrieve local url from config
	local = config.Localurl

	// Retrieve app id from config
	instanceName = config.AppInstanceId
	contextState.AppId = instanceName

	// Retrieve sandbox url from config
	mecUrl = config.SandboxUrl

	// Find mec platform mec app is on
	resp := strings.LastIndex(mecUrl, "/")
	if resp == -1 {
		log.Error("Error finding mec platform")
	} else {
		mep = mecUrl[resp+1:]
	}
	contextState.Mep = mep

	// Retreieve local url from config
	localPort = config.Port

	// Retrieve service name config otherwise use default service name
	if config.ServiceName != "" {
		serviceName = config.ServiceName
	}

	log.Info("Starting Demo 3 instance on Port=", localPort, " using app instance id=", instanceName, " mec platform=", mep)

	// Create application support client
	appSupportClientCfg := asc.NewConfiguration()
	appSupportClientCfg.BasePath = mecUrl + "/mec_app_support/v1"
	appSupportClient = asc.NewAPIClient(appSupportClientCfg)
	appSupportClientPath = appSupportClientCfg.BasePath
	if appSupportClient == nil {
		return "", errors.New("Failed to create App Enablement App Support REST API client")
	}

	// Create service management client
	srvMgmtClientCfg := smc.NewConfiguration()
	srvMgmtClientCfg.BasePath = mecUrl + "/mec_service_mgmt/v1"
	srvMgmtClient = smc.NewAPIClient(srvMgmtClientCfg)
	srvMgmtClientPath = srvMgmtClientCfg.BasePath
	if srvMgmtClient == nil {
		return "", errors.New("Failed to create App Enablement Service Management REST API client")
	}

	return localPort, nil
}

// REST API
// Discover mec services & subscribe to service availilable subscription
func servicesSubscriptionPOST(w http.ResponseWriter, r *http.Request) {

	// Retrieving mec services
	err := getMecServices()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Check subscription if sent to prevent resending subscription
	if !subscriptionSent {
		callBackReference := local + localPort + "/services/callback/service-availability"
		err := subscribeAvailability(instanceName, callBackReference)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		subscriptionSent = true
	}

	// Send response
	w.WriteHeader(http.StatusOK)
}

// REST API
// Handle service subscription callback notification
func notificationPOST(w http.ResponseWriter, r *http.Request) {

	// Decode request body
	var notification smc.ServiceAvailabilityNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&notification)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("Received service availability notification")

	// Parse request param to show on logs
	msg := ""
	if notification.ServiceReferences[0].ChangeType == "ADDED" {
		msg = "Available"
	} else {
		msg = "Unavailable"
	}

	state := ""
	if *notification.ServiceReferences[0].State == smc.ACTIVE_ServiceState {
		state = "ACTIVE"
	} else {
		state = "UNACTIVE"
	}
	log.Info(notification.ServiceReferences[0].SerName + " " + msg + " (" + state + ")")

	w.WriteHeader(http.StatusOK)
}

// Rest API
// Create mec service only if none created
func servicePOST(w http.ResponseWriter, r *http.Request) {

	// Lock registered service to prevent creating more than one mec service from multiple client concurrently
	mutex.Lock()
	defer mutex.Unlock()
	if !registeredService {
		err := registerService(instanceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		registeredService = true
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully created a service")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Service already created")
}

// Rest API
// Delete mec service only if present
func serviceDELETE(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	defer mutex.Unlock()

	if registeredService {
		err := unregisterService(instanceName, appEnablementServiceId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		registeredService = false
		log.Info(serviceName, " service deleted")
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Fprintf(w, "Need to create a service first")
}

// Rest API
// Handle user-app termination call-back notification
func terminateNotificatonPOST(w http.ResponseWriter, r *http.Request) {
	var notification asc.AppTerminationNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&notification)
	if err != nil {
		log.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Info("Received user-app termination notification")
	w.WriteHeader(http.StatusOK)
	terminateNotification = true
	Terminate()
}

// Rest API
// Register MEC Application instances with AMS & consume servicee
func amsCreatePOST(w http.ResponseWriter, r *http.Request) {

	// Cofigure AMS mec client
	// Create application mobility suppport client
	if !amsServiceCreated {

		amsClientcfg := ams.NewConfiguration()
		amsUrl := mecServices["mec021-1"]
		if amsUrl == "" {
			log.Info("Could not find ams services try discovering available services ")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Could not find ams services try discovering available services")
			return
		}
		amsClientcfg.BasePath = amsUrl
		amsClient = ams.NewAPIClient(amsClientcfg)
		if amsClient == nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Invoke client
		err := amsSendService(instanceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		amsServiceCreated = true
		w.WriteHeader(http.StatusOK)
		return
	}
	fmt.Fprintf(w, "AMS service created already")
	w.WriteHeader(http.StatusOK)
}

// Rest API
// Submit AMS subscription to mec platform
func amsSubscriptionPOST(w http.ResponseWriter, r *http.Request) {

	if !amsSubscriptionSent && amsServiceCreated {
		err := amsSendSubscription(instanceName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Need to create a service or already have a subscription")
}

// Rest API
// Handle AMS notification
func amsNotificationPOST(w http.ResponseWriter, r *http.Request) {
	var amsNotification ams.MobilityProcedureNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&amsNotification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	amsTargetId = amsNotification.TargetAppInfo.AppInstanceId

	log.Info("AMS event received for ", amsNotification.AssociateId[0].Value, " moved to app ", amsTargetId)

	// Find ams target service resource url using mec011
	serviceInfo, _, serviceErr := srvMgmtClient.MecServiceMgmtApi.AppServicesGET(context.TODO(), amsTargetId, nil)
	if serviceErr != nil {
		log.Debug("Failed to get target app mec service resource on mec platform")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var notifyUrl string
	for i := 0; i < len(serviceInfo); i++ {
		if serviceInfo[i].SerName == serviceName {
			notifyUrl = serviceInfo[i].TransportInfo.Endpoint.Uris[0]
		}
	}

	sendContextTransfer(notifyUrl, amsNotification.TargetAppInfo.AppInstanceId)

	// Update ams
	amsErr := amsUpdate(amsServiceId, instanceName, amsTargetId, 1, true)
	if amsErr != nil {
		log.Error("Failed to update ams")
	}

	w.WriteHeader(http.StatusOK)
}

// Rest API
// Handle context state transfer
func stateTransferPOST(w http.ResponseWriter, r *http.Request) {
	var targetContextState ContextState
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&targetContextState)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update AMS
	amsErr := amsUpdate(amsServiceId, instanceName, targetContextState.AppId, 0, false)
	if amsErr != nil {
		log.Info("Failed to update ams")
	}

	log.Info("AMS context info counter = ", targetContextState.Counter, " received from user app ", targetContextState.AppId, " on ", targetContextState.Mep)
	w.WriteHeader(http.StatusOK)
}

// Client request to sent context state transfer
func sendContextTransfer(notifyUrl string, targetId string) {
	log.Info("Sending context state counter = ", contextState.Counter, " to user app ", targetId)
	jsonCounter, err := json.Marshal(contextState)
	if err != nil {
		log.Error("Failed to marshal context state ", err.Error())
	}

	resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(jsonCounter))
	log.Info(notifyUrl)
	if err != nil {
		log.Error(resp.Status, err)
		return
	}

	defer resp.Body.Close()

}

// CLient request to create a new application mobility service
func amsSendService(appInstanceId string) error {
	log.Debug("Sending request to mec platform create ams api")
	var bodyRegisterationInfo ams.RegistrationInfo
	bodyRegisterationInfo.ServiceConsumerId = &ams.RegistrationInfoServiceConsumerId{
		AppInstanceId: appInstanceId,
	}

	var associateId ams.AssociateId
	associateId.Type_ = 1
	associateId.Value = "10.100.0.3"
	bodyRegisterationInfo.DeviceInformation = append(bodyRegisterationInfo.DeviceInformation, ams.RegistrationInfoDeviceInformation{AssociateId: &associateId,
		AppMobilityServiceLevel: 3,
		ContextTransferState:    0,
	})

	registerationInfo, _, err := amsClient.AmsiApi.AppMobilityServicePOST(context.TODO(), bodyRegisterationInfo)

	// Store ams service id
	amsServiceId = registerationInfo.AppMobilityServiceId

	if err != nil {
		log.Error(err)
		return err
	} else {
		log.Info("Created app mobility service resource on user app instance ", instanceName[0:6], "...", " tracking ", associateId.Value)
	}

	return nil
}

// Client request to update device context transfer state
func amsUpdate(amsId string, appInstanceId string, targetId string, contextState int32, leaving bool) error {
	var bodyRegisterationInfo ams.RegistrationInfo
	bodyRegisterationInfo.AppMobilityServiceId = amsId
	bodyRegisterationInfo.ServiceConsumerId = &ams.RegistrationInfoServiceConsumerId{
		AppInstanceId: appInstanceId,
	}

	// Provide device info only specific to mep1 platform
	var associateId ams.AssociateId
	associateId.Type_ = 1
	associateId.Value = "10.100.0.3"
	bodyRegisterationInfo.DeviceInformation = append(bodyRegisterationInfo.DeviceInformation, ams.RegistrationInfoDeviceInformation{AssociateId: &associateId,
		AppMobilityServiceLevel: 3,
		ContextTransferState:    contextState,
	})

	_, _, err := amsClient.AmsiApi.AppMobilityServiceByIdPUT(context.TODO(), bodyRegisterationInfo, amsId)
	if err != nil {
		log.Error(err)
		return err
	}

	if leaving {
		log.Info("Completed AMS context transfer on ", instanceName[0:6], "... to user-app ", targetId[0:6], "...")
	} else {
		log.Info("Completed AMS context transfer from ", targetId[0:6], "...")

	}
	return nil
}

// CLient request to create an ams subscription
func amsSendSubscription(appInstanceId string) error {
	log.Debug("Sending request to mec platform adding ams subscription api")

	var mobilityProcedureSubscription ams.MobilityProcedureSubscription

	// Add body param callback ref
	mobilityProcedureSubscription.CallbackReference = local + localPort + "/services/callback/amsevent"
	mobilityProcedureSubscription.SubscriptionType = "MobilityProcedureSubscription"

	// Default tracking ue set to 10.100.0.3
	var associateId ams.AssociateId
	associateId.Type_ = 1
	associateId.Value = "10.100.0.3"

	// Filter criteria
	var mobilityFiler ams.MobilityProcedureSubscriptionFilterCriteria
	mobilityFiler.AppInstanceId = appInstanceId
	mobilityFiler.AssociateId = append(mobilityFiler.AssociateId, associateId)

	mobilityProcedureSubscription.FilterCriteria = &mobilityFiler

	inlineSubscription := ams.ConvertMobilityProcedureSubscriptionToInlineSubscription(&mobilityProcedureSubscription)

	mobilitySubscription, resp, err := amsClient.AmsiApi.SubPOST(context.TODO(), *inlineSubscription)
	hRefLink := mobilitySubscription.Links.Self.Href

	// Find subscription id from response
	idPosition := strings.LastIndex(hRefLink, "/")
	if idPosition == -1 {
		log.Error("Error parsing subscription id for subscription")
		return err
	}
	amsSubscriptionId = hRefLink[idPosition+1:]

	if err != nil {
		log.Error(resp.Status)
		return err
	} else {
		amsSubscriptionSent = true
		log.Info("Successfully created ams subscription")
	}
	return nil
}

// Client request to notify mec platform of mec app
func sendReadyConfirmation(appInstanceId string) error {
	log.Debug("Sending request to mec platform user-app confirm-ready api")
	var appReady asc.AppReadyConfirmation
	appReady.Indication = "READY"
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsConfirmReadyPOST(context.TODO(), appReady, appInstanceId)
	if err != nil {
		log.Error("Failed to send ready confirm acknowlegement: ", resp.Status)
		return err
	}
	return nil
}

// Client request to retrieve list of mec service resources on sandbox
func getMecServices() error {
	log.Debug("Sending request to mec platform get service resources api ")
	appServicesResponse, resp, err := srvMgmtClient.MecServiceMgmtApi.ServicesGET(context.TODO(), nil)
	if err != nil {
		log.Error("Failed to fetch services on mec platform ", resp.Status)
		return err
	}

	log.Info("Returning available mec service resources on mec platform")

	// Store mec services & log service urls
	for i := 0; i < len(appServicesResponse); i++ {
		mecServices[appServicesResponse[i].SerName] = appServicesResponse[i].TransportInfo.Endpoint.Uris[0]
		log.Info(appServicesResponse[i].SerName, " URL: ", appServicesResponse[i].TransportInfo.Endpoint.Uris[0])
	}

	return nil
}

// Client request to create a mec-service resource
func registerService(appInstanceId string) error {
	log.Debug("Sending request to mec platform post service resource api ")
	var srvInfo smc.ServiceInfoPost
	//serName
	srvInfo.SerName = serviceName
	//version
	srvInfo.Version = serviceAppVersion
	//state
	state := smc.ACTIVE_ServiceState
	srvInfo.State = &state
	//serializer
	serializer := smc.JSON_SerializerType
	srvInfo.Serializer = &serializer

	//transportInfo
	var transportInfo smc.TransportInfo
	transportInfo.Id = "transport"
	transportInfo.Name = "REST"
	transportType := smc.REST_HTTP_TransportType
	transportInfo.Type_ = &transportType
	transportInfo.Protocol = "HTTP"
	transportInfo.Version = "2.0"
	var endpoint smc.OneOfTransportInfoEndpoint
	endpointPath := local + localPort + "/services/callback/incoming-context"
	endpoint.Uris = append(endpoint.Uris, endpointPath)
	transportInfo.Endpoint = &endpoint
	srvInfo.TransportInfo = &transportInfo

	//serCategory
	var category smc.CategoryRef
	category.Href = "catalogueHref"
	category.Id = "amsId"
	category.Name = "AMSI"
	category.Version = "v1"
	srvInfo.SerCategory = &category

	//scopeOfLocality
	localityType := smc.LocalityType(scopeOfLocality)
	srvInfo.ScopeOfLocality = &localityType

	//consumedLocalOnly
	srvInfo.ConsumedLocalOnly = consumedLocalOnly

	appServicesPostResponse, resp, err := srvMgmtClient.MecServiceMgmtApi.AppServicesPOST(context.TODO(), srvInfo, appInstanceId)
	if err != nil {
		log.Error("Failed to register service resource on mec app enablement registry: ", resp.Status)
		return err
	}
	log.Info(serviceName, " service resource created with instance id: ", appServicesPostResponse.SerInstanceId)
	appEnablementServiceId = appServicesPostResponse.SerInstanceId
	registeredService = true
	return nil
}

// Client request to delete a mec-service resource
func unregisterService(appInstanceId string, serviceId string) error {
	//log.Debug("Sending request to mec platform delete service api")
	resp, err := srvMgmtClient.MecServiceMgmtApi.AppServicesServiceIdDELETE(context.TODO(), appInstanceId, serviceId)
	if err != nil {
		log.Debug("Failed to send request to delete service resource on mec platform ", resp.Status)
		return err
	}
	return nil
}

// Client request to subscribe service-availability notifications
func subscribeAvailability(appInstanceId string, callbackReference string) error {
	log.Debug("Sending request to mec platform add service-avail subscription api")
	var filter smc.SerAvailabilityNotificationSubscriptionFilteringCriteria
	filter.SerNames = nil
	filter.IsLocal = true
	subscription := smc.SerAvailabilityNotificationSubscription{"SerAvailabilityNotificationSubscription", callbackReference, nil, &filter}
	serAvailabilityNotificationSubscription, resp, err := srvMgmtClient.MecServiceMgmtApi.ApplicationsSubscriptionsPOST(context.TODO(), subscription, appInstanceId)
	if err != nil {
		log.Error("Failed to send service subscription: ", resp.Status)
		return err
	}

	hRefLink := serAvailabilityNotificationSubscription.Links.Self.Href

	// Find subscription id from response
	idPosition := strings.LastIndex(hRefLink, "/")
	if idPosition == -1 {
		log.Error("Error parsing subscription id for subscription")
	}
	serviceSubscriptionId = hRefLink[idPosition+1:]

	log.Info("Subscribed to service availibility notification on mec platform")

	return nil
}

// Client request to sent confirm terminate
func confirmTerminate(appInstanceId string) {

	operationAction := asc.TERMINATING_OperationActionType
	var terminationBody asc.AppTerminationConfirmation
	terminationBody.OperationAction = &operationAction
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsConfirmTerminationPOST(context.TODO(), terminationBody, appInstanceId)
	if err != nil {
		log.Error("Failed to send confirm termination ", resp.Status)
	} else {
		log.Info("Confirm Terminated")
	}

}

// Client request to subscribe app-termination notifications
func subscribeAppTermination(appInstanceId string, callBackReference string) error {
	log.Debug("Sending request to mec platform add app terminate subscription api")
	var appTerminationBody asc.AppTerminationNotificationSubscription
	appTerminationBody.SubscriptionType = "AppTerminationNotificationSubscription"
	appTerminationBody.CallbackReference = callBackReference
	appTerminationBody.AppInstanceId = appInstanceId
	appTerminationResponse, resp, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionsPOST(context.TODO(), appTerminationBody, appInstanceId)
	if err != nil {
		log.Error("Failed to send termination subscription: ", resp.Status)
		return err
	}

	hRefLink := appTerminationResponse.Links.Self.Href

	// Find subscription id from response
	idPosition := strings.LastIndex(hRefLink, "/")
	if idPosition == -1 {
		log.Error("Error parsing subscription id for subscription")
	}
	terminationSubscriptionId = hRefLink[idPosition+1:]

	return nil
}

// Client request to delete app-termination subscriptions
func delAppTerminationSubscription(appInstanceId string, subscriptionId string) error {
	resp, err := appSupportClient.MecAppSupportApi.ApplicationsSubscriptionDELETE(context.TODO(), appInstanceId, subscriptionId)
	if err != nil {
		log.Error("Failed to clear app termination subscription ", resp.Status)
		return err
	}
	return nil
}

// Client request to delete subscription of service-availability notifications
func delsubscribeAvailability(appInstanceId string, subscriptionId string) error {
	resp, err := srvMgmtClient.MecServiceMgmtApi.ApplicationsSubscriptionDELETE(context.TODO(), appInstanceId, subscriptionId)
	if err != nil {
		log.Error("Failed to clear service availability subscriptions: ", resp.Status)
		return err
	}
	return nil
}

// Client request to delete ams service
func delAmsService(serviceId string) error {
	resp, err := amsClient.AmsiApi.AppMobilityServiceByIdDELETE(context.TODO(), serviceId)
	if err != nil {
		log.Error("Failed to cleared ams service ", resp.Status)
		return err
	}

	return nil
}

// Client request to delete ams subscription
func deleteAmsSubscription(subscriptionId string) error {
	//log.Debug("Sending request to mec platform delete ams susbcription api")
	if amsSubscriptionSent {
		resp, err := amsClient.AmsiApi.SubByIdDELETE(context.TODO(), subscriptionId)
		if err != nil {
			log.Error("Failed to clear ams subcription ", resp.Status)
			return err
		}
	}
	return nil
}

// Confirm app readiness & app termination subscription
func Run() {

	// Confirm application readiness
	if !confirmReady {
		err := sendReadyConfirmation(instanceName)
		if err != nil {
			log.Fatal("Check configurations if valid")
		} else {
			log.Info("User app instance ", instanceName[0:6], ".... is ready to mec platform")
		}
	}

	// Subscribe for App Termination notifications
	if !terminationSubscription {
		callBackReference := local + localPort + "/application/termination"
		err := subscribeAppTermination(instanceName, callBackReference)
		if err == nil {
			log.Info("Subscribed to app termination notification on mec platform")
		}
	}

}

// Terminate by deleting all resources allocated on MEC platform & mec app
func Terminate() {

	// Only invoke graceful termination if not terminated
	if !terminated {
		//Delete app subscriptions
		err := delAppTerminationSubscription(instanceName, terminationSubscriptionId)
		if err == nil {
			log.Info("Cleared app-termination subscription on mec platform")
		}

		// Delete service subscriptions
		if subscriptionSent {
			err := delsubscribeAvailability(instanceName, serviceSubscriptionId)
			if err == nil {
				log.Info("Cleared service-avail subscription on mec platform")
			}
		}

		// Delete service
		if registeredService {
			err := unregisterService(instanceName, appEnablementServiceId)
			if err == nil {
				log.Info("Cleared user-app services on mec platform")
			}
		}

		// Delete ams service
		if amsServiceCreated {
			err := delAmsService(amsServiceId)
			if err == nil {
				log.Info("Cleared ams service on mec platform")
			}
		}

		// Delete ams subscriptions
		if amsSubscriptionSent {
			err := deleteAmsSubscription(amsSubscriptionId)
			if err == nil {
				log.Info("Cleared ams subcription on mec platform")
			}
		}

		//Send Confirm Terminate if received notification
		if terminateNotification {
			confirmTerminate(instanceName)
		}

		terminated = true
	}

}
