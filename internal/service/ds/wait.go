package ds

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/directoryservice"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

const (
	directoryCreatedTimeout      = 60 * time.Minute
	directoryDeletedTimeout      = 60 * time.Minute
	directoryShareDeletedTimeout = 60 * time.Minute
)

func waitDirectoryCreated(conn *directoryservice.DirectoryService, id string) (*directoryservice.DirectoryDescription, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{directoryservice.DirectoryStageRequested, directoryservice.DirectoryStageCreating, directoryservice.DirectoryStageCreated},
		Target:  []string{directoryservice.DirectoryStageActive},
		Refresh: statusDirectoryStage(conn, id),
		Timeout: directoryCreatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*directoryservice.DirectoryDescription); ok {
		tfresource.SetLastError(err, errors.New(aws.StringValue(output.StageReason)))

		return output, err
	}

	return nil, err
}

func waitDirectoryDeleted(conn *directoryservice.DirectoryService, id string) (*directoryservice.DirectoryDescription, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{directoryservice.DirectoryStageActive, directoryservice.DirectoryStageDeleting},
		Target:  []string{},
		Refresh: statusDirectoryStage(conn, id),
		Timeout: directoryDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*directoryservice.DirectoryDescription); ok {
		tfresource.SetLastError(err, errors.New(aws.StringValue(output.StageReason)))

		return output, err
	}

	return nil, err
}

func waitDirectoryShareDeleted(ctx context.Context, conn *directoryservice.DirectoryService, ownerId, sharedId string) (*directoryservice.SharedDirectory, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{
			directoryservice.ShareStatusDeleting,
			directoryservice.ShareStatusShared,
			directoryservice.ShareStatusPendingAcceptance,
			directoryservice.ShareStatusRejectFailed,
			directoryservice.ShareStatusRejected,
			directoryservice.ShareStatusRejecting,
		},
		Target:  []string{},
		Refresh: statusSharedDirectory(ctx, conn, ownerId, sharedId),
		Timeout: directoryShareDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForStateContext(ctx)

	if output, ok := outputRaw.(*directoryservice.SharedDirectory); ok {
		tfresource.SetLastError(err, errors.New(aws.StringValue(output.ShareStatus)))

		return output, err
	}

	return nil, err
}
