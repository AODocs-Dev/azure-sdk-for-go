//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

// ClientOptions contains the optional parameters when creating a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// Client represents a client to an Azure Storage page blob;
type Client base.CompositeClient[generated.BlobClient, generated.PageBlobClient]

// NewClient creates a ServiceClient object using the specified URL, Azure AD credential, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewClient(blobURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{shared.TokenScope}, nil)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewPageBlobClient(blobURL, pl, nil)), nil
}

// NewClientWithNoCredential creates a ServiceClient object using the specified URL and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net?<SAS token>
func NewClientWithNoCredential(blobURL string, options *ClientOptions) (*Client, error) {
	conOptions := shared.GetClientOptions(options)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewPageBlobClient(blobURL, pl, nil)), nil
}

// NewClientWithSharedKeyCredential creates a ServiceClient object using the specified URL, shared key, and options.
// Example of serviceURL: https://<your_storage_account>.blob.core.windows.net
func NewClientWithSharedKeyCredential(blobURL string, cred *blob.SharedKeyCredential, options *ClientOptions) (*Client, error) {
	authPolicy := exported.NewSharedKeyCredPolicy(cred)
	conOptions := shared.GetClientOptions(options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, authPolicy)
	pl := runtime.NewPipeline(exported.ModuleName, exported.ModuleVersion, runtime.PipelineOptions{}, &conOptions.ClientOptions)

	return (*Client)(base.NewPageBlobClient(blobURL, pl, cred)), nil
}

// NewClientFromConnectionString creates Client from a connection String
func NewClientFromConnectionString(connectionString, containerName, blobName string, options *ClientOptions) (*Client, error) {
	parsed, err := shared.ParseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	parsed.ServiceURL = runtime.JoinPaths(parsed.ServiceURL, containerName, blobName)

	if parsed.AccountKey != "" && parsed.AccountName != "" {
		credential, err := exported.NewSharedKeyCredential(parsed.AccountName, parsed.AccountKey)
		if err != nil {
			return nil, err
		}
		return NewClientWithSharedKeyCredential(parsed.ServiceURL, credential, options)
	}

	return NewClientWithNoCredential(parsed.ServiceURL, options)
}

func (pb *Client) generated() *generated.PageBlobClient {
	_, pageBlob := base.InnerClients((*base.CompositeClient[generated.BlobClient, generated.PageBlobClient])(pb))
	return pageBlob
}

// URL returns the URL endpoint used by the Client object.
func (pb *Client) URL() string {
	return pb.generated().Endpoint()
}

// BlobClient returns the embedded blob client for this AppendBlob client.
func (pb *Client) BlobClient() *blob.Client {
	innerBlob, _ := base.InnerClients((*base.CompositeClient[generated.BlobClient, generated.PageBlobClient])(pb))
	return (*blob.Client)(innerBlob)
}

func (pb *Client) sharedKey() *blob.SharedKeyCredential {
	return base.SharedKeyComposite((*base.CompositeClient[generated.BlobClient, generated.PageBlobClient])(pb))
}

// WithSnapshot creates a new PageBlobURL object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the snapshot returning a URL to the base blob.
func (pb *Client) WithSnapshot(snapshot string) (*Client, error) {
	p, err := exported.ParseURL(pb.URL())
	if err != nil {
		return nil, err
	}
	p.Snapshot = snapshot

	return (*Client)(base.NewPageBlobClient(p.String(), pb.generated().Pipeline(), pb.sharedKey())), nil
}

// WithVersionID creates a new PageBlobURL object identical to the source but with the specified snapshot timestamp.
// Pass "" to remove the version returning a URL to the base blob.
func (pb *Client) WithVersionID(versionID string) (*Client, error) {
	p, err := exported.ParseURL(pb.URL())
	if err != nil {
		return nil, err
	}
	p.VersionID = versionID

	return (*Client)(base.NewPageBlobClient(p.String(), pb.generated().Pipeline(), pb.sharedKey())), nil
}

// Create creates a page blob of the specified length. Call PutPage to upload data to a page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-blob.
func (pb *Client) Create(ctx context.Context, size int64, o *CreateOptions) (CreateResponse, error) {
	createOptions, HTTPHeaders, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions := o.format()

	resp, err := pb.generated().Create(ctx, 0, size, createOptions, HTTPHeaders,
		leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)
	return resp, err
}

// UploadPages writes 1 or more pages to the page blob. The start offset and the stream size must be a multiple of 512 bytes.
// This method panics if the stream is not at position 0.
// Note that the http client closes the body stream after the request is sent to the service.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-page.
func (pb *Client) UploadPages(ctx context.Context, body io.ReadSeekCloser, options *UploadPagesOptions) (UploadPagesResponse, error) {
	count, err := shared.ValidateSeekableStreamAt0AndGetCount(body)

	if err != nil {
		return UploadPagesResponse{}, err
	}

	uploadPagesOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions := options.format()

	resp, err := pb.generated().UploadPages(ctx, count, body, uploadPagesOptions, leaseAccessConditions,
		cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions)

	return resp, err
}

// UploadPagesFromURL copies 1 or more pages from a source URL to the page blob.
// The sourceOffset specifies the start offset of source data to copy from.
// The destOffset specifies the start offset of data in page blob will be written to.
// The count must be a multiple of 512 bytes.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-page-from-url.
func (pb *Client) UploadPagesFromURL(ctx context.Context, source string, sourceOffset, destOffset, count int64,
	o *UploadPagesFromURLOptions) (UploadPagesFromURLResponse, error) {

	uploadPagesFromURLOptions, cpkInfo, cpkScopeInfo, leaseAccessConditions, sequenceNumberAccessConditions,
		modifiedAccessConditions, sourceModifiedAccessConditions := o.format()

	resp, err := pb.generated().UploadPagesFromURL(ctx, source, shared.RangeToString(sourceOffset, count), 0,
		shared.RangeToString(destOffset, count), uploadPagesFromURLOptions, cpkInfo, cpkScopeInfo, leaseAccessConditions,
		sequenceNumberAccessConditions, modifiedAccessConditions, sourceModifiedAccessConditions)

	return resp, err
}

// ClearPages frees the specified pages from the page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/put-page.
func (pb *Client) ClearPages(ctx context.Context, offset, count int64, options *ClearPagesOptions) (ClearPagesResponse, error) {
	clearOptions := &generated.PageBlobClientClearPagesOptions{
		Range: shared.HTTPRange{Offset: offset, Count: count}.Format(),
	}

	leaseAccessConditions, cpkInfo, cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions := options.format()

	resp, err := pb.generated().ClearPages(ctx, 0, clearOptions, leaseAccessConditions, cpkInfo,
		cpkScopeInfo, sequenceNumberAccessConditions, modifiedAccessConditions)

	return resp, err
}

// NewGetPageRangesPager returns the list of valid page ranges for a page blob or snapshot of a page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-page-ranges.
func (pb *Client) NewGetPageRangesPager(o *GetPageRangesOptions) *runtime.Pager[GetPageRangesResponse] {
	opts, leaseAccessConditions, modifiedAccessConditions := o.format()

	return runtime.NewPager(runtime.PagingHandler[GetPageRangesResponse]{
		More: func(page GetPageRangesResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *GetPageRangesResponse) (GetPageRangesResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = pb.generated().GetPageRangesCreateRequest(ctx, opts, leaseAccessConditions, modifiedAccessConditions)
			} else {
				opts.Marker = page.NextMarker
				req, err = pb.generated().GetPageRangesCreateRequest(ctx, opts, leaseAccessConditions, modifiedAccessConditions)
			}
			if err != nil {
				return GetPageRangesResponse{}, err
			}
			resp, err := pb.generated().Pipeline().Do(req)
			if err != nil {
				return GetPageRangesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return GetPageRangesResponse{}, runtime.NewResponseError(resp)
			}
			return pb.generated().GetPageRangesHandleResponse(resp)
		},
	})
}

// NewGetPageRangesDiffPager gets the collection of page ranges that differ between a specified snapshot and this page blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-page-ranges.
func (pb *Client) NewGetPageRangesDiffPager(o *GetPageRangesDiffOptions) *runtime.Pager[GetPageRangesDiffResponse] {
	opts, leaseAccessConditions, modifiedAccessConditions := o.format()

	return runtime.NewPager(runtime.PagingHandler[GetPageRangesDiffResponse]{
		More: func(page GetPageRangesDiffResponse) bool {
			return page.NextMarker != nil && len(*page.NextMarker) > 0
		},
		Fetcher: func(ctx context.Context, page *GetPageRangesDiffResponse) (GetPageRangesDiffResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = pb.generated().GetPageRangesDiffCreateRequest(ctx, opts, leaseAccessConditions, modifiedAccessConditions)
			} else {
				opts.Marker = page.NextMarker
				req, err = pb.generated().GetPageRangesDiffCreateRequest(ctx, opts, leaseAccessConditions, modifiedAccessConditions)
			}
			if err != nil {
				return GetPageRangesDiffResponse{}, err
			}
			resp, err := pb.generated().Pipeline().Do(req)
			if err != nil {
				return GetPageRangesDiffResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return GetPageRangesDiffResponse{}, runtime.NewResponseError(resp)
			}
			return pb.generated().GetPageRangesDiffHandleResponse(resp)
		},
	})
}

// Resize resizes the page blob to the specified size (which must be a multiple of 512).
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-blob-properties.
func (pb *Client) Resize(ctx context.Context, size int64, options *ResizeOptions) (ResizeResponse, error) {
	resizeOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions := options.format()

	resp, err := pb.generated().Resize(ctx, size, resizeOptions, leaseAccessConditions, cpkInfo, cpkScopeInfo, modifiedAccessConditions)

	return resp, err
}

// UpdateSequenceNumber sets the page blob's sequence number.
func (pb *Client) UpdateSequenceNumber(ctx context.Context, options *UpdateSequenceNumberOptions) (UpdateSequenceNumberResponse, error) {
	actionType, updateOptions, lac, mac := options.format()
	resp, err := pb.generated().UpdateSequenceNumber(ctx, *actionType, updateOptions, lac, mac)

	return resp, err
}

// StartCopyIncremental begins an operation to start an incremental copy from one-page blob's snapshot to this page blob.
// The snapshot is copied such that only the differential changes between the previously copied snapshot are transferred to the destination.
// The copied snapshots are complete copies of the original snapshot and can be read or copied from as usual.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/incremental-copy-blob and
// https://docs.microsoft.com/en-us/azure/virtual-machines/windows/incremental-snapshots.
func (pb *Client) StartCopyIncremental(ctx context.Context, copySource string, prevSnapshot string, options *CopyIncrementalOptions) (CopyIncrementalResponse, error) {
	copySourceURL, err := url.Parse(copySource)
	if err != nil {
		return CopyIncrementalResponse{}, err
	}

	queryParams := copySourceURL.Query()
	queryParams.Set("snapshot", prevSnapshot)
	copySourceURL.RawQuery = queryParams.Encode()

	pageBlobCopyIncrementalOptions, modifiedAccessConditions := options.format()
	resp, err := pb.generated().CopyIncremental(ctx, copySourceURL.String(), pageBlobCopyIncrementalOptions, modifiedAccessConditions)

	return resp, err
}

// Redeclared APIs

// Delete marks the specified blob or snapshot for deletion. The blob is later deleted during garbage collection.
// Note that deleting a blob also deletes all its snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/delete-blob.
func (pb *Client) Delete(ctx context.Context, o *blob.DeleteOptions) (blob.DeleteResponse, error) {
	return pb.BlobClient().Delete(ctx, o)
}

// Undelete restores the contents and metadata of a soft-deleted blob and any associated soft-deleted snapshots.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/undelete-blob.
func (pb *Client) Undelete(ctx context.Context, o *blob.UndeleteOptions) (blob.UndeleteResponse, error) {
	return pb.BlobClient().Undelete(ctx, o)
}

// SetTier operation sets the tier on a blob. The operation is allowed on a page
// blob in a premium storage account and on a block blob in a blob storage account (locally
// redundant storage only). A premium page blob's tier determines the allowed size, IOPS, and
// bandwidth of the blob. A block blob's tier determines Hot/Cool/Archive storage type. This operation
// does not update the blob's ETag.
// For detailed information about block blob level tier-ing see https://docs.microsoft.com/en-us/azure/storage/blobs/storage-blob-storage-tiers.
func (pb *Client) SetTier(ctx context.Context, tier blob.AccessTier, o *blob.SetTierOptions) (blob.SetTierResponse, error) {
	return pb.BlobClient().SetTier(ctx, tier, o)
}

// GetProperties returns the blob's properties.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob-properties.
func (pb *Client) GetProperties(ctx context.Context, o *blob.GetPropertiesOptions) (blob.GetPropertiesResponse, error) {
	return pb.BlobClient().GetProperties(ctx, o)
}

// SetHTTPHeaders changes a blob's HTTP headers.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/set-blob-properties.
func (pb *Client) SetHTTPHeaders(ctx context.Context, HTTPHeaders blob.HTTPHeaders, o *blob.SetHTTPHeadersOptions) (blob.SetHTTPHeadersResponse, error) {
	return pb.BlobClient().SetHTTPHeaders(ctx, HTTPHeaders, o)
}

// SetMetadata changes a blob's metadata.
// https://docs.microsoft.com/rest/api/storageservices/set-blob-metadata.
func (pb *Client) SetMetadata(ctx context.Context, metadata map[string]string, o *blob.SetMetadataOptions) (blob.SetMetadataResponse, error) {
	return pb.BlobClient().SetMetadata(ctx, metadata, o)
}

// CreateSnapshot creates a read-only snapshot of a blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/snapshot-blob.
func (pb *Client) CreateSnapshot(ctx context.Context, o *blob.CreateSnapshotOptions) (blob.CreateSnapshotResponse, error) {
	return pb.BlobClient().CreateSnapshot(ctx, o)
}

// StartCopyFromURL copies the data at the source URL to a blob.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/copy-blob.
func (pb *Client) StartCopyFromURL(ctx context.Context, copySource string, o *blob.StartCopyFromURLOptions) (blob.StartCopyFromURLResponse, error) {
	return pb.BlobClient().StartCopyFromURL(ctx, copySource, o)
}

// AbortCopyFromURL stops a pending copy that was previously started and leaves a destination blob with 0 length and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/abort-copy-blob.
func (pb *Client) AbortCopyFromURL(ctx context.Context, copyID string, o *blob.AbortCopyFromURLOptions) (blob.AbortCopyFromURLResponse, error) {
	return pb.BlobClient().AbortCopyFromURL(ctx, copyID, o)
}

// SetTags operation enables users to set tags on a blob or specific blob version, but not snapshot.
// Each call to this operation replaces all existing tags attached to the blob.
// To remove all tags from the blob, call this operation with no tags set.
// https://docs.microsoft.com/en-us/rest/api/storageservices/set-blob-tags
func (pb *Client) SetTags(ctx context.Context, tags map[string]string, o *blob.SetTagsOptions) (blob.SetTagsResponse, error) {
	return pb.BlobClient().SetTags(ctx, tags, o)
}

// GetTags operation enables users to get tags on a blob or specific blob version, or snapshot.
// https://docs.microsoft.com/en-us/rest/api/storageservices/get-blob-tags
func (pb *Client) GetTags(ctx context.Context, o *blob.GetTagsOptions) (blob.GetTagsResponse, error) {
	return pb.BlobClient().GetTags(ctx, o)
}

// CopyFromURL synchronously copies the data at the source URL to a block blob, with sizes up to 256 MB.
// For more information, see https://docs.microsoft.com/en-us/rest/api/storageservices/copy-blob-from-url.
func (pb *Client) CopyFromURL(ctx context.Context, copySource string, o *blob.CopyFromURLOptions) (blob.CopyFromURLResponse, error) {
	return pb.BlobClient().CopyFromURL(ctx, copySource, o)
}

// Concurrent Download Functions -----------------------------------------------------------------------------------------

// DownloadStream reads a range of bytes from a blob. The response also includes the blob's properties and metadata.
// For more information, see https://docs.microsoft.com/rest/api/storageservices/get-blob.
func (pb *Client) DownloadStream(ctx context.Context, o *blob.DownloadStreamOptions) (blob.DownloadStreamResponse, error) {
	return pb.BlobClient().DownloadStream(ctx, o)
}

// DownloadBuffer downloads an Azure blob to a buffer with parallel.
func (pb *Client) DownloadBuffer(ctx context.Context, buffer []byte, o *blob.DownloadBufferOptions) (int64, error) {
	return pb.BlobClient().DownloadBuffer(ctx, shared.NewBytesWriter(buffer), o)
}

// DownloadFile downloads an Azure blob to a local file.
// The file would be truncated if the size doesn't match.
func (pb *Client) DownloadFile(ctx context.Context, file *os.File, o *blob.DownloadFileOptions) (int64, error) {
	return pb.BlobClient().DownloadFile(ctx, file, o)
}
