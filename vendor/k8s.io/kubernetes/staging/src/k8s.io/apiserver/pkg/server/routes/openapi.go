/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package routes

import (
	"github.com/emicklei/go-restful"
	"github.com/go-openapi/spec"
	"k8s.io/klog"
	"k8s.io/kube-openapi/pkg/builder"

	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/kube-openapi/pkg/common"
	"k8s.io/kube-openapi/pkg/handler"
)

// OpenAPI installs spec endpoints for each web service.
type OpenAPI struct {
	Config *common.Config
}

// Install adds the SwaggerUI webservice to the given mux.
func (oa OpenAPI) Install(c *restful.Container, mux *mux.PathRecorderMux) (*handler.OpenAPIService, *spec.Swagger) {
	spec, err := builder.BuildOpenAPISpec(c.RegisteredWebServices(), oa.Config)
	if err != nil {
		klog.Fatalf("Failed to build open api spec for root: %v", err)
	}
	// NOTE: [DEPRECATION] We will announce deprecation for format-separated endpoints for OpenAPI spec,
	// and switch to a single /openapi/v2 endpoint in Kubernetes 1.10. The design doc and deprecation process
	// are tracked at: https://docs.google.com/document/d/19lEqE9lc4yHJ3WJAJxS_G7TcORIJXGHyq3wpwcH28nU.
	if _, err := handler.BuildAndRegisterOpenAPIService("/swagger.json", c.RegisteredWebServices(), oa.Config, mux); err != nil {
		klog.Fatalf("Failed to register open api spec for root: %v", err)
	}
	openAPIVersionedService, err := handler.RegisterOpenAPIVersionedService(spec, "/openapi/v2", mux)
	if err != nil {
		klog.Fatalf("Failed to register versioned open api spec for root: %v", err)
	}
	return openAPIVersionedService, spec
}
