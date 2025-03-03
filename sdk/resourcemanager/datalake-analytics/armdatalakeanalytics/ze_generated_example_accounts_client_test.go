//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatalakeanalytics_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datalake-analytics/armdatalakeanalytics"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/datalake-analytics/resource-manager/Microsoft.DataLakeAnalytics/preview/2019-11-01-preview/examples/Accounts_List.json
func ExampleAccountsClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatalakeanalytics.NewAccountsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager(&armdatalakeanalytics.AccountsClientListOptions{Filter: to.Ptr("test_filter"),
		Top:     to.Ptr[int32](1),
		Skip:    to.Ptr[int32](1),
		Select:  to.Ptr("test_select"),
		Orderby: to.Ptr("test_orderby"),
		Count:   to.Ptr(false),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/datalake-analytics/resource-manager/Microsoft.DataLakeAnalytics/preview/2019-11-01-preview/examples/Accounts_ListByResourceGroup.json
func ExampleAccountsClient_NewListByResourceGroupPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatalakeanalytics.NewAccountsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListByResourceGroupPager("contosorg",
		&armdatalakeanalytics.AccountsClientListByResourceGroupOptions{Filter: to.Ptr("test_filter"),
			Top:     to.Ptr[int32](1),
			Skip:    to.Ptr[int32](1),
			Select:  to.Ptr("test_select"),
			Orderby: to.Ptr("test_orderby"),
			Count:   to.Ptr(false),
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/datalake-analytics/resource-manager/Microsoft.DataLakeAnalytics/preview/2019-11-01-preview/examples/Accounts_Create.json
func ExampleAccountsClient_BeginCreate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatalakeanalytics.NewAccountsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginCreate(ctx,
		"contosorg",
		"contosoadla",
		armdatalakeanalytics.CreateDataLakeAnalyticsAccountParameters{
			Location: to.Ptr("eastus2"),
			Properties: &armdatalakeanalytics.CreateDataLakeAnalyticsAccountProperties{
				ComputePolicies: []*armdatalakeanalytics.CreateComputePolicyWithAccountParameters{
					{
						Name: to.Ptr("test_policy"),
						Properties: &armdatalakeanalytics.CreateOrUpdateComputePolicyProperties{
							MaxDegreeOfParallelismPerJob: to.Ptr[int32](1),
							MinPriorityPerJob:            to.Ptr[int32](1),
							ObjectID:                     to.Ptr("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345"),
							ObjectType:                   to.Ptr(armdatalakeanalytics.AADObjectTypeUser),
						},
					}},
				DataLakeStoreAccounts: []*armdatalakeanalytics.AddDataLakeStoreWithAccountParameters{
					{
						Name: to.Ptr("test_adls"),
						Properties: &armdatalakeanalytics.AddDataLakeStoreProperties{
							Suffix: to.Ptr("test_suffix"),
						},
					}},
				DefaultDataLakeStoreAccount: to.Ptr("test_adls"),
				FirewallAllowAzureIPs:       to.Ptr(armdatalakeanalytics.FirewallAllowAzureIPsStateEnabled),
				FirewallRules: []*armdatalakeanalytics.CreateFirewallRuleWithAccountParameters{
					{
						Name: to.Ptr("test_rule"),
						Properties: &armdatalakeanalytics.CreateOrUpdateFirewallRuleProperties{
							EndIPAddress:   to.Ptr("2.2.2.2"),
							StartIPAddress: to.Ptr("1.1.1.1"),
						},
					}},
				FirewallState:                to.Ptr(armdatalakeanalytics.FirewallStateEnabled),
				MaxDegreeOfParallelism:       to.Ptr[int32](30),
				MaxDegreeOfParallelismPerJob: to.Ptr[int32](1),
				MaxJobCount:                  to.Ptr[int32](3),
				MinPriorityPerJob:            to.Ptr[int32](1),
				NewTier:                      to.Ptr(armdatalakeanalytics.TierTypeConsumption),
				QueryStoreRetention:          to.Ptr[int32](30),
				StorageAccounts: []*armdatalakeanalytics.AddStorageAccountWithAccountParameters{
					{
						Name: to.Ptr("test_storage"),
						Properties: &armdatalakeanalytics.AddStorageAccountProperties{
							AccessKey: to.Ptr("34adfa4f-cedf-4dc0-ba29-b6d1a69ab346"),
							Suffix:    to.Ptr("test_suffix"),
						},
					}},
			},
			Tags: map[string]*string{
				"test_key": to.Ptr("test_value"),
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/datalake-analytics/resource-manager/Microsoft.DataLakeAnalytics/preview/2019-11-01-preview/examples/Accounts_Get.json
func ExampleAccountsClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatalakeanalytics.NewAccountsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"contosorg",
		"contosoadla",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/datalake-analytics/resource-manager/Microsoft.DataLakeAnalytics/preview/2019-11-01-preview/examples/Accounts_Update.json
func ExampleAccountsClient_BeginUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatalakeanalytics.NewAccountsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginUpdate(ctx,
		"contosorg",
		"contosoadla",
		&armdatalakeanalytics.AccountsClientBeginUpdateOptions{Parameters: &armdatalakeanalytics.UpdateDataLakeAnalyticsAccountParameters{
			Properties: &armdatalakeanalytics.UpdateDataLakeAnalyticsAccountProperties{
				ComputePolicies: []*armdatalakeanalytics.UpdateComputePolicyWithAccountParameters{
					{
						Name: to.Ptr("test_policy"),
						Properties: &armdatalakeanalytics.UpdateComputePolicyProperties{
							MaxDegreeOfParallelismPerJob: to.Ptr[int32](1),
							MinPriorityPerJob:            to.Ptr[int32](1),
							ObjectID:                     to.Ptr("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345"),
							ObjectType:                   to.Ptr(armdatalakeanalytics.AADObjectTypeUser),
						},
					}},
				FirewallAllowAzureIPs: to.Ptr(armdatalakeanalytics.FirewallAllowAzureIPsStateEnabled),
				FirewallRules: []*armdatalakeanalytics.UpdateFirewallRuleWithAccountParameters{
					{
						Name: to.Ptr("test_rule"),
						Properties: &armdatalakeanalytics.UpdateFirewallRuleProperties{
							EndIPAddress:   to.Ptr("2.2.2.2"),
							StartIPAddress: to.Ptr("1.1.1.1"),
						},
					}},
				FirewallState:                to.Ptr(armdatalakeanalytics.FirewallStateEnabled),
				MaxDegreeOfParallelism:       to.Ptr[int32](1),
				MaxDegreeOfParallelismPerJob: to.Ptr[int32](1),
				MaxJobCount:                  to.Ptr[int32](1),
				MinPriorityPerJob:            to.Ptr[int32](1),
				NewTier:                      to.Ptr(armdatalakeanalytics.TierTypeConsumption),
				QueryStoreRetention:          to.Ptr[int32](1),
			},
			Tags: map[string]*string{
				"test_key": to.Ptr("test_value"),
			},
		},
		})
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/datalake-analytics/resource-manager/Microsoft.DataLakeAnalytics/preview/2019-11-01-preview/examples/Accounts_Delete.json
func ExampleAccountsClient_BeginDelete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatalakeanalytics.NewAccountsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	poller, err := client.BeginDelete(ctx,
		"contosorg",
		"contosoadla",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/datalake-analytics/resource-manager/Microsoft.DataLakeAnalytics/preview/2019-11-01-preview/examples/Accounts_CheckNameAvailability.json
func ExampleAccountsClient_CheckNameAvailability() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armdatalakeanalytics.NewAccountsClient("34adfa4f-cedf-4dc0-ba29-b6d1a69ab345", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CheckNameAvailability(ctx,
		"EastUS2",
		armdatalakeanalytics.CheckNameAvailabilityParameters{
			Name: to.Ptr("contosoadla"),
			Type: to.Ptr(armdatalakeanalytics.CheckNameAvailabilityParametersTypeMicrosoftDataLakeAnalyticsAccounts),
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}
