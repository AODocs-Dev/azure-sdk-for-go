//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package pageblob_test

import (
	"bytes"
	"context"
	"io"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/bloberror"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/container"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/testcommon"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/pageblob"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func Test(t *testing.T) {
	suite.Run(t, &PageBlobRecordedTestsSuite{})
	//suite.Run(t, &PageBlobUnrecordedTestsSuite{})
}

// nolint
func (s *PageBlobRecordedTestsSuite) BeforeTest(suite string, test string) {
	testcommon.BeforeTest(s.T(), suite, test)
}

// nolint
func (s *PageBlobRecordedTestsSuite) AfterTest(suite string, test string) {
	testcommon.AfterTest(s.T(), suite, test)
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) BeforeTest(suite string, test string) {

}

// nolint
func (s *PageBlobUnrecordedTestsSuite) AfterTest(suite string, test string) {

}

type PageBlobRecordedTestsSuite struct {
	suite.Suite
}

type PageBlobUnrecordedTestsSuite struct {
	suite.Suite
}

func getPageBlobClient(pageBlobName string, containerClient *container.Client) *pageblob.Client {
	return containerClient.NewPageBlobClient(pageBlobName)
}

func createNewPageBlob(ctx context.Context, _require *require.Assertions, pageBlobName string, containerClient *container.Client) *pageblob.Client {
	return createNewPageBlobWithSize(ctx, _require, pageBlobName, containerClient, pageblob.PageBytes*10)
}

func createNewPageBlobWithSize(ctx context.Context, _require *require.Assertions, pageBlobName string, containerClient *container.Client, sizeInBytes int64) *pageblob.Client {
	pbClient := getPageBlobClient(pageBlobName, containerClient)

	_, err := pbClient.Create(ctx, sizeInBytes, nil)
	_require.Nil(err)
	//_require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	return pbClient
}

func createNewPageBlobWithCPK(ctx context.Context, _require *require.Assertions, pageBlobName string, container *container.Client, sizeInBytes int64, cpkInfo *blob.CpkInfo, cpkScopeInfo *blob.CpkScopeInfo) (pbClient *pageblob.Client) {
	pbClient = getPageBlobClient(pageBlobName, container)

	_, err := pbClient.Create(ctx, sizeInBytes, &pageblob.CreateOptions{
		CpkInfo:      cpkInfo,
		CpkScopeInfo: cpkScopeInfo,
	})
	_require.Nil(err)
	// _require.Equal(resp.RawResponse.StatusCode, 201)
	return
}

func (s *PageBlobRecordedTestsSuite) TestPutGetPages() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	offset, count := int64(0), int64(1024)
	reader, _ := testcommon.GenerateData(1024)
	putResp, err := pbClient.UploadPages(context.Background(), reader, &pageblob.UploadPagesOptions{
		Offset: &offset,
		Count:  &count,
	})
	_require.Nil(err)
	_require.NotNil(putResp.LastModified)
	_require.Equal((*putResp.LastModified).IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.Nil(putResp.ContentMD5)
	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
	_require.NotNil(*putResp.RequestID)
	_require.NotNil(*putResp.Version)
	_require.NotNil(putResp.Date)
	_require.Equal((*putResp.Date).IsZero(), false)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Offset: to.Ptr(int64(0)),
		Count:  to.Ptr(int64(1023)),
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.NotNil(pageListResp.LastModified)
		_require.Equal((*pageListResp.LastModified).IsZero(), false)
		_require.NotNil(pageListResp.ETag)
		_require.Equal(*pageListResp.BlobContentLength, int64(512*10))
		_require.NotNil(*pageListResp.RequestID)
		_require.NotNil(*pageListResp.Version)
		_require.NotNil(pageListResp.Date)
		_require.Equal((*pageListResp.Date).IsZero(), false)
		_require.NotNil(pageListResp.PageList)
		pageRangeResp := pageListResp.PageList.PageRange
		_require.Len(pageRangeResp, 1)
		rawStart, rawEnd := rawPageRange((pageRangeResp)[0])
		_require.Equal(rawStart, offset)
		_require.Equal(rawEnd, count-1)
		if err != nil {
			break
		}
	}
}

//nolint
//func (s *PageBlobUnrecordedTestsSuite) TestUploadPagesFromURL() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(contentSize)
//	srcBlob := createNewPageBlobWithSize(_require, "srcblob", containerClient, int64(contentSize))
//	destBlob := createNewPageBlobWithSize(_require, "dstblob", containerClient, int64(contentSize))
//
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadSrcResp1, err := srcBlob.UploadPages(context.Background(), streaming.NopCloser(r), &pageblob.UploadPagesOptions{
//		Offset: to.Ptr(offset),
//		Count: to.Ptr(count),
//	}
//	_require.Nil(err)
//	_require.NotNil(uploadSrcResp1.LastModified)
//	_require.Equal((*uploadSrcResp1.LastModified).IsZero(), false)
//	_require.NotNil(uploadSrcResp1.ETag)
//	_require.Nil(uploadSrcResp1.ContentMD5)
//	_require.Equal(*uploadSrcResp1.BlobSequenceNumber, int64(0))
//	_require.NotNil(*uploadSrcResp1.RequestID)
//	_require.NotNil(*uploadSrcResp1.Version)
//	_require.NotNil(uploadSrcResp1.Date)
//	_require.Equal((*uploadSrcResp1.Date).IsZero(), false)
//
//	// Get source pbClient URL with SAS for UploadPagesFromURL.
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//	srcBlobParts.SAS, err = BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		_require.Error(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Upload page from URL.
//	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), nil)
//	_require.Nil(err)
//	// _require.Equal(pResp1.RawResponse.StatusCode, 201)
//	_require.NotNil(pResp1.ETag)
//	_require.NotNil(pResp1.LastModified)
//	_require.NotNil(pResp1.ContentMD5)
//	_require.NotNil(pResp1.RequestID)
//	_require.NotNil(pResp1.Version)
//	_require.NotNil(pResp1.Date)
//	_require.Equal((*pResp1.Date).IsZero(), false)
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{}))
//	_require.Nil(err)
//	_require.EqualValues(destData, sourceData)
//}
//
////nolint
//func (s *PageBlobUnrecordedTestsSuite) TestUploadPagesFromURLWithMD5() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	contentSize := 4 * 1024 * 1024 // 4MB
//	r, sourceData := getRandomDataAndReader(contentSize)
//	md5Value := md5.Sum(sourceData)
//	contentMD5 := md5Value[:]
//	ctx := ctx // Use default Background context
//	srcBlob := createNewPageBlobWithSize(_require, "srcblob", containerClient, int64(contentSize))
//	destBlob := createNewPageBlobWithSize(_require, "dstblob", containerClient, int64(contentSize))
//
//	// Prepare source pbClient for copy.
//	offset, _, count := int64(0), int64(contentSize-1), int64(contentSize)
//	uploadPagesOptions := pageblob.UploadPagesOptions{Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),}
//	_, err = srcBlob.UploadPages(context.Background(), streaming.NopCloser(r), &uploadPagesOptions)
//	_require.Nil(err)
//	// _require.Equal(uploadSrcResp1.RawResponse.StatusCode, 201)
//
//	// Get source pbClient URL with SAS for UploadPagesFromURL.
//	credential, err := getGenericCredential(nil, testcommon.TestAccountDefault)
//	_require.Nil(err)
//	srcBlobParts, _ := NewBlobURLParts(srcBlob.URL())
//
//	srcBlobParts.SAS, err = azblob.BlobSASSignatureValues{
//		Protocol:      SASProtocolHTTPS,                     // Users MUST use HTTPS (not HTTP)
//		ExpiryTime:    time.Now().UTC().Add(48 * time.Hour), // 48-hours before expiration
//		ContainerName: srcBlobParts.ContainerName,
//		BlobName:      srcBlobParts.BlobName,
//		Permissions:   BlobSASPermissions{Read: true}.String(),
//	}.Sign(credential)
//	if err != nil {
//		_require.Error(err)
//	}
//
//	srcBlobURLWithSAS := srcBlobParts.URL()
//
//	// Upload page from URL with MD5.
//	uploadPagesFromURLOptions := pageblob.UploadPagesFromURLOptions{
//		SourceContentMD5: contentMD5,
//	}
//	pResp1, err := destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_require.Nil(err)
//	// _require.Equal(pResp1.RawResponse.StatusCode, 201)
//	_require.NotNil(pResp1.ETag)
//	_require.NotNil(pResp1.LastModified)
//	_require.NotNil(pResp1.ContentMD5)
//	_require.EqualValues(pResp1.ContentMD5, contentMD5)
//	_require.NotNil(pResp1.RequestID)
//	_require.NotNil(pResp1.Version)
//	_require.NotNil(pResp1.Date)
//	_require.Equal((*pResp1.Date).IsZero(), false)
//	_require.Equal(*pResp1.BlobSequenceNumber, int64(0))
//
//	// Check data integrity through downloading.
//	downloadResp, err := destBlob.Download(ctx, nil)
//	_require.Nil(err)
//	destData, err := io.ReadAll(downloadResp.BodyReader(&blob.RetryReaderOptions{}))
//	_require.Nil(err)
//	_require.EqualValues(destData, sourceData)
//
//	// Upload page from URL with bad MD5
//	_, badMD5 := getRandomDataAndReader(16)
//	badContentMD5 := badMD5[:]
//	uploadPagesFromURLOptions = pageblob.UploadPagesFromURLOptions{
//		SourceContentMD5: badContentMD5,
//	}
//	_, err = destBlob.UploadPagesFromURL(ctx, srcBlobURLWithSAS, 0, 0, int64(contentSize), &uploadPagesFromURLOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
//}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestClearDiffPages() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	contentSize := 2 * 1024
	r := testcommon.GetReaderToGeneratedBytes(contentSize)
	_, err = pbClient.UploadPages(context.Background(), r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(0)),
		Count:  to.Ptr(int64(contentSize)),
	})
	_require.Nil(err)

	snapshotResp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	r1 := testcommon.GetReaderToGeneratedBytes(contentSize)
	_, err = pbClient.UploadPages(context.Background(), r1, &pageblob.UploadPagesOptions{Offset: to.Ptr(int64(contentSize)), Count: to.Ptr(int64(contentSize))})
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset:       to.Ptr(int64(0)),
		Count:        to.Ptr(int64(4096)),
		PrevSnapshot: snapshotResp.Snapshot,
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		pageRangeResp := pageListResp.PageList.PageRange
		_require.NotNil(pageRangeResp)
		_require.Len(pageRangeResp, 1)
		rawStart, rawEnd := rawPageRange((pageRangeResp)[0])
		_require.Equal(rawStart, int64(2048))
		_require.Equal(rawEnd, int64(4095))
		if err != nil {
			break
		}
	}

	_, err = pbClient.ClearPages(context.Background(), int64(2048), int64(2048), nil)
	_require.Nil(err)

	pager = pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset:       to.Ptr(int64(0)),
		Count:        to.Ptr(int64(4096)),
		PrevSnapshot: snapshotResp.Snapshot,
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		pageRangeResp := pageListResp.PageList.PageRange
		_require.Len(pageRangeResp, 0)
		if err != nil {
			break
		}
	}
}

// nolint
func waitForIncrementalCopy(_require *require.Assertions, copyBlobClient *pageblob.Client, blobCopyResponse *pageblob.CopyIncrementalResponse) *string {
	status := *blobCopyResponse.CopyStatus
	var getPropertiesAndMetadataResult blob.GetPropertiesResponse
	// Wait for the copy to finish
	start := time.Now()
	for status != blob.CopyStatusTypeSuccess {
		getPropertiesAndMetadataResult, _ = copyBlobClient.GetProperties(context.Background(), nil)
		status = *getPropertiesAndMetadataResult.CopyStatus
		currentTime := time.Now()
		if currentTime.Sub(start) >= time.Minute {
			_require.Fail("")
		}
	}
	return getPropertiesAndMetadataResult.DestinationSnapshot
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestIncrementalCopy() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err = containerClient.SetAccessPolicy(context.Background(), nil, &container.SetAccessPolicyOptions{Access: to.Ptr(container.PublicAccessTypeBlob)})
	_require.Nil(err)

	srcBlob := createNewPageBlob(context.Background(), _require, "src"+testcommon.GenerateBlobName(testName), containerClient)

	contentSize := 1024
	r := testcommon.GetReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize)
	_, err = srcBlob.UploadPages(context.Background(), r, &pageblob.UploadPagesOptions{Offset: to.Ptr(offset), Count: to.Ptr(count)})
	_require.Nil(err)

	snapshotResp, err := srcBlob.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	dstBlob := containerClient.NewPageBlobClient("dst" + testcommon.GenerateBlobName(testName))

	resp, err := dstBlob.StartCopyIncremental(context.Background(), srcBlob.URL(), *snapshotResp.Snapshot, nil)
	_require.Nil(err)
	_require.NotNil(resp.LastModified)
	_require.Equal((*resp.LastModified).IsZero(), false)
	_require.NotNil(resp.ETag)
	_require.NotEqual(*resp.RequestID, "")
	_require.NotEqual(*resp.Version, "")
	_require.NotNil(resp.Date)
	_require.Equal((*resp.Date).IsZero(), false)
	_require.NotEqual(*resp.CopyID, "")
	_require.Equal(*resp.CopyStatus, blob.CopyStatusTypePending)

	waitForIncrementalCopy(_require, dstBlob, &resp)
}

func (s *PageBlobRecordedTestsSuite) TestResizePageBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	_, err = pbClient.Resize(context.Background(), 2048, nil)
	_require.Nil(err)

	_, err = pbClient.Resize(context.Background(), 8192, nil)
	_require.Nil(err)

	resp2, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp2.ContentLength, int64(8192))
}

func (s *PageBlobRecordedTestsSuite) TestPageSequenceNumbers() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(0)
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	sequenceNumber = int64(7)
	actionType = pageblob.SequenceNumberActionTypeMax
	updateSequenceNumberPageBlob = pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}

	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	updateSequenceNumberPageBlob = pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: to.Ptr(int64(11)),
		ActionType:     to.Ptr(pageblob.SequenceNumberActionTypeUpdate),
	}

	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)
}

//nolint
//func (s *PageBlobUnrecordedTestsSuite) TestPutPagesWithMD5() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	blobName := testcommon.GenerateBlobName(testName)
//	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)
//
//	// put page with valid MD5
//	contentSize := 1024
//	readerToBody, body := getRandomDataAndReader(contentSize)
//	offset, _, count := int64(0), int64(0)+int64(contentSize-1), int64(contentSize)
//	md5Value := md5.Sum(body)
//	_ = body
//	contentMD5 := md5Value[:]
//
//	putResp, err := pbClient.UploadPages(context.Background(), streaming.NopCloser(readerToBody), &pageblob.UploadPagesOptions{
//		Offset:                  to.Ptr(offset),
//		Count:                   to.Ptr(count),
//		TransactionalContentMD5: contentMD5,
//	})
//	_require.Nil(err)
//	// _require.Equal(putResp.RawResponse.StatusCode, 201)
//	_require.NotNil(putResp.LastModified)
//	_require.Equal((*putResp.LastModified).IsZero(), false)
//	_require.NotNil(putResp.ETag)
//	_require.NotNil(putResp.ContentMD5)
//	_require.EqualValues(putResp.ContentMD5, contentMD5)
//	_require.Equal(*putResp.BlobSequenceNumber, int64(0))
//	_require.NotNil(*putResp.RequestID)
//	_require.NotNil(*putResp.Version)
//	_require.NotNil(putResp.Date)
//	_require.Equal((*putResp.Date).IsZero(), false)
//
//	// put page with bad MD5
//	readerToBody, _ = getRandomDataAndReader(1024)
//	_, badMD5 := getRandomDataAndReader(16)
//	basContentMD5 := badMD5[:]
//	putResp, err = pbClient.UploadPages(context.Background(), streaming.NopCloser(readerToBody), &pageblob.UploadPagesOptions{
//		Offset:                  to.Ptr(offset),
//		Count:                   to.Ptr(count),
//		TransactionalContentMD5: basContentMD5,
//	})
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.MD5Mismatch)
//}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageSizeInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(context.Background(), 1, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageSequenceInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(-1)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageMetadataNonEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageMetadataEmpty() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       map[string]string{},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Nil(resp.Metadata)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageMetadataInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       map[string]string{"In valid1": "bar"},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)
	_require.Contains(err.Error(), testcommon.InvalidHeaderErrorSubstring)

}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageHTTPHeaders() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		HTTPHeaders:    &testcommon.BasicHeaders,
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	h := blob.ParseHTTPHeaders(resp)
	_require.EqualValues(h, testcommon.BasicHeaders)
}

func validatePageBlobPut(_require *require.Assertions, pbClient *pageblob.Client) {
	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Metadata)
	_require.EqualValues(resp.Metadata, testcommon.BasicMetadata)
	_require.EqualValues(blob.ParseHTTPHeaders(resp), testcommon.BasicHeaders)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(context.Background(), pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(context.Background(), pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(context.Background(), pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResp.Date, 10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResp, err := pbClient.Create(context.Background(), pageblob.PageBytes, nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResp.Date, -10)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(0)
	eTag := "garbage"
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(0)
	eTag := "garbage"
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.Nil(err)

	validatePageBlobPut(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobCreatePageIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	sequenceNumber := int64(0)
	createPageBlobOptions := pageblob.CreateOptions{
		SequenceNumber: &sequenceNumber,
		Metadata:       testcommon.BasicMetadata,
		HTTPHeaders:    &testcommon.BasicHeaders,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Create(context.Background(), pageblob.PageBytes, &createPageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobPutPagesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	contentSize := 1024
	r := testcommon.GetReaderToGeneratedBytes(contentSize)
	offset, count := int64(0), int64(contentSize/2)
	uploadPagesOptions := pageblob.UploadPagesOptions{Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count))}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)
}

//// Body cannot be nil check already added in the request preparer
////func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesNilBody() {
////  svcClient := testcommon.GetServiceClient()
////  containerClient, _ := createNewContainer(c, svcClient)
////  defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
////  pbClient, _ := createNewPageBlob(c, containerClient)
////
////  _, err := pbClient.UploadPages(context.Background(), nil, nil)
////  _require.NotNil(err)
////}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesEmptyBody() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r := bytes.NewReader([]byte{})
	offset, count := int64(0), int64(0)
	uploadPagesOptions := pageblob.UploadPagesOptions{Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count))}
	_, err = pbClient.UploadPages(context.Background(), streaming.NopCloser(r), &uploadPagesOptions)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesNonExistentBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count))}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.BlobNotFound)
}

func validateUploadPages(_require *require.Assertions, pbClient *pageblob.Client) {
	// This will only validate a single put page at 0-PageBlobPageBytes-1
	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{
		Offset: to.Ptr(int64(0)),
		Count:  to.Ptr(int64(blob.CountToEnd)),
	})

	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)

		start, end := int64(0), int64(pageblob.PageBytes-1)
		rawStart, rawEnd := *(pageListResp.PageList.PageRange[0].Start), *(pageListResp.PageList.PageRange[0].End)
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	_, err = pbClient.UploadPages(context.Background(), r, &pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	})
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	eTag := "garbage"
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	eTag := "garbage"
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)
	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLessThanTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThan := int64(10)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLessThanFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThan := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLessThanNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}

	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLTETrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLTEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberLTENegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberEqualTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)

	validateUploadPages(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	ifSequenceNumberEqualTo := int64(1)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

//func (s *PageBlobRecordedTestsSuite) TestBlobPutPagesIfSequenceNumberEqualNegOne() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	blobName := testcommon.GenerateBlobName(testName)
//	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)
//
//	r, _ := testcommon.GenerateData(pageblob.PageBytes)
//	offset, count := int64(0), int64(pageblob.PageBytes)
//	ifSequenceNumberEqualTo := int64(-1)
//	uploadPagesOptions := pageblob.UploadPagesOptions{
//		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
//		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
//			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
//		},
//	}
//	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions) // This will cause the library to set the value of the header to 0
//	_require.Nil(err)
//}

func setupClearPagesTest(t *testing.T, _require *require.Assertions, testName string) (*container.Client, *pageblob.Client) {
	svcClient, err := testcommon.GetServiceClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)

	return containerClient, pbClient
}

func validateClearPagesTest(_require *require.Assertions, pbClient *pageblob.Client) {
	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		pageListResp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Nil(pageListResp.PageRange)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes+1), nil)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		}})
	_require.Nil(err)
	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	})
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: getPropertiesResp.ETag,
			},
		},
	}
	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	eTag := "garbage"
	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	eTag := "garbage"
	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	clearPageOptions := pageblob.ClearPagesOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLessThanTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberLessThan := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLessThanFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberLessThan := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLessThanNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberLessThan := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThan: &ifSequenceNumberLessThan,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLTETrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLTEFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberLessThanOrEqualTo := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberLTENegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberLessThanOrEqualTo := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberLessThanOrEqualTo: &ifSequenceNumberLessThanOrEqualTo,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions) // This will cause the library to set the value of the header to 0
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberEqualTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberEqualTo := int64(10)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.Nil(err)

	validateClearPagesTest(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberEqualFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	sequenceNumber := int64(10)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err := pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	ifSequenceNumberEqualTo := int64(1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err = pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.SequenceNumberConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobClearPagesIfSequenceNumberEqualNegOne() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupClearPagesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	ifSequenceNumberEqualTo := int64(-1)
	clearPageOptions := pageblob.ClearPagesOptions{
		SequenceNumberAccessConditions: &pageblob.SequenceNumberAccessConditions{
			IfSequenceNumberEqualTo: &ifSequenceNumberEqualTo,
		},
	}
	_, err := pbClient.ClearPages(context.Background(), int64(0), int64(pageblob.PageBytes), &clearPageOptions) // This will cause the library to set the value of the header to 0
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidInput)
}

func setupGetPageRangesTest(t *testing.T, _require *require.Assertions, testName string) (containerClient *container.Client, pbClient *pageblob.Client) {
	svcClient, err := testcommon.GetServiceClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient = testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient = createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)
	return
}

func validateBasicGetPageRanges(_require *require.Assertions, resp pageblob.PageList, err error) {
	_require.Nil(err)
	_require.NotNil(resp.PageRange)
	_require.Len(resp.PageRange, 1)
	start, end := int64(0), int64(pageblob.PageBytes-1)
	rawStart, rawEnd := rawPageRange((resp.PageRange)[0])
	_require.Equal(rawStart, start)
	_require.Equal(rawEnd, end)
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesEmptyBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Nil(resp.PageRange)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesEmptyRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(-2)), Count: to.Ptr(int64(500))})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.Nil(err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesNonContiguousRanges() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	r, _ := testcommon.GenerateData(pageblob.PageBytes)
	offset, count := int64(2*pageblob.PageBytes), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(int64(offset)), Count: to.Ptr(int64(count)),
	}
	_, err := pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		pageListResp := resp.PageList.PageRange
		_require.NotNil(pageListResp)
		_require.Len(pageListResp, 2)

		start, end := int64(0), int64(pageblob.PageBytes-1)
		rawStart, rawEnd := rawPageRange(pageListResp[0])
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)

		start, end = int64(pageblob.PageBytes*2), int64((pageblob.PageBytes*3)-1)
		rawStart, rawEnd = rawPageRange(pageListResp[1])
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesNotPageAligned() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(2000))})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(resp.Snapshot)

	snapshotURL, _ := pbClient.WithSnapshot(*resp.Snapshot)
	pager := snapshotURL.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0))})
	for pager.More() {
		resp2, err := pager.NextPage(context.Background())
		_require.Nil(err)

		validateBasicGetPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfModifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, 10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	getPropertiesResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	currentTime := testcommon.GetRelativeTimeFromAnchor(getPropertiesResp.Date, -10)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfUnmodifiedSince: &currentTime,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfMatch: resp.ETag,
		},
	}})
	for pager.More() {
		resp2, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfMatch: to.Ptr("garbage"),
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfNoneMatch: to.Ptr("garbage"),
		},
	}})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateBasicGetPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobGetPageRangesIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient := setupGetPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)), AccessConditions: &blob.AccessConditions{
		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
			IfNoneMatch: resp.ETag,
		},
	}})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		if err != nil {
			break
		}
	}

	//serr := err.(StorageError)
	//_require.(serr.RawResponse.StatusCode, chk.Equals, 304) // Service Code not returned in the body for a HEAD
}

// nolint
func setupDiffPageRangesTest(t *testing.T, _require *require.Assertions, testName string) (containerClient *container.Client, pbClient *pageblob.Client, snapshot string) {
	svcClient, err := testcommon.GetServiceClient(t, testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient = testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient = createNewPageBlob(context.Background(), _require, blobName, containerClient)

	r := testcommon.GetReaderToGeneratedBytes(pageblob.PageBytes)
	offset, count := int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(offset),
		Count:  to.Ptr(count),
	}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)

	resp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)
	snapshot = *resp.Snapshot

	r = testcommon.GetReaderToGeneratedBytes(pageblob.PageBytes)
	offset, count = int64(0), int64(pageblob.PageBytes)
	uploadPagesOptions = pageblob.UploadPagesOptions{Offset: to.Ptr(offset), Count: to.Ptr(count)}
	_, err = pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)
	return
}

func rawPageRange(pr *pageblob.PageRange) (start, end int64) {
	if pr.Start != nil {
		start = *pr.Start
	}
	if pr.End != nil {
		end = *pr.End
	}
	return
}

func validateDiffPageRanges(_require *require.Assertions, resp pageblob.PageList, err error) {
	_require.Nil(err)
	_require.NotNil(resp.PageRange)
	_require.Len(resp.PageRange, 1)
	rawStart, rawEnd := rawPageRange(resp.PageRange[0])
	_require.EqualValues(rawStart, int64(0))
	_require.EqualValues(rawEnd, int64(pageblob.PageBytes-1))
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangesNonExistentSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	snapshotTime, _ := time.Parse(blob.SnapshotTimeFormat, snapshot)
	snapshotTime = snapshotTime.Add(time.Minute)
	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		PrevSnapshot: to.Ptr(snapshotTime.Format(blob.SnapshotTimeFormat))})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.PreviousSnapshotNotFound)
		if err != nil {
			break
		}
	}

}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeInvalidRange() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{Offset: to.Ptr(int64(-22)), Count: to.Ptr(int64(14)), Snapshot: &snapshot})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.Nil(err)
		if err != nil {
			break
		}
	}
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeGMT(-10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfModifiedSince: &currentTime}},
	})
	for pager.More() {
		resp2, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateDiffPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeGMT(10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeGMT(10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateDiffPageRanges(_require, resp.PageList, err)
		if err != nil {
			break
		}
	}
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	currentTime := testcommon.GetRelativeTimeGMT(-10)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfUnmodifiedSince: &currentTime},
		},
	})
	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}
	}

}

////nolint
//func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	resp, err := pbClient.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//
//	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
//		Snapshot: to.Ptr(snapshot),
//		LeaseAccessConditions: &blob.LeaseAccessConditions{
//			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//				IfMatch: resp.ETag,
//			},
//		},
//	})
//	for pager.More() {
//		resp2, err := pager.NextPage(context.Background())
//		_require.Nil(err)
//		validateDiffPageRanges(_require, resp2.PageList, err)
//		if err != nil {
//			break
//		}
//	}
//}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshotStr := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		Snapshot: to.Ptr(snapshotStr),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: to.Ptr("garbage"),
			},
		}})

	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
		if err != nil {
			break
		}

	}
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshotStr := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		PrevSnapshot: to.Ptr(snapshotStr),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: to.Ptr("garbage"),
			},
		}})

	for pager.More() {
		resp2, err := pager.NextPage(context.Background())
		_require.Nil(err)
		validateDiffPageRanges(_require, resp2.PageList, err)
		if err != nil {
			break
		}
	}
}

// nolint
func (s *PageBlobUnrecordedTestsSuite) TestBlobDiffPageRangeIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	containerClient, pbClient, snapshot := setupDiffPageRangesTest(s.T(), _require, testName)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	pager := pbClient.NewGetPageRangesDiffPager(&pageblob.GetPageRangesDiffOptions{
		Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(0)),
		PrevSnapshot: to.Ptr(snapshot),
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfNoneMatch: resp.ETag},
		},
	})

	for pager.More() {
		_, err := pager.NextPage(context.Background())
		_require.NotNil(err)
		if err != nil {
			break
		}
	}
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeZero() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	// The default pbClient is created with size > 0, so this should actually update
	_, err = pbClient.Resize(context.Background(), 0, nil)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.ContentLength, int64(0))
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeInvalidSizeNegative() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	_, err = pbClient.Resize(context.Background(), -4, nil)
	_require.NotNil(err)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeInvalidSizeMisaligned() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	_, err = pbClient.Resize(context.Background(), 12, nil)
	_require.NotNil(err)
}

func validateResize(_require *require.Assertions, pbClient *pageblob.Client) {
	resp, _ := pbClient.GetProperties(context.Background(), nil)
	_require.Equal(*resp.ContentLength, int64(pageblob.PageBytes))
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	eTag := "garbage"
	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	validateResize(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	resizePageBlobOptions := pageblob.ResizeOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberActionTypeInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	sequenceNumber := int64(1)
	actionType := pageblob.SequenceNumberActionType("garbage")
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberSequenceNumberInvalid() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	defer func() { // Invalid sequence number should panic
		_ = recover()
	}()

	sequenceNumber := int64(-1)
	actionType := pageblob.SequenceNumberActionTypeUpdate
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		SequenceNumber: &sequenceNumber,
		ActionType:     &actionType,
	}

	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.InvalidHeaderValue)
}

func validateSequenceNumberSet(_require *require.Assertions, pbClient *pageblob.Client) {
	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.BlobSequenceNumber, int64(1))
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfModifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfModifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfModifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, 10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfUnmodifiedSinceFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := getPageBlobClient(blobName, containerClient)

	pageBlobCreateResponse, err := pbClient.Create(context.Background(), pageblob.PageBytes*10, nil)
	_require.Nil(err)
	// _require.Equal(pageBlobCreateResponse.RawResponse.StatusCode, 201)
	_require.NotNil(pageBlobCreateResponse.Date)

	currentTime := testcommon.GetRelativeTimeFromAnchor(pageBlobCreateResponse.Date, -10)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfUnmodifiedSince: &currentTime,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, blobName, containerClient)

	eTag := "garbage"
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfNoneMatchTrue() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, "src"+blobName, containerClient)

	eTag := "garbage"
	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: &eTag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.Nil(err)

	validateSequenceNumberSet(_require, pbClient)
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetSequenceNumberIfNoneMatchFalse() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	blobName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, "src"+blobName, containerClient)

	resp, _ := pbClient.GetProperties(context.Background(), nil)

	actionType := pageblob.SequenceNumberActionTypeIncrement
	updateSequenceNumberPageBlob := pageblob.UpdateSequenceNumberOptions{
		ActionType: &actionType,
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfNoneMatch: resp.ETag,
			},
		},
	}
	_, err = pbClient.UpdateSequenceNumber(context.Background(), &updateSequenceNumberPageBlob)
	_require.NotNil(err)

	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
}

//func setupStartIncrementalCopyTest(_require *require.Assertions, testName string) (containerClient *container.Client,
//	pbClient *pageblob.Client, copyPBClient *pageblob.Client, snapshot string) {
////	var recording *testframework.Recording
//	if _context != nil {
//		recording = _context.recording
//	}
//	svcClient, err := testcommon.GetServiceClient(recording, testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient = testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	accessType := container.PublicAccessTypeBlob
//	setAccessPolicyOptions := container.SetAccessPolicyOptions{
//		Access: &accessType,
//	}
//	_, err = containerClient.SetAccessPolicy(context.Background(), &setAccessPolicyOptions)
//	_require.Nil(err)
//
//	pbClient = createNewPageBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)
//	resp, _ := pbClient.CreateSnapshot(context.Background(), nil)
//
//	copyPBClient = getPageBlobClient("copy"+testcommon.GenerateBlobName(testName), containerClient)
//
//	// Must create the incremental copy pbClient so that the access conditions work on it
//	resp2, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), *resp.Snapshot, nil)
//	_require.Nil(err)
//	waitForIncrementalCopy(_require, copyPBClient, &resp2)
//
//	resp, _ = pbClient.CreateSnapshot(context.Background(), nil) // Take a new snapshot so the next copy will succeed
//	snapshot = *resp.Snapshot
//	return
//}

//func validateIncrementalCopy(_require *require.Assertions, copyPBClient *pageblob.Client, resp *pageblob.CopyIncrementalResponse) {
//	t := waitForIncrementalCopy(_require, copyPBClient, resp)
//
//	// If we can access the snapshot without error, we are satisfied that it was created as a result of the copy
//	copySnapshotURL, err := copyPBClient.WithSnapshot(*t)
//	_require.Nil(err)
//	_, err = copySnapshotURL.GetProperties(context.Background(), nil)
//	_require.Nil(err)
//}

//func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopySnapshotNotExist() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
////	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
//	if err != nil {
//		_require.Fail("Unable to fetch service client because " + err.Error())
//	}
//
//	containerName := testcommon.GenerateContainerName(testName)
//	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	blobName := testcommon.GenerateBlobName(testName)
//	pbClient := createNewPageBlob(context.Background(), _require, "src"+blobName, containerClient)
//	copyPBClient := getPageBlobClient("dst"+blobName, containerClient)
//
//	snapshot := time.Now().UTC().Format(blob.SnapshotTimeFormat)
//	_, err = copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, nil)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.CannotVerifyCopySource)
//}

//func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopyIfModifiedSinceTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	currentTime := testcommon.GetRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
//}
//
//func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopyIfModifiedSinceFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	currentTime := testcommon.GetRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfModifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
//}
//
//func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	currentTime := testcommon.GetRelativeTimeGMT(20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
//}
//
//func (s *PageBlobRecordedTestsSuite) TestBlobStartIncrementalCopyIfUnmodifiedSinceFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	currentTime := testcommon.GetRelativeTimeGMT(-20)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfUnmodifiedSince: &currentTime,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
//}
//
//// nolint
//func (s *PageBlobUnrecordedTestsSuite) TestBlobStartIncrementalCopyIfMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	resp, _ := copyPBClient.GetProperties(context.Background(), nil)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfMatch: resp.ETag,
//		},
//	}
//	resp2, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp2)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//}
//
////nolint
//func (s *PageBlobUnrecordedTestsSuite) TestBlobStartIncrementalCopyIfMatchFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfMatch: &eTag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.TargetConditionNotMet)
//}

////nolint
//func (s *PageBlobUnrecordedTestsSuite) TestBlobStartIncrementalCopyIfNoneMatchTrue() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	eTag := "garbage"
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfNoneMatch: &eTag,
//		},
//	}
//	resp, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.Nil(err)
//
//	validateIncrementalCopy(_require, copyPBClient, &resp)
//}

////nolint
//func (s *PageBlobUnrecordedTestsSuite) TestBlobStartIncrementalCopyIfNoneMatchFalse() {
//	_require := require.New(s.T())
//	testName := s.T().Name()
//	containerClient, pbClient, copyPBClient, snapshot := setupStartIncrementalCopyTest(_require, testName)
//	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)
//
//	resp, _ := copyPBClient.GetProperties(context.Background(), nil)
//
//	copyIncrementalPageBlobOptions := pageblob.CopyIncrementalOptions{
//		ModifiedAccessConditions: &blob.ModifiedAccessConditions{
//			IfNoneMatch: resp.ETag,
//		},
//	}
//	_, err := copyPBClient.StartCopyIncremental(context.Background(), pbClient.URL(), snapshot, &copyIncrementalPageBlobOptions)
//	_require.NotNil(err)
//
//	testcommon.ValidateBlobErrorCode(_require, err, bloberror.ConditionNotMet)
//}

func setAndCheckPageBlobTier(_require *require.Assertions, pbClient *pageblob.Client, tier blob.AccessTier) {
	_, err := pbClient.SetTier(context.Background(), tier, nil)
	_require.Nil(err)

	resp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp.AccessTier, string(tier))
}

func (s *PageBlobRecordedTestsSuite) TestBlobSetTierAllTiersOnPageBlob() {
	_require := require.New(s.T())
	testName := s.T().Name()
	premiumServiceClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountPremium, nil)
	_require.NoError(err)

	premContainerName := "prem" + testcommon.GenerateContainerName(testName)
	premContainerClient := testcommon.CreateNewContainer(context.Background(), _require, premContainerName, premiumServiceClient)
	defer testcommon.DeleteContainer(context.Background(), _require, premContainerClient)

	pbName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlob(context.Background(), _require, pbName, premContainerClient)

	possibleTiers := []blob.AccessTier{
		blob.AccessTierP4,
		blob.AccessTierP6,
		blob.AccessTierP10,
		blob.AccessTierP20,
		blob.AccessTierP30,
		blob.AccessTierP40,
		blob.AccessTierP50,
	}
	for _, possibleTier := range possibleTiers {
		setAndCheckPageBlobTier(_require, pbClient, possibleTier)
	}
}

func (s *PageBlobUnrecordedTestsSuite) TestPageBlockWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := testcommon.GenerateData(contentSize)
	pbName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(context.Background(), _require, pbName, containerClient, int64(contentSize), &testcommon.TestCPKByValue, nil)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(offset), Count: to.Ptr(count),
		CpkInfo: &testcommon.TestCPKByValue,
	}
	uploadResp, err := pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.EqualValues(uploadResp.EncryptionKeySHA256, testcommon.TestCPKByValue.EncryptionKeySHA256)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(blob.CountToEnd))})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		pageListResp := resp.PageList.PageRange
		start, end := int64(0), int64(contentSize-1)
		rawStart, rawEnd := rawPageRange(pageListResp[0])
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}

	// Get blob content without encryption key should fail the request.
	_, err = pbClient.DownloadStream(context.Background(), nil)
	_require.NotNil(err)

	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestInvalidCPKByValue,
	}
	_, err = pbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.NotNil(err)

	// Download blob to do data integrity check.
	downloadBlobOptions = blob.DownloadStreamOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	downloadResp, err := pbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Nil(err)

	destData, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionKeySHA256, *testcommon.TestCPKByValue.EncryptionKeySHA256)
}

func (s *PageBlobUnrecordedTestsSuite) TestPageBlockWithCPKScope() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	contentSize := 4 * 1024 * 1024 // 4MB
	r, srcData := testcommon.GenerateData(contentSize)
	pbName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(context.Background(), _require, pbName, containerClient, int64(contentSize), nil, &testcommon.TestCPKByScope)

	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset:       to.Ptr(int64(offset)),
		Count:        to.Ptr(int64(count)),
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	uploadResp, err := pbClient.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)
	// _require.Equal(uploadResp.RawResponse.StatusCode, 201)
	_require.EqualValues(uploadResp.EncryptionScope, testcommon.TestCPKByScope.EncryptionScope)

	pager := pbClient.NewGetPageRangesPager(&pageblob.GetPageRangesOptions{Offset: to.Ptr(int64(0)), Count: to.Ptr(int64(blob.CountToEnd))})
	for pager.More() {
		resp, err := pager.NextPage(context.Background())
		_require.Nil(err)
		pageListResp := resp.PageList.PageRange
		start, end := int64(0), int64(contentSize-1)
		rawStart, rawEnd := rawPageRange(pageListResp[0])
		_require.Equal(rawStart, start)
		_require.Equal(rawEnd, end)
		if err != nil {
			break
		}
	}

	// Download blob to do data integrity check.
	downloadBlobOptions := blob.DownloadStreamOptions{
		CpkScopeInfo: &testcommon.TestCPKByScope,
	}
	downloadResp, err := pbClient.DownloadStream(context.Background(), &downloadBlobOptions)
	_require.Nil(err)

	destData, err := io.ReadAll(downloadResp.Body)
	_require.Nil(err)
	_require.EqualValues(destData, srcData)
	_require.EqualValues(*downloadResp.EncryptionScope, *testcommon.TestCPKByScope.EncryptionScope)
}

func (s *PageBlobUnrecordedTestsSuite) TestCreatePageBlobWithTags() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClient := createNewPageBlob(context.Background(), _require, "src"+testcommon.GenerateBlobName(testName), containerClient)

	contentSize := 1 * 1024
	offset, count := int64(0), int64(contentSize)
	putResp, err := pbClient.UploadPages(context.Background(), testcommon.GetReaderToGeneratedBytes(1024), &pageblob.UploadPagesOptions{
		Offset: to.Ptr(offset), Count: to.Ptr(count)})
	_require.Nil(err)
	//_require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.Equal(putResp.LastModified.IsZero(), false)
	_require.NotEqual(putResp.ETag, "")
	_require.NotEqual(putResp.Version, "")

	_, err = pbClient.SetTags(context.Background(), testcommon.BasicBlobTagsMap, nil)
	_require.Nil(err)
	//_require.Equal(setTagResp.RawResponse.StatusCode, 204)

	gpResp, err := pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(gpResp)
	_require.Equal(*gpResp.TagCount, int64(len(testcommon.BasicBlobTagsMap)))

	blobGetTagsResponse, err := pbClient.GetTags(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(testcommon.BasicBlobTagsMap))
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.BasicBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}

	modifiedBlobTags := map[string]string{
		"a0z1u2r3e4": "b0l1o2b3",
		"b0l1o2b3":   "s0d1k2",
	}

	_, err = pbClient.SetTags(context.Background(), modifiedBlobTags, nil)
	_require.Nil(err)
	//_require.Equal(setTagResp.RawResponse.StatusCode, 204)

	gpResp, err = pbClient.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(gpResp)
	_require.Equal(*gpResp.TagCount, int64(len(modifiedBlobTags)))

	blobGetTagsResponse, err = pbClient.GetTags(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet = blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(modifiedBlobTags))
	for _, blobTag := range blobTagsSet {
		_require.Equal(modifiedBlobTags[*blobTag.Key], *blobTag.Value)
	}
}

func (s *PageBlobUnrecordedTestsSuite) TestPageBlobSetBlobTagForSnapshot() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName), svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClient := createNewPageBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	_, err = pbClient.SetTags(context.Background(), testcommon.SpecialCharBlobTagsMap, nil)
	_require.Nil(err)

	resp, err := pbClient.CreateSnapshot(context.Background(), nil)
	_require.Nil(err)

	snapshotURL, _ := pbClient.WithSnapshot(*resp.Snapshot)
	resp2, err := snapshotURL.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.Equal(*resp2.TagCount, int64(len(testcommon.SpecialCharBlobTagsMap)))

	blobGetTagsResponse, err := pbClient.GetTags(context.Background(), nil)
	_require.Nil(err)
	// _require.Equal(blobGetTagsResponse.RawResponse.StatusCode, 200)
	blobTagsSet := blobGetTagsResponse.BlobTagSet
	_require.NotNil(blobTagsSet)
	_require.Len(blobTagsSet, len(testcommon.SpecialCharBlobTagsMap))
	for _, blobTag := range blobTagsSet {
		_require.Equal(testcommon.SpecialCharBlobTagsMap[*blobTag.Key], *blobTag.Value)
	}
}

func (s *PageBlobRecordedTestsSuite) TestCreatePageBlobReturnsVID() {
	_require := require.New(s.T())
	testName := s.T().Name()

	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)

	containerName := testcommon.GenerateContainerName(testName)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, containerName, svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbClob := createNewPageBlob(context.Background(), _require, testcommon.GenerateBlobName(testName), containerClient)

	contentSize := 1 * 1024
	r, _ := testcommon.GenerateData(contentSize)
	offset, count := int64(0), int64(contentSize)
	uploadPagesOptions := pageblob.UploadPagesOptions{
		Offset: to.Ptr(offset),
		Count:  to.Ptr(count),
	}
	putResp, err := pbClob.UploadPages(context.Background(), r, &uploadPagesOptions)
	_require.Nil(err)
	//_require.Equal(putResp.RawResponse.StatusCode, 201)
	_require.Equal(putResp.LastModified.IsZero(), false)
	_require.NotNil(putResp.ETag)
	_require.NotEqual(*putResp.Version, "")

	gpResp, err := pbClob.GetProperties(context.Background(), nil)
	_require.Nil(err)
	_require.NotNil(gpResp)
}

func (s *PageBlobRecordedTestsSuite) TestBlobResizeWithCPK() {
	_require := require.New(s.T())
	testName := s.T().Name()
	svcClient, err := testcommon.GetServiceClient(s.T(), testcommon.TestAccountDefault, nil)
	_require.NoError(err)
	containerClient := testcommon.CreateNewContainer(context.Background(), _require, testcommon.GenerateContainerName(testName)+"01", svcClient)
	defer testcommon.DeleteContainer(context.Background(), _require, containerClient)

	pbName := testcommon.GenerateBlobName(testName)
	pbClient := createNewPageBlobWithCPK(context.Background(), _require, pbName, containerClient, pageblob.PageBytes*10, &testcommon.TestCPKByValue, nil)

	resizePageBlobOptions := pageblob.ResizeOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	_, err = pbClient.Resize(context.Background(), pageblob.PageBytes, &resizePageBlobOptions)
	_require.Nil(err)

	getBlobPropertiesOptions := blob.GetPropertiesOptions{
		CpkInfo: &testcommon.TestCPKByValue,
	}
	resp, _ := pbClient.GetProperties(context.Background(), &getBlobPropertiesOptions)
	_require.Equal(*resp.ContentLength, int64(pageblob.PageBytes))
}
