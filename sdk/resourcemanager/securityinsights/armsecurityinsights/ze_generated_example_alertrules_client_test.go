//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurityinsights_test

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
)

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/alertRules/GetAllAlertRules.json
func ExampleAlertRulesClient_NewListPager() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewAlertRulesClient("d0cfe6b2-9ac0-4464-9919-dccaee2e48c0", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	pager := client.NewListPager("myRg",
		"myWorkspace",
		nil)
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

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/alertRules/GetFusionAlertRule.json
func ExampleAlertRulesClient_Get() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewAlertRulesClient("d0cfe6b2-9ac0-4464-9919-dccaee2e48c0", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.Get(ctx,
		"myRg",
		"myWorkspace",
		"myFirstFusionRule",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/alertRules/CreateFusionAlertRuleWithFusionScenarioExclusion.json
func ExampleAlertRulesClient_CreateOrUpdate() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewAlertRulesClient("d0cfe6b2-9ac0-4464-9919-dccaee2e48c0", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	res, err := client.CreateOrUpdate(ctx,
		"myRg",
		"myWorkspace",
		"myFirstFusionRule",
		&armsecurityinsights.FusionAlertRule{
			Etag: to.Ptr("3d00c3ca-0000-0100-0000-5d42d5010000"),
			Kind: to.Ptr(armsecurityinsights.AlertRuleKindFusion),
			Properties: &armsecurityinsights.FusionAlertRuleProperties{
				AlertRuleTemplateName: to.Ptr("f71aba3d-28fb-450b-b192-4e76a83015c8"),
				Enabled:               to.Ptr(true),
				SourceSettings: []*armsecurityinsights.FusionSourceSettings{
					{
						Enabled:    to.Ptr(true),
						SourceName: to.Ptr("Anomalies"),
					},
					{
						Enabled:    to.Ptr(true),
						SourceName: to.Ptr("Alert providers"),
						SourceSubTypes: []*armsecurityinsights.FusionSourceSubTypeSetting{
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Azure Active Directory Identity Protection"),
							},
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Azure Defender"),
							},
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Azure Defender for IoT"),
							},
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Microsoft 365 Defender"),
							},
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Microsoft Cloud App Security"),
							},
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Microsoft Defender for Endpoint"),
							},
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Microsoft Defender for Identity"),
							},
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Microsoft Defender for Office 365"),
							},
							{
								Enabled: to.Ptr(true),
								SeverityFilters: &armsecurityinsights.FusionSubTypeSeverityFilter{
									Filters: []*armsecurityinsights.FusionSubTypeSeverityFiltersItem{
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityHigh),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityMedium),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityLow),
										},
										{
											Enabled:  to.Ptr(true),
											Severity: to.Ptr(armsecurityinsights.AlertSeverityInformational),
										}},
								},
								SourceSubTypeName: to.Ptr("Azure Sentinel scheduled analytics rules"),
							}},
					},
					{
						Enabled:    to.Ptr(true),
						SourceName: to.Ptr("Raw logs from other sources"),
						SourceSubTypes: []*armsecurityinsights.FusionSourceSubTypeSetting{
							{
								Enabled:           to.Ptr(true),
								SeverityFilters:   &armsecurityinsights.FusionSubTypeSeverityFilter{},
								SourceSubTypeName: to.Ptr("Palo Alto Networks"),
							}},
					}},
			},
		},
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// TODO: use response item
	_ = res
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/tree/main/specification/securityinsights/resource-manager/Microsoft.SecurityInsights/preview/2022-05-01-preview/examples/alertRules/DeleteAlertRule.json
func ExampleAlertRulesClient_Delete() {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("failed to obtain a credential: %v", err)
	}
	ctx := context.Background()
	client, err := armsecurityinsights.NewAlertRulesClient("d0cfe6b2-9ac0-4464-9919-dccaee2e48c0", cred, nil)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}
	_, err = client.Delete(ctx,
		"myRg",
		"myWorkspace",
		"73e01a99-5cd7-4139-a149-9f2736ff2ab5",
		nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
}
