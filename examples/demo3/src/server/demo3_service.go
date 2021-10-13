/*
 * MEC Use-Case 3 API
 *
 * This section describes a use case that the user can accomplish using the MEC Sandbox APIs from a MEC application
 *
 * API version: 0.0.1
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package server

import (
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

var srvMgmtClient *smc.APIClient
var srvMgmtClientPath string
var appSupportClient *asc.APIClient
var appSupportClientPath string

var amsClient *ams.APIClient
var amsClientPath string

var instanceName string
var mecUrl string
var localPort string
var subscriptionSent bool
var confirmReady bool
var appEnablementServiceId string
var serviceSubscriptionId string
var registeredService bool
var terminationSubscriptionId string
var mep string
var amsSubscriptionId string

var serviceName string = "user-app"
var scopeOfLocality string = defaultScopeOfLocality
var consumedLocalOnly bool = defaultConsumedLocalOnly
var terminationSubscription bool = false
var terminated bool = false
var terminateNotification bool = false
var amsSubscriptionSent bool = false

const serviceAppVersion = "2.1.1"
const local = "http://10.190.115.162"
const defaultScopeOfLocality = "MEC_SYSTEM"
const defaultConsumedLocalOnly = true

func Init(envPath string, envName string) (port string, err error) {
	var config util.Config
	var configErr error

	log.Info("Using config from ", envPath, "/", envName)
	config, configErr = util.LoadConfig(envPath, envName)

	if configErr != nil {
		log.Fatal(configErr)
	}

	// Retrieve app id from config
	instanceName = config.AppInstanceId

	// Retrieve sandbox url from config
	mecUrl = config.SandboxUrl

	// Parse mec platfor mep no. use for ams service
	resp := strings.LastIndex(mecUrl, "/")
	if resp == -1 {
		log.Fatal("Error parsing mep no. from config")
	} else {
		mep = mecUrl[resp+1:]
	}

	// Retreieve local url from config
	localPort = config.Port

	// Retrieve service name config if present
	if config.ServiceName != "" {
		serviceName = config.ServiceName
	}

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

	// Create application mobility suppport client
	amsClientcfg := ams.NewConfiguration()

	// Replace amsUrl with mep1 to demonstrate use-case of ams api
	amsUrl := strings.Replace(mecUrl, "mep2", "mep1", 1)
	amsClientcfg.BasePath = amsUrl + "/amsi/v1"
	amsClient = ams.NewAPIClient(amsClientcfg)
	amsClientPath = amsClientcfg.BasePath
	if amsClient == nil {
		return "", errors.New("Failed to create Application Mobility Support REST API client")
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

	w.WriteHeader(http.StatusNoContent)
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
		fmt.Fprintf(w, "Sucessfully created a service")
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
	}
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
// Handle AMS notification
func amsNotificationPOST(w http.ResponseWriter, r *http.Request) {
	var amsNotification ams.MobilityProcedureNotification
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&amsNotification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Info("AMS event received for ", amsNotification.AssociateId[0].Value, " to ", amsNotification.TargetAppInfo.AppInstanceId)
	w.WriteHeader(http.StatusOK)
}

// Rest API
// Submit AMS subscription to mec platform
func amsSubscriptionPOST(w http.ResponseWriter, r *http.Request) {
	amsSendSubscription(instanceName)
	w.WriteHeader(http.StatusOK)
}

// Rest API
// Register MEC Application instances with AMS & consume servicee
func amsCreatePOST(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	amsSendService(instanceName)
}

// CLient request to create a new application mobility service
func amsSendService(appInstanceId string) {
	log.Debug("Sending request to mec platform create ams api")
	var bodyRegisterationInfo ams.RegistrationInfo
	bodyRegisterationInfo.ServiceConsumerId = &ams.RegistrationInfoServiceConsumerId{
		AppInstanceId: appInstanceId,
	}

	// Provide device info only specific to mep1 platform
	if mep == "mep1" {
		var associateId ams.AssociateId
		associateId.Type_ = 1
		associateId.Value = "10.100.0.10"
		bodyRegisterationInfo.DeviceInformation = append(bodyRegisterationInfo.DeviceInformation, ams.RegistrationInfoDeviceInformation{AssociateId: &associateId,
			AppMobilityServiceLevel: 3,
			ContextTransferState:    0,
		})
	}

	_, _, err := amsClient.AmsiApi.AppMobilityServicePOST(context.TODO(), bodyRegisterationInfo)

	if err != nil {
		log.Error(err)
	} else {
		log.Info("Consumed AMS service sucessfully")
	}
}

// CLient request to create a new application mobility service
func amsSendSubscription(appInstanceId string) {
	log.Debug("Sending request to mec platform add ams subscription api")

	var mobilityProcedureSubscription ams.MobilityProcedureSubscription

	// Add body param callback ref
	mobilityProcedureSubscription.CallbackReference = local + localPort + "/subscriptions"
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

	mobilitySubscription, resp, err := amsClient.AmsiApi.SubPOST(context.TODO(), mobilityProcedureSubscription)
	hRefLink := mobilitySubscription.Links.Self.Href

	// Find subscription id from response
	idPosition := strings.LastIndex(hRefLink, "/")
	if idPosition == -1 {
		log.Error("Error parsing subscription id for subscription")
	}
	amsSubscriptionId = hRefLink[idPosition+1:]

	if err != nil {
		log.Error(resp.Status)
	} else {
		amsSubscriptionSent = true
		log.Info("Successfully created ams subscription")
	}

}

// Client request to delete ams subscription
func deleteAmsSubscription(subscriptionId string) error {
	log.Debug("Sending request to mec platform delete ams susbcription api")
	if amsSubscriptionSent {
		resp, err := amsClient.AmsiApi.AppMobilityServiceByIdDELETE(context.TODO(), subscriptionId)
		if err != nil {
			log.Info("Failed to delete ams subcription ", resp.Status)
			return err
		}
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
	appServicesPostResponse, resp, err := srvMgmtClient.MecServiceMgmtApi.ServicesGET(context.TODO(), nil)
	log.Debug("Sending request to mec platform get services api ")
	if err != nil {
		log.Error("Failed to fetch services on mec platform ", resp.Status)
		return err
	}

	log.Info("Returning available mec services on mec platform")
	servicesName := make([]string, len(appServicesPostResponse))
	for i := 0; i < len(appServicesPostResponse); i++ {
		servicesName[i] = appServicesPostResponse[i].SerName + " URL: " + appServicesPostResponse[i].TransportInfo.Endpoint.Uris[0]
	}

	for _, v := range servicesName {
		log.Info(v)
	}

	return nil
}

// Client request to create a mec-service resource
func registerService(appInstanceId string) error {
	log.Debug("Sending request to mec platform post service api ")
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
	endpointPath := local + "/" + localPort
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
		log.Error("Failed to register service on mec app enablement registry: ", resp.Status)
		return err
	}
	log.Info(serviceName, " service created with instance id: ", appServicesPostResponse.SerInstanceId)
	appEnablementServiceId = appServicesPostResponse.SerInstanceId
	registeredService = true
	return nil
}

// Client request to delete a mec-service resource
func unregisterService(appInstanceId string, serviceId string) error {
	log.Debug("Sending request to mec platform delete service api")
	resp, err := srvMgmtClient.MecServiceMgmtApi.AppServicesServiceIdDELETE(context.TODO(), appInstanceId, serviceId)
	if err != nil {
		log.Debug("Failed to send request to delete service on mec platform ", resp.Status)
		return err
	}
	return nil
}

// Client request to subscribe service-availability notifications
func subscribeAvailability(appInstanceId string, callbackReference string) error {
	log.Debug("Sending request to mec platform add subscription api")
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
	log.Debug("Sending request to mec platform confirm terminate subscription api")
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

// Confirm app readiness & app termination subscription
func Run() {

	// Confirm application readiness
	if !confirmReady {
		err := sendReadyConfirmation(instanceName)
		if err != nil {
			log.Fatal("Check configurations if valid")
		} else {
			log.Info("User app is ready to mec platform")
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

	// Only invoke graceful termination if terminated is false
	if !terminated {
		//Delete app subscriptions
		err := delAppTerminationSubscription(instanceName, terminationSubscriptionId)
		if err == nil {
			log.Info("Cleared app enablement subscription on mec platform")
		}

		//Send Confirm Terminate if received notification
		if terminateNotification {
			confirmTerminate(instanceName)
		}

		// Delete service subscriptions
		if subscriptionSent {
			err := delsubscribeAvailability(instanceName, serviceSubscriptionId)
			if err == nil {
				log.Info("Cleared service subscription on mec platform")
			}
		}

		// Delete service
		if registeredService {
			err := unregisterService(instanceName, appEnablementServiceId)
			if err == nil {
				log.Info("Cleared user-app services on mec platform")
			}
		}

		// Delete ams subscriptions
		if amsSubscriptionSent {
			err := deleteAmsSubscription(amsSubscriptionId)
			if err == nil {
				log.Info("Cleared ams subcription on mec platform")
			}
		}

		terminated = true
	}

}
