// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package rds

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindDBProxyTarget(ctx context.Context, conn *rds.RDS, dbProxyName, targetGroupName, targetType, rdsResourceId string) (*rds.DBProxyTarget, error) {
	input := &rds.DescribeDBProxyTargetsInput{
		DBProxyName:     aws.String(dbProxyName),
		TargetGroupName: aws.String(targetGroupName),
	}
	var dbProxyTarget *rds.DBProxyTarget

	err := conn.DescribeDBProxyTargetsPagesWithContext(ctx, input, func(page *rds.DescribeDBProxyTargetsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, target := range page.Targets {
			if aws.StringValue(target.Type) == targetType && aws.StringValue(target.RdsResourceId) == rdsResourceId {
				dbProxyTarget = target
				return false
			}
		}

		return !lastPage
	})

	return dbProxyTarget, err
}

func FindDBProxyEndpoint(ctx context.Context, conn *rds.RDS, id string) (*rds.DBProxyEndpoint, error) {
	dbProxyName, dbProxyEndpointName, err := ProxyEndpointParseID(id)
	if err != nil {
		return nil, err
	}

	input := &rds.DescribeDBProxyEndpointsInput{
		DBProxyName:         aws.String(dbProxyName),
		DBProxyEndpointName: aws.String(dbProxyEndpointName),
	}
	var dbProxyEndpoint *rds.DBProxyEndpoint

	err = conn.DescribeDBProxyEndpointsPagesWithContext(ctx, input, func(page *rds.DescribeDBProxyEndpointsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, endpoint := range page.DBProxyEndpoints {
			if aws.StringValue(endpoint.DBProxyEndpointName) == dbProxyEndpointName &&
				aws.StringValue(endpoint.DBProxyName) == dbProxyName {
				dbProxyEndpoint = endpoint
				return false
			}
		}

		return !lastPage
	})

	return dbProxyEndpoint, err
}

func FindDBClusterRoleByDBClusterIDAndRoleARN(ctx context.Context, conn *rds.RDS, dbClusterID, roleARN string) (*rds.DBClusterRole, error) {
	dbCluster, err := FindDBClusterByID(ctx, conn, dbClusterID)
	if err != nil {
		return nil, err
	}

	for _, associatedRole := range dbCluster.AssociatedRoles {
		if aws.StringValue(associatedRole.RoleArn) == roleARN {
			if status := aws.StringValue(associatedRole.Status); status == ClusterRoleStatusDeleted {
				return nil, &retry.NotFoundError{
					Message: status,
				}
			}

			return associatedRole, nil
		}
	}

	return nil, &retry.NotFoundError{}
}

func FindDBProxyByName(ctx context.Context, conn *rds.RDS, name string) (*rds.DBProxy, error) {
	input := &rds.DescribeDBProxiesInput{
		DBProxyName: aws.String(name),
	}

	output, err := conn.DescribeDBProxiesWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBProxyNotFoundFault) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.DBProxies) == 0 || output.DBProxies[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.DBProxies); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	dbProxy := output.DBProxies[0]

	// Eventual consistency check.
	if aws.StringValue(dbProxy.DBProxyName) != name {
		return nil, &retry.NotFoundError{
			LastRequest: input,
		}
	}

	return dbProxy, nil
}

func FindDBSubnetGroupByName(ctx context.Context, conn *rds.RDS, name string) (*rds.DBSubnetGroup, error) {
	input := &rds.DescribeDBSubnetGroupsInput{
		DBSubnetGroupName: aws.String(name),
	}

	output, err := conn.DescribeDBSubnetGroupsWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeDBSubnetGroupNotFoundFault) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.DBSubnetGroups) == 0 || output.DBSubnetGroups[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.DBSubnetGroups); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	dbSubnetGroup := output.DBSubnetGroups[0]

	// Eventual consistency check.
	if aws.StringValue(dbSubnetGroup.DBSubnetGroupName) != name {
		return nil, &retry.NotFoundError{
			LastRequest: input,
		}
	}

	return dbSubnetGroup, nil
}

func FindEventSubscriptionByID(ctx context.Context, conn *rds.RDS, id string) (*rds.EventSubscription, error) {
	input := &rds.DescribeEventSubscriptionsInput{
		SubscriptionName: aws.String(id),
	}

	output, err := conn.DescribeEventSubscriptionsWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeSubscriptionNotFoundFault) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.EventSubscriptionsList) == 0 || output.EventSubscriptionsList[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.EventSubscriptionsList); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	return output.EventSubscriptionsList[0], nil
}

func FindReservedDBInstanceByID(ctx context.Context, conn *rds.RDS, id string) (*rds.ReservedDBInstance, error) {
	input := &rds.DescribeReservedDBInstancesInput{
		ReservedDBInstanceId: aws.String(id),
	}

	output, err := conn.DescribeReservedDBInstancesWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, rds.ErrCodeReservedDBInstanceNotFoundFault) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || len(output.ReservedDBInstances) == 0 || output.ReservedDBInstances[0] == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	if count := len(output.ReservedDBInstances); count > 1 {
		return nil, tfresource.NewTooManyResultsError(count, input)
	}

	return output.ReservedDBInstances[0], nil
}
