/*
 * Copyright (c) 2019  InterDigital Communications, Inc
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
 * AdvantEDGE Platform Controller REST API
 * This API is the main platform API and mainly used by the AdvantEDGE frontend to interact with scenarios <p>**Micro-service**<br>[meep-ctrl-engine](https://github.com/InterDigitalInc/AdvantEDGE/tree/master/go-apps/meep-ctrl-engine) <p>**Type & Usage**<br>Platform main interface used by controller software that want to interact with the AdvantEDGE platform <p>**Details**<br>API details available at _your-AdvantEDGE-ip-address:30000/api_ <p>**Default Port**<br>`30000` 
 *
 * OpenAPI spec version: 1.0.0
 * Contact: AdvantEDGE@InterDigital.com
 *
 * NOTE: This class is auto generated by the swagger code generator program.
 * https://github.com/swagger-api/swagger-codegen.git
 *
 * Swagger Codegen version: 2.4.9
 *
 * Do not edit the class manually.
 *
 */

(function(root, factory) {
  if (typeof define === 'function' && define.amd) {
    // AMD. Register as an anonymous module.
    define(['ApiClient', 'model/Replay', 'model/ReplayFileList', 'model/ScenarioInfo'], factory);
  } else if (typeof module === 'object' && module.exports) {
    // CommonJS-like environments that support module.exports, like Node.
    module.exports = factory(require('../ApiClient'), require('../model/Replay'), require('../model/ReplayFileList'), require('../model/ScenarioInfo'));
  } else {
    // Browser globals (root is window)
    if (!root.AdvantEdgePlatformControllerRestApi) {
      root.AdvantEdgePlatformControllerRestApi = {};
    }
    root.AdvantEdgePlatformControllerRestApi.EventReplayApi = factory(root.AdvantEdgePlatformControllerRestApi.ApiClient, root.AdvantEdgePlatformControllerRestApi.Replay, root.AdvantEdgePlatformControllerRestApi.ReplayFileList, root.AdvantEdgePlatformControllerRestApi.ScenarioInfo);
  }
}(this, function(ApiClient, Replay, ReplayFileList, ScenarioInfo) {
  'use strict';

  /**
   * EventReplay service.
   * @module api/EventReplayApi
   * @version 1.0.0
   */

  /**
   * Constructs a new EventReplayApi. 
   * @alias module:api/EventReplayApi
   * @class
   * @param {module:ApiClient} [apiClient] Optional API client implementation to use,
   * default to {@link module:ApiClient#instance} if unspecified.
   */
  var exports = function(apiClient) {
    this.apiClient = apiClient || ApiClient.instance;


    /**
     * Callback function to receive the result of the createReplayFile operation.
     * @callback module:api/EventReplayApi~createReplayFileCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Add a replay file
     * Add a replay file to the platform store
     * @param {String} name Replay file name
     * @param {module:model/Replay} replayFile Replay-file
     * @param {module:api/EventReplayApi~createReplayFileCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.createReplayFile = function(name, replayFile, callback) {
      var postBody = replayFile;

      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling createReplayFile");
      }

      // verify the required parameter 'replayFile' is set
      if (replayFile === undefined || replayFile === null) {
        throw new Error("Missing the required parameter 'replayFile' when calling createReplayFile");
      }


      var pathParams = {
        'name': name
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/replay/{name}', 'POST',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the createReplayFileFromScenarioExec operation.
     * @callback module:api/EventReplayApi~createReplayFileFromScenarioExecCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Generate a replay file from scenario execution events
     * Generate a replay file using events from the latest execution of a scenario
     * @param {String} name Replay file name
     * @param {module:model/ScenarioInfo} scenarioInfo Scenario information
     * @param {module:api/EventReplayApi~createReplayFileFromScenarioExecCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.createReplayFileFromScenarioExec = function(name, scenarioInfo, callback) {
      var postBody = scenarioInfo;

      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling createReplayFileFromScenarioExec");
      }

      // verify the required parameter 'scenarioInfo' is set
      if (scenarioInfo === undefined || scenarioInfo === null) {
        throw new Error("Missing the required parameter 'scenarioInfo' when calling createReplayFileFromScenarioExec");
      }


      var pathParams = {
        'name': name
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/replay/{name}/generate', 'POST',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the deleteReplayFile operation.
     * @callback module:api/EventReplayApi~deleteReplayFileCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Delete a replay file
     * Delete a replay file by name from the platform store
     * @param {String} name replay file name
     * @param {module:api/EventReplayApi~deleteReplayFileCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.deleteReplayFile = function(name, callback) {
      var postBody = null;

      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling deleteReplayFile");
      }


      var pathParams = {
        'name': name
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/replay/{name}', 'DELETE',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the deleteReplayFileList operation.
     * @callback module:api/EventReplayApi~deleteReplayFileListCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Delete all replay files
     * Delete all replay files present in the platform store
     * @param {module:api/EventReplayApi~deleteReplayFileListCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.deleteReplayFileList = function(callback) {
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/replay', 'DELETE',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getReplayFile operation.
     * @callback module:api/EventReplayApi~getReplayFileCallback
     * @param {String} error Error message, if any.
     * @param {module:model/Replay} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get a specific replay file
     * Get a replay file by name from the platform store
     * @param {String} name Replay file name
     * @param {module:api/EventReplayApi~getReplayFileCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/Replay}
     */
    this.getReplayFile = function(name, callback) {
      var postBody = null;

      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling getReplayFile");
      }


      var pathParams = {
        'name': name
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = Replay;

      return this.apiClient.callApi(
        '/replay/{name}', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the getReplayFileList operation.
     * @callback module:api/EventReplayApi~getReplayFileListCallback
     * @param {String} error Error message, if any.
     * @param {module:model/ReplayFileList} data The data returned by the service call.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Get all replay file names
     * Returns a list of all replay files names present in the platform store
     * @param {module:api/EventReplayApi~getReplayFileListCallback} callback The callback function, accepting three arguments: error, data, response
     * data is of type: {@link module:model/ReplayFileList}
     */
    this.getReplayFileList = function(callback) {
      var postBody = null;


      var pathParams = {
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = ReplayFileList;

      return this.apiClient.callApi(
        '/replay', 'GET',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the loopReplay operation.
     * @callback module:api/EventReplayApi~loopReplayCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Loop-Execute a replay file present in the platform store
     * Loop-Execute a replay file present in the platform store
     * @param {String} name Replay file name
     * @param {module:api/EventReplayApi~loopReplayCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.loopReplay = function(name, callback) {
      var postBody = null;

      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling loopReplay");
      }


      var pathParams = {
        'name': name
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/replay/{name}/loop', 'POST',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the playReplayFile operation.
     * @callback module:api/EventReplayApi~playReplayFileCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Execute a replay file present in the platform store
     * Execute a replay file present in the platform store
     * @param {String} name Replay file name
     * @param {module:api/EventReplayApi~playReplayFileCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.playReplayFile = function(name, callback) {
      var postBody = null;

      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling playReplayFile");
      }


      var pathParams = {
        'name': name
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/replay/{name}/play', 'POST',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }

    /**
     * Callback function to receive the result of the stopReplayFile operation.
     * @callback module:api/EventReplayApi~stopReplayFileCallback
     * @param {String} error Error message, if any.
     * @param data This operation does not return a value.
     * @param {String} response The complete HTTP response.
     */

    /**
     * Stop execution of a replay file
     * Stop execution a replay file
     * @param {String} name Replay file name
     * @param {module:api/EventReplayApi~stopReplayFileCallback} callback The callback function, accepting three arguments: error, data, response
     */
    this.stopReplayFile = function(name, callback) {
      var postBody = null;

      // verify the required parameter 'name' is set
      if (name === undefined || name === null) {
        throw new Error("Missing the required parameter 'name' when calling stopReplayFile");
      }


      var pathParams = {
        'name': name
      };
      var queryParams = {
      };
      var collectionQueryParams = {
      };
      var headerParams = {
      };
      var formParams = {
      };

      var authNames = [];
      var contentTypes = ['application/json'];
      var accepts = ['application/json'];
      var returnType = null;

      return this.apiClient.callApi(
        '/replay/{name}/stop', 'POST',
        pathParams, queryParams, collectionQueryParams, headerParams, formParams, postBody,
        authNames, contentTypes, accepts, returnType, callback
      );
    }
  };

  return exports;
}));
