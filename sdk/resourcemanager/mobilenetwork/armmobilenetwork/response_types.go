//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armmobilenetwork

// AttachedDataNetworksClientCreateOrUpdateResponse contains the response from method AttachedDataNetworksClient.CreateOrUpdate.
type AttachedDataNetworksClientCreateOrUpdateResponse struct {
	AttachedDataNetwork
}

// AttachedDataNetworksClientDeleteResponse contains the response from method AttachedDataNetworksClient.Delete.
type AttachedDataNetworksClientDeleteResponse struct {
	// placeholder for future response values
}

// AttachedDataNetworksClientGetResponse contains the response from method AttachedDataNetworksClient.Get.
type AttachedDataNetworksClientGetResponse struct {
	AttachedDataNetwork
}

// AttachedDataNetworksClientListByPacketCoreDataPlaneResponse contains the response from method AttachedDataNetworksClient.ListByPacketCoreDataPlane.
type AttachedDataNetworksClientListByPacketCoreDataPlaneResponse struct {
	AttachedDataNetworkListResult
}

// AttachedDataNetworksClientUpdateTagsResponse contains the response from method AttachedDataNetworksClient.UpdateTags.
type AttachedDataNetworksClientUpdateTagsResponse struct {
	AttachedDataNetwork
}

// DataNetworksClientCreateOrUpdateResponse contains the response from method DataNetworksClient.CreateOrUpdate.
type DataNetworksClientCreateOrUpdateResponse struct {
	DataNetwork
}

// DataNetworksClientDeleteResponse contains the response from method DataNetworksClient.Delete.
type DataNetworksClientDeleteResponse struct {
	// placeholder for future response values
}

// DataNetworksClientGetResponse contains the response from method DataNetworksClient.Get.
type DataNetworksClientGetResponse struct {
	DataNetwork
}

// DataNetworksClientListByMobileNetworkResponse contains the response from method DataNetworksClient.ListByMobileNetwork.
type DataNetworksClientListByMobileNetworkResponse struct {
	DataNetworkListResult
}

// DataNetworksClientUpdateTagsResponse contains the response from method DataNetworksClient.UpdateTags.
type DataNetworksClientUpdateTagsResponse struct {
	DataNetwork
}

// MobileNetworksClientCreateOrUpdateResponse contains the response from method MobileNetworksClient.CreateOrUpdate.
type MobileNetworksClientCreateOrUpdateResponse struct {
	MobileNetwork
}

// MobileNetworksClientDeleteResponse contains the response from method MobileNetworksClient.Delete.
type MobileNetworksClientDeleteResponse struct {
	// placeholder for future response values
}

// MobileNetworksClientGetResponse contains the response from method MobileNetworksClient.Get.
type MobileNetworksClientGetResponse struct {
	MobileNetwork
}

// MobileNetworksClientListByResourceGroupResponse contains the response from method MobileNetworksClient.ListByResourceGroup.
type MobileNetworksClientListByResourceGroupResponse struct {
	ListResult
}

// MobileNetworksClientListBySubscriptionResponse contains the response from method MobileNetworksClient.ListBySubscription.
type MobileNetworksClientListBySubscriptionResponse struct {
	ListResult
}

// MobileNetworksClientListSimIDsResponse contains the response from method MobileNetworksClient.ListSimIDs.
type MobileNetworksClientListSimIDsResponse struct {
	SimIDListResult
}

// MobileNetworksClientUpdateTagsResponse contains the response from method MobileNetworksClient.UpdateTags.
type MobileNetworksClientUpdateTagsResponse struct {
	MobileNetwork
}

// OperationsClientListResponse contains the response from method OperationsClient.List.
type OperationsClientListResponse struct {
	OperationList
}

// PacketCoreControlPlaneVersionsClientGetResponse contains the response from method PacketCoreControlPlaneVersionsClient.Get.
type PacketCoreControlPlaneVersionsClientGetResponse struct {
	PacketCoreControlPlaneVersion
}

// PacketCoreControlPlaneVersionsClientListByResourceGroupResponse contains the response from method PacketCoreControlPlaneVersionsClient.ListByResourceGroup.
type PacketCoreControlPlaneVersionsClientListByResourceGroupResponse struct {
	PacketCoreControlPlaneVersionListResult
}

// PacketCoreControlPlanesClientCreateOrUpdateResponse contains the response from method PacketCoreControlPlanesClient.CreateOrUpdate.
type PacketCoreControlPlanesClientCreateOrUpdateResponse struct {
	PacketCoreControlPlane
}

// PacketCoreControlPlanesClientDeleteResponse contains the response from method PacketCoreControlPlanesClient.Delete.
type PacketCoreControlPlanesClientDeleteResponse struct {
	// placeholder for future response values
}

// PacketCoreControlPlanesClientGetResponse contains the response from method PacketCoreControlPlanesClient.Get.
type PacketCoreControlPlanesClientGetResponse struct {
	PacketCoreControlPlane
}

// PacketCoreControlPlanesClientListByResourceGroupResponse contains the response from method PacketCoreControlPlanesClient.ListByResourceGroup.
type PacketCoreControlPlanesClientListByResourceGroupResponse struct {
	PacketCoreControlPlaneListResult
}

// PacketCoreControlPlanesClientListBySubscriptionResponse contains the response from method PacketCoreControlPlanesClient.ListBySubscription.
type PacketCoreControlPlanesClientListBySubscriptionResponse struct {
	PacketCoreControlPlaneListResult
}

// PacketCoreControlPlanesClientUpdateTagsResponse contains the response from method PacketCoreControlPlanesClient.UpdateTags.
type PacketCoreControlPlanesClientUpdateTagsResponse struct {
	PacketCoreControlPlane
}

// PacketCoreDataPlanesClientCreateOrUpdateResponse contains the response from method PacketCoreDataPlanesClient.CreateOrUpdate.
type PacketCoreDataPlanesClientCreateOrUpdateResponse struct {
	PacketCoreDataPlane
}

// PacketCoreDataPlanesClientDeleteResponse contains the response from method PacketCoreDataPlanesClient.Delete.
type PacketCoreDataPlanesClientDeleteResponse struct {
	// placeholder for future response values
}

// PacketCoreDataPlanesClientGetResponse contains the response from method PacketCoreDataPlanesClient.Get.
type PacketCoreDataPlanesClientGetResponse struct {
	PacketCoreDataPlane
}

// PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse contains the response from method PacketCoreDataPlanesClient.ListByPacketCoreControlPlane.
type PacketCoreDataPlanesClientListByPacketCoreControlPlaneResponse struct {
	PacketCoreDataPlaneListResult
}

// PacketCoreDataPlanesClientUpdateTagsResponse contains the response from method PacketCoreDataPlanesClient.UpdateTags.
type PacketCoreDataPlanesClientUpdateTagsResponse struct {
	PacketCoreDataPlane
}

// ServicesClientCreateOrUpdateResponse contains the response from method ServicesClient.CreateOrUpdate.
type ServicesClientCreateOrUpdateResponse struct {
	Service
}

// ServicesClientDeleteResponse contains the response from method ServicesClient.Delete.
type ServicesClientDeleteResponse struct {
	// placeholder for future response values
}

// ServicesClientGetResponse contains the response from method ServicesClient.Get.
type ServicesClientGetResponse struct {
	Service
}

// ServicesClientListByMobileNetworkResponse contains the response from method ServicesClient.ListByMobileNetwork.
type ServicesClientListByMobileNetworkResponse struct {
	ServiceListResult
}

// ServicesClientUpdateTagsResponse contains the response from method ServicesClient.UpdateTags.
type ServicesClientUpdateTagsResponse struct {
	Service
}

// SimGroupsClientCreateOrUpdateResponse contains the response from method SimGroupsClient.CreateOrUpdate.
type SimGroupsClientCreateOrUpdateResponse struct {
	SimGroup
}

// SimGroupsClientDeleteResponse contains the response from method SimGroupsClient.Delete.
type SimGroupsClientDeleteResponse struct {
	// placeholder for future response values
}

// SimGroupsClientGetResponse contains the response from method SimGroupsClient.Get.
type SimGroupsClientGetResponse struct {
	SimGroup
}

// SimGroupsClientListByResourceGroupResponse contains the response from method SimGroupsClient.ListByResourceGroup.
type SimGroupsClientListByResourceGroupResponse struct {
	SimGroupListResult
}

// SimGroupsClientListBySubscriptionResponse contains the response from method SimGroupsClient.ListBySubscription.
type SimGroupsClientListBySubscriptionResponse struct {
	SimGroupListResult
}

// SimGroupsClientUpdateTagsResponse contains the response from method SimGroupsClient.UpdateTags.
type SimGroupsClientUpdateTagsResponse struct {
	SimGroup
}

// SimPoliciesClientCreateOrUpdateResponse contains the response from method SimPoliciesClient.CreateOrUpdate.
type SimPoliciesClientCreateOrUpdateResponse struct {
	SimPolicy
}

// SimPoliciesClientDeleteResponse contains the response from method SimPoliciesClient.Delete.
type SimPoliciesClientDeleteResponse struct {
	// placeholder for future response values
}

// SimPoliciesClientGetResponse contains the response from method SimPoliciesClient.Get.
type SimPoliciesClientGetResponse struct {
	SimPolicy
}

// SimPoliciesClientListByMobileNetworkResponse contains the response from method SimPoliciesClient.ListByMobileNetwork.
type SimPoliciesClientListByMobileNetworkResponse struct {
	SimPolicyListResult
}

// SimPoliciesClientUpdateTagsResponse contains the response from method SimPoliciesClient.UpdateTags.
type SimPoliciesClientUpdateTagsResponse struct {
	SimPolicy
}

// SimsClientCreateOrUpdateResponse contains the response from method SimsClient.CreateOrUpdate.
type SimsClientCreateOrUpdateResponse struct {
	Sim
}

// SimsClientDeleteResponse contains the response from method SimsClient.Delete.
type SimsClientDeleteResponse struct {
	// placeholder for future response values
}

// SimsClientGetResponse contains the response from method SimsClient.Get.
type SimsClientGetResponse struct {
	Sim
}

// SimsClientListBySimGroupResponse contains the response from method SimsClient.ListBySimGroup.
type SimsClientListBySimGroupResponse struct {
	SimListResult
}

// SitesClientCreateOrUpdateResponse contains the response from method SitesClient.CreateOrUpdate.
type SitesClientCreateOrUpdateResponse struct {
	Site
}

// SitesClientDeleteResponse contains the response from method SitesClient.Delete.
type SitesClientDeleteResponse struct {
	// placeholder for future response values
}

// SitesClientGetResponse contains the response from method SitesClient.Get.
type SitesClientGetResponse struct {
	Site
}

// SitesClientListByMobileNetworkResponse contains the response from method SitesClient.ListByMobileNetwork.
type SitesClientListByMobileNetworkResponse struct {
	SiteListResult
}

// SitesClientUpdateTagsResponse contains the response from method SitesClient.UpdateTags.
type SitesClientUpdateTagsResponse struct {
	Site
}

// SlicesClientCreateOrUpdateResponse contains the response from method SlicesClient.CreateOrUpdate.
type SlicesClientCreateOrUpdateResponse struct {
	Slice
}

// SlicesClientDeleteResponse contains the response from method SlicesClient.Delete.
type SlicesClientDeleteResponse struct {
	// placeholder for future response values
}

// SlicesClientGetResponse contains the response from method SlicesClient.Get.
type SlicesClientGetResponse struct {
	Slice
}

// SlicesClientListByMobileNetworkResponse contains the response from method SlicesClient.ListByMobileNetwork.
type SlicesClientListByMobileNetworkResponse struct {
	SliceListResult
}

// SlicesClientUpdateTagsResponse contains the response from method SlicesClient.UpdateTags.
type SlicesClientUpdateTagsResponse struct {
	Slice
}
