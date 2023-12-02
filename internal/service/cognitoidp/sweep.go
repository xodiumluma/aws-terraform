// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package cognitoidp

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_cognito_user_pool_domain", &resource.Sweeper{
		Name: "aws_cognito_user_pool_domain",
		F:    sweepUserPoolDomains,
	})

	resource.AddTestSweepers("aws_cognito_user_pool", &resource.Sweeper{
		Name: "aws_cognito_user_pool",
		F:    sweepUserPools,
		Dependencies: []string{
			"aws_cognito_user_pool_domain",
		},
	})
}

func sweepUserPoolDomains(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("Error getting client: %s", err)
	}
	conn := client.CognitoIDPConn(ctx)

	input := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(50),
	}

	err = conn.ListUserPoolsPagesWithContext(ctx, input, func(resp *cognitoidentityprovider.ListUserPoolsOutput, lastPage bool) bool {
		if len(resp.UserPools) == 0 {
			log.Print("[DEBUG] No Cognito user pools (i.e. domains) to sweep")
			return false
		}

		for _, u := range resp.UserPools {
			output, err := conn.DescribeUserPoolWithContext(ctx, &cognitoidentityprovider.DescribeUserPoolInput{
				UserPoolId: u.Id,
			})
			if err != nil {
				log.Printf("[ERROR] Failed describing Cognito user pool (%s): %s", aws.StringValue(u.Name), err)
				continue
			}
			if output.UserPool != nil && output.UserPool.Domain != nil {
				domain := aws.StringValue(output.UserPool.Domain)

				log.Printf("[INFO] Deleting Cognito user pool domain: %s", domain)
				_, err := conn.DeleteUserPoolDomainWithContext(ctx, &cognitoidentityprovider.DeleteUserPoolDomainInput{
					Domain:     output.UserPool.Domain,
					UserPoolId: u.Id,
				})
				if err != nil {
					log.Printf("[ERROR] Failed deleting Cognito user pool domain (%s): %s", domain, err)
				}
			}
		}
		return !lastPage
	})

	if err != nil {
		if awsv1.SkipSweepError(err) {
			log.Printf("[WARN] Skipping Cognito User Pool Domain sweep for %s: %s", region, err)
			return nil
		}
		return fmt.Errorf("Error retrieving Cognito User Pools: %s", err)
	}

	return nil
}

func sweepUserPools(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("error getting client: %w", err)
	}
	conn := client.CognitoIDPConn(ctx)

	input := &cognitoidentityprovider.ListUserPoolsInput{
		MaxResults: aws.Int64(50),
	}

	err = conn.ListUserPoolsPagesWithContext(ctx, input, func(resp *cognitoidentityprovider.ListUserPoolsOutput, lastPage bool) bool {
		if len(resp.UserPools) == 0 {
			log.Print("[DEBUG] No Cognito User Pools to sweep")
			return false
		}

		for _, userPool := range resp.UserPools {
			name := aws.StringValue(userPool.Name)

			log.Printf("[INFO] Deleting Cognito User Pool: %s", name)
			_, err := conn.DeleteUserPoolWithContext(ctx, &cognitoidentityprovider.DeleteUserPoolInput{
				UserPoolId: userPool.Id,
			})
			if err != nil {
				log.Printf("[ERROR] Failed deleting Cognito User Pool (%s): %s", name, err)
			}
		}
		return !lastPage
	})

	if err != nil {
		if awsv1.SkipSweepError(err) {
			log.Printf("[WARN] Skipping Cognito User Pool sweep for %s: %s", region, err)
			return nil
		}
		return fmt.Errorf("Error retrieving Cognito User Pools: %w", err)
	}

	return nil
}
