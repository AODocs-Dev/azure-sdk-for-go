package hybridkubernetes

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// AuthenticationMethod enumerates the values for authentication method.
type AuthenticationMethod string

const (
	// AuthenticationMethodAAD ...
	AuthenticationMethodAAD AuthenticationMethod = "AAD"
	// AuthenticationMethodToken ...
	AuthenticationMethodToken AuthenticationMethod = "Token"
)

// PossibleAuthenticationMethodValues returns an array of possible values for the AuthenticationMethod const type.
func PossibleAuthenticationMethodValues() []AuthenticationMethod {
	return []AuthenticationMethod{AuthenticationMethodAAD, AuthenticationMethodToken}
}

// ConnectivityStatus enumerates the values for connectivity status.
type ConnectivityStatus string

const (
	// ConnectivityStatusConnected ...
	ConnectivityStatusConnected ConnectivityStatus = "Connected"
	// ConnectivityStatusConnecting ...
	ConnectivityStatusConnecting ConnectivityStatus = "Connecting"
	// ConnectivityStatusExpired ...
	ConnectivityStatusExpired ConnectivityStatus = "Expired"
	// ConnectivityStatusOffline ...
	ConnectivityStatusOffline ConnectivityStatus = "Offline"
)

// PossibleConnectivityStatusValues returns an array of possible values for the ConnectivityStatus const type.
func PossibleConnectivityStatusValues() []ConnectivityStatus {
	return []ConnectivityStatus{ConnectivityStatusConnected, ConnectivityStatusConnecting, ConnectivityStatusExpired, ConnectivityStatusOffline}
}

// CreatedByType enumerates the values for created by type.
type CreatedByType string

const (
	// CreatedByTypeApplication ...
	CreatedByTypeApplication CreatedByType = "Application"
	// CreatedByTypeKey ...
	CreatedByTypeKey CreatedByType = "Key"
	// CreatedByTypeManagedIdentity ...
	CreatedByTypeManagedIdentity CreatedByType = "ManagedIdentity"
	// CreatedByTypeUser ...
	CreatedByTypeUser CreatedByType = "User"
)

// PossibleCreatedByTypeValues returns an array of possible values for the CreatedByType const type.
func PossibleCreatedByTypeValues() []CreatedByType {
	return []CreatedByType{CreatedByTypeApplication, CreatedByTypeKey, CreatedByTypeManagedIdentity, CreatedByTypeUser}
}

// LastModifiedByType enumerates the values for last modified by type.
type LastModifiedByType string

const (
	// LastModifiedByTypeApplication ...
	LastModifiedByTypeApplication LastModifiedByType = "Application"
	// LastModifiedByTypeKey ...
	LastModifiedByTypeKey LastModifiedByType = "Key"
	// LastModifiedByTypeManagedIdentity ...
	LastModifiedByTypeManagedIdentity LastModifiedByType = "ManagedIdentity"
	// LastModifiedByTypeUser ...
	LastModifiedByTypeUser LastModifiedByType = "User"
)

// PossibleLastModifiedByTypeValues returns an array of possible values for the LastModifiedByType const type.
func PossibleLastModifiedByTypeValues() []LastModifiedByType {
	return []LastModifiedByType{LastModifiedByTypeApplication, LastModifiedByTypeKey, LastModifiedByTypeManagedIdentity, LastModifiedByTypeUser}
}

// ProvisioningState enumerates the values for provisioning state.
type ProvisioningState string

const (
	// ProvisioningStateAccepted ...
	ProvisioningStateAccepted ProvisioningState = "Accepted"
	// ProvisioningStateCanceled ...
	ProvisioningStateCanceled ProvisioningState = "Canceled"
	// ProvisioningStateDeleting ...
	ProvisioningStateDeleting ProvisioningState = "Deleting"
	// ProvisioningStateFailed ...
	ProvisioningStateFailed ProvisioningState = "Failed"
	// ProvisioningStateProvisioning ...
	ProvisioningStateProvisioning ProvisioningState = "Provisioning"
	// ProvisioningStateSucceeded ...
	ProvisioningStateSucceeded ProvisioningState = "Succeeded"
	// ProvisioningStateUpdating ...
	ProvisioningStateUpdating ProvisioningState = "Updating"
)

// PossibleProvisioningStateValues returns an array of possible values for the ProvisioningState const type.
func PossibleProvisioningStateValues() []ProvisioningState {
	return []ProvisioningState{ProvisioningStateAccepted, ProvisioningStateCanceled, ProvisioningStateDeleting, ProvisioningStateFailed, ProvisioningStateProvisioning, ProvisioningStateSucceeded, ProvisioningStateUpdating}
}

// ResourceIdentityType enumerates the values for resource identity type.
type ResourceIdentityType string

const (
	// ResourceIdentityTypeNone ...
	ResourceIdentityTypeNone ResourceIdentityType = "None"
	// ResourceIdentityTypeSystemAssigned ...
	ResourceIdentityTypeSystemAssigned ResourceIdentityType = "SystemAssigned"
)

// PossibleResourceIdentityTypeValues returns an array of possible values for the ResourceIdentityType const type.
func PossibleResourceIdentityTypeValues() []ResourceIdentityType {
	return []ResourceIdentityType{ResourceIdentityTypeNone, ResourceIdentityTypeSystemAssigned}
}
