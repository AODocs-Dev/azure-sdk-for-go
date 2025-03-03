//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementListApiDiagnostics.json
func ExampleAPIDiagnosticClient_NewListByServicePager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIDiagnosticClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByServicePager("rg1",
		"apimService1",
		"echo-api",
		&armapimanagement.APIDiagnosticClientListByServiceOptions{Filter: nil,
			Top:  nil,
			Skip: nil,
		})
	for pager.More() {
		nextResult, err := pager.NextPage(ctx)
		if err != nil {
			log.Fatalf("failed to advance page: %v", err)
		}
		for _, v := range nextResult.Value {
			// TODO: use page item
			_ = v
		}
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementHeadApiDiagnostic.json
func ExampleAPIDiagnosticClient_GetEntityTag() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIDiagnosticClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.GetEntityTag(ctx,
		"rg1",
		"apimService1",
		"57d1f7558aa04f15146d9d8a",
		"applicationinsights",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementGetApiDiagnostic.json
func ExampleAPIDiagnosticClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIDiagnosticClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"rg1",
		"apimService1",
		"57d1f7558aa04f15146d9d8a",
		"applicationinsights",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementCreateApiDiagnostic.json
func ExampleAPIDiagnosticClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIDiagnosticClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateOrUpdate(ctx,
		"rg1",
		"apimService1",
		"57d1f7558aa04f15146d9d8a",
		"applicationinsights",
		armapimanagement.DiagnosticContract{
			Properties: &armapimanagement.DiagnosticContractProperties{
				AlwaysLog: to.Ptr(armapimanagement.AlwaysLogAllErrors),
				Backend: &armapimanagement.PipelineDiagnosticSettings{
					Response: &armapimanagement.HTTPMessageDiagnostic{
						Body: &armapimanagement.BodyDiagnosticSettings{
							Bytes: to.Ptr[int32](512),
						},
						Headers: []*string{
							to.Ptr("Content-type")},
					},
					Request: &armapimanagement.HTTPMessageDiagnostic{
						Body: &armapimanagement.BodyDiagnosticSettings{
							Bytes: to.Ptr[int32](512),
						},
						Headers: []*string{
							to.Ptr("Content-type")},
					},
				},
				Frontend: &armapimanagement.PipelineDiagnosticSettings{
					Response: &armapimanagement.HTTPMessageDiagnostic{
						Body: &armapimanagement.BodyDiagnosticSettings{
							Bytes: to.Ptr[int32](512),
						},
						Headers: []*string{
							to.Ptr("Content-type")},
					},
					Request: &armapimanagement.HTTPMessageDiagnostic{
						Body: &armapimanagement.BodyDiagnosticSettings{
							Bytes: to.Ptr[int32](512),
						},
						Headers: []*string{
							to.Ptr("Content-type")},
					},
				},
				LoggerID: to.Ptr("/loggers/applicationinsights"),
				Sampling: &armapimanagement.SamplingSettings{
					Percentage:   to.Ptr[float64](50),
					SamplingType: to.Ptr(armapimanagement.SamplingTypeFixed),
				},
			},
		},
		&armapimanagement.APIDiagnosticClientCreateOrUpdateOptions{IfMatch: nil})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementUpdateApiDiagnostic.json
func ExampleAPIDiagnosticClient_Update() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIDiagnosticClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Update(ctx,
		"rg1",
		"apimService1",
		"echo-api",
		"applicationinsights",
		"*",
		armapimanagement.DiagnosticContract{
			Properties: &armapimanagement.DiagnosticContractProperties{
				AlwaysLog: to.Ptr(armapimanagement.AlwaysLogAllErrors),
				Backend: &armapimanagement.PipelineDiagnosticSettings{
					Response: &armapimanagement.HTTPMessageDiagnostic{
						Body: &armapimanagement.BodyDiagnosticSettings{
							Bytes: to.Ptr[int32](512),
						},
						Headers: []*string{
							to.Ptr("Content-type")},
					},
					Request: &armapimanagement.HTTPMessageDiagnostic{
						Body: &armapimanagement.BodyDiagnosticSettings{
							Bytes: to.Ptr[int32](512),
						},
						Headers: []*string{
							to.Ptr("Content-type")},
					},
				},
				Frontend: &armapimanagement.PipelineDiagnosticSettings{
					Response: &armapimanagement.HTTPMessageDiagnostic{
						Body: &armapimanagement.BodyDiagnosticSettings{
							Bytes: to.Ptr[int32](512),
						},
						Headers: []*string{
							to.Ptr("Content-type")},
					},
					Request: &armapimanagement.HTTPMessageDiagnostic{
						Body: &armapimanagement.BodyDiagnosticSettings{
							Bytes: to.Ptr[int32](512),
						},
						Headers: []*string{
							to.Ptr("Content-type")},
					},
				},
				LoggerID: to.Ptr("/loggers/applicationinsights"),
				Sampling: &armapimanagement.SamplingSettings{
					Percentage:   to.Ptr[float64](50),
					SamplingType: to.Ptr(armapimanagement.SamplingTypeFixed),
				},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/apimanagement/resource-manager/Microsoft.ApiManagement/stable/2021-08-01/examples/ApiManagementDeleteApiDiagnostic.json
func ExampleAPIDiagnosticClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armapimanagement.NewAPIDiagnosticClient("subid", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"rg1",
		"apimService1",
		"57d1f7558aa04f15146d9d8a",
		"applicationinsights",
		"*",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
