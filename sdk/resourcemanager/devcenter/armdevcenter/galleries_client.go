//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdevcenter

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// GalleriesClient contains the methods for the Galleries group.
// Don't use this type directly, use NewGalleriesClient() instead.
type GalleriesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewGalleriesClient creates a new instance of GalleriesClient with the specified values.
// subscriptionID - Unique identifier of the Azure subscription. This is a GUID-formatted string (e.g. 00000000-0000-0000-0000-000000000000).
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewGalleriesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*GalleriesClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &GalleriesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a gallery.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-08-01-preview
// resourceGroupName - Name of the resource group within the Azure subscription.
// devCenterName - The name of the devcenter.
// galleryName - The name of the gallery.
// body - Represents a gallery.
// options - GalleriesClientBeginCreateOrUpdateOptions contains the optional parameters for the GalleriesClient.BeginCreateOrUpdate
// method.
func (client *GalleriesClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, devCenterName string, galleryName string, body Gallery, options *GalleriesClientBeginCreateOrUpdateOptions) (*runtime.Poller[GalleriesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, devCenterName, galleryName, body, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[GalleriesClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[GalleriesClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - Creates or updates a gallery.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-08-01-preview
func (client *GalleriesClient) createOrUpdate(ctx context.Context, resourceGroupName string, devCenterName string, galleryName string, body Gallery, options *GalleriesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, devCenterName, galleryName, body, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *GalleriesClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, devCenterName string, galleryName string, body Gallery, options *GalleriesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devcenters/{devCenterName}/galleries/{galleryName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if devCenterName == "" {
		return nil, errors.New("parameter devCenterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{devCenterName}", url.PathEscape(devCenterName))
	if galleryName == "" {
		return nil, errors.New("parameter galleryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// BeginDelete - Deletes a gallery resource.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-08-01-preview
// resourceGroupName - Name of the resource group within the Azure subscription.
// devCenterName - The name of the devcenter.
// galleryName - The name of the gallery.
// options - GalleriesClientBeginDeleteOptions contains the optional parameters for the GalleriesClient.BeginDelete method.
func (client *GalleriesClient) BeginDelete(ctx context.Context, resourceGroupName string, devCenterName string, galleryName string, options *GalleriesClientBeginDeleteOptions) (*runtime.Poller[GalleriesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, devCenterName, galleryName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[GalleriesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
	} else {
		return runtime.NewPollerFromResumeToken[GalleriesClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Deletes a gallery resource.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-08-01-preview
func (client *GalleriesClient) deleteOperation(ctx context.Context, resourceGroupName string, devCenterName string, galleryName string, options *GalleriesClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, devCenterName, galleryName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *GalleriesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, devCenterName string, galleryName string, options *GalleriesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devcenters/{devCenterName}/galleries/{galleryName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if devCenterName == "" {
		return nil, errors.New("parameter devCenterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{devCenterName}", url.PathEscape(devCenterName))
	if galleryName == "" {
		return nil, errors.New("parameter galleryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets a gallery
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-08-01-preview
// resourceGroupName - Name of the resource group within the Azure subscription.
// devCenterName - The name of the devcenter.
// galleryName - The name of the gallery.
// options - GalleriesClientGetOptions contains the optional parameters for the GalleriesClient.Get method.
func (client *GalleriesClient) Get(ctx context.Context, resourceGroupName string, devCenterName string, galleryName string, options *GalleriesClientGetOptions) (GalleriesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, devCenterName, galleryName, options)
	if err != nil {
		return GalleriesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return GalleriesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return GalleriesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *GalleriesClient) getCreateRequest(ctx context.Context, resourceGroupName string, devCenterName string, galleryName string, options *GalleriesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devcenters/{devCenterName}/galleries/{galleryName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if devCenterName == "" {
		return nil, errors.New("parameter devCenterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{devCenterName}", url.PathEscape(devCenterName))
	if galleryName == "" {
		return nil, errors.New("parameter galleryName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{galleryName}", url.PathEscape(galleryName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *GalleriesClient) getHandleResponse(resp *http.Response) (GalleriesClientGetResponse, error) {
	result := GalleriesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Gallery); err != nil {
		return GalleriesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByDevCenterPager - Lists galleries for a devcenter.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-08-01-preview
// resourceGroupName - Name of the resource group within the Azure subscription.
// devCenterName - The name of the devcenter.
// options - GalleriesClientListByDevCenterOptions contains the optional parameters for the GalleriesClient.ListByDevCenter
// method.
func (client *GalleriesClient) NewListByDevCenterPager(resourceGroupName string, devCenterName string, options *GalleriesClientListByDevCenterOptions) *runtime.Pager[GalleriesClientListByDevCenterResponse] {
	return runtime.NewPager(runtime.PagingHandler[GalleriesClientListByDevCenterResponse]{
		More: func(page GalleriesClientListByDevCenterResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *GalleriesClientListByDevCenterResponse) (GalleriesClientListByDevCenterResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByDevCenterCreateRequest(ctx, resourceGroupName, devCenterName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return GalleriesClientListByDevCenterResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return GalleriesClientListByDevCenterResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return GalleriesClientListByDevCenterResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByDevCenterHandleResponse(resp)
		},
	})
}

// listByDevCenterCreateRequest creates the ListByDevCenter request.
func (client *GalleriesClient) listByDevCenterCreateRequest(ctx context.Context, resourceGroupName string, devCenterName string, options *GalleriesClientListByDevCenterOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devcenters/{devCenterName}/galleries"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if devCenterName == "" {
		return nil, errors.New("parameter devCenterName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{devCenterName}", url.PathEscape(devCenterName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByDevCenterHandleResponse handles the ListByDevCenter response.
func (client *GalleriesClient) listByDevCenterHandleResponse(resp *http.Response) (GalleriesClientListByDevCenterResponse, error) {
	result := GalleriesClientListByDevCenterResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GalleryListResult); err != nil {
		return GalleriesClientListByDevCenterResponse{}, err
	}
	return result, nil
}
