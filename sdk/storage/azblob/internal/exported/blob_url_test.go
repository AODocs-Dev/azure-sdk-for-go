//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseURL(t *testing.T) {
	testStorageAccount := "fakestorageaccount"
	host := fmt.Sprintf("%s.blob.core.windows.net", testStorageAccount)
	testContainer := "fakecontainer"
	fileNames := []string{"/._.TESTT.txt", "/.gitignore/dummyfile1"}

	sas := "sv=2019-12-12&sr=b&st=2111-01-09T01:42:34.936Z&se=2222-03-09T01:42:34.936Z&sp=rw&sip=168.1.5.60-168.1.5.70&spr=https,http&si=myIdentifier&ss=bf&srt=s&sig=clNxbtnkKSHw7f3KMEVVc4agaszoRFdbZr%2FWBmPNsrw%3D"

	for _, fileName := range fileNames {
		snapshotID, versionID := "", "2021-10-25T05:41:32.5526810Z"
		sasWithVersionID := "?versionId=" + versionID + "&" + sas
		urlWithVersion := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s%s", testStorageAccount, testContainer, fileName, sasWithVersionID)
		blobURLParts, err := ParseURL(urlWithVersion)
		require.NoError(t, err)

		require.Equal(t, blobURLParts.Scheme, "https")
		require.Equal(t, blobURLParts.Host, host)
		require.Equal(t, blobURLParts.ContainerName, testContainer)
		require.Equal(t, blobURLParts.VersionID, versionID)
		require.Equal(t, blobURLParts.Snapshot, snapshotID)

		validateSAS(t, sas, blobURLParts.SAS)
	}

	for _, fileName := range fileNames {
		snapshotID, versionID := "2011-03-09T01:42:34Z", ""
		sasWithSnapshotID := "?snapshot=" + snapshotID + "&" + sas
		urlWithVersion := fmt.Sprintf("https://%s.blob.core.windows.net/%s%s%s", testStorageAccount, testContainer, fileName, sasWithSnapshotID)
		blobURLParts, err := ParseURL(urlWithVersion)
		require.NoError(t, err)

		require.Equal(t, blobURLParts.Scheme, "https")
		require.Equal(t, blobURLParts.Host, host)
		require.Equal(t, blobURLParts.ContainerName, testContainer)
		require.Equal(t, blobURLParts.VersionID, versionID)
		require.Equal(t, blobURLParts.Snapshot, snapshotID)

		validateSAS(t, sas, blobURLParts.SAS)
	}

	//urlWithIP := "https://127.0.0.1:5000/"
}
