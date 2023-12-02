// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ses

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep"
	"github.com/hashicorp/terraform-provider-aws/internal/sweep/awsv1"
)

func RegisterSweepers() {
	resource.AddTestSweepers("aws_ses_configuration_set", &resource.Sweeper{
		Name: "aws_ses_configuration_set",
		F:    sweepConfigurationSets,
	})

	resource.AddTestSweepers("aws_ses_domain_identity", &resource.Sweeper{
		Name: "aws_ses_domain_identity",
		F:    func(region string) error { return sweepIdentities(region, ses.IdentityTypeDomain) },
	})

	resource.AddTestSweepers("aws_ses_email_identity", &resource.Sweeper{
		Name: "aws_ses_email_identity",
		F:    func(region string) error { return sweepIdentities(region, ses.IdentityTypeEmailAddress) },
	})

	resource.AddTestSweepers("aws_ses_receipt_rule_set", &resource.Sweeper{
		Name: "aws_ses_receipt_rule_set",
		F:    sweepReceiptRuleSets,
	})
}

func sweepConfigurationSets(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("getting client: %w", err)
	}
	conn := client.SESConn(ctx)
	input := &ses.ListConfigurationSetsInput{}
	var sweeperErrs *multierror.Error

	for {
		output, err := conn.ListConfigurationSetsWithContext(ctx, input)
		if awsv1.SkipSweepError(err) {
			log.Printf("[WARN] Skipping SES Configuration Sets sweep for %s: %s", region, err)
			return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
		}
		if err != nil {
			sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("retrieving SES Configuration Sets: %w", err))
			return sweeperErrs
		}

		for _, configurationSet := range output.ConfigurationSets {
			name := aws.StringValue(configurationSet.Name)

			log.Printf("[INFO] Deleting SES Configuration Set: %s", name)
			_, err := conn.DeleteConfigurationSetWithContext(ctx, &ses.DeleteConfigurationSetInput{
				ConfigurationSetName: aws.String(name),
			})
			if tfawserr.ErrCodeEquals(err, ses.ErrCodeConfigurationSetDoesNotExistException) {
				continue
			}
			if err != nil {
				sweeperErr := fmt.Errorf("deleting SES Configuration Set (%s): %w", name, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}
		input.NextToken = output.NextToken
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepIdentities(region, identityType string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("getting client: %w", err)
	}
	conn := client.SESConn(ctx)
	input := &ses.ListIdentitiesInput{
		IdentityType: aws.String(identityType),
	}
	var sweeperErrs *multierror.Error

	err = conn.ListIdentitiesPagesWithContext(ctx, input, func(page *ses.ListIdentitiesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, identity := range page.Identities {
			identity := aws.StringValue(identity)

			log.Printf("[INFO] Deleting SES Identity: %s", identity)
			_, err = conn.DeleteIdentityWithContext(ctx, &ses.DeleteIdentityInput{
				Identity: aws.String(identity),
			})
			if err != nil {
				sweeperErr := fmt.Errorf("deleting SES Identity (%s): %w", identity, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}
		}

		return !lastPage
	})
	if awsv1.SkipSweepError(err) {
		log.Printf("[WARN] Skipping SES Identities sweep for %s: %s", region, err)
		return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
	}
	if err != nil {
		sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("retrieving SES Identities: %w", err))
	}

	return sweeperErrs.ErrorOrNil()
}

func sweepReceiptRuleSets(region string) error {
	ctx := sweep.Context(region)
	client, err := sweep.SharedRegionalSweepClient(ctx, region)
	if err != nil {
		return fmt.Errorf("getting client: %w", err)
	}
	conn := client.SESConn(ctx)

	// You cannot delete the receipt rule set that is currently active.
	// Setting the name of the active receipt rule set to null disables all email receiving.
	log.Printf("[INFO] Disabling any currently active SES Receipt Rule Set")
	_, err = conn.SetActiveReceiptRuleSetWithContext(ctx, &ses.SetActiveReceiptRuleSetInput{})
	// In some regions, this will return "InvalidAction" with no message
	if awsv1.SkipSweepError(err) || tfawserr.ErrCodeEquals(err, "InvalidAction") {
		log.Printf("[WARN] Skipping SES Receipt Rule Sets sweep for %s: %s", region, err)
		return nil
	}
	if err != nil {
		return fmt.Errorf("disabling any currently active SES Receipt Rule Set: %w", err)
	}

	input := &ses.ListReceiptRuleSetsInput{}
	var sweeperErrs *multierror.Error

	for {
		output, err := conn.ListReceiptRuleSetsWithContext(ctx, input)
		if awsv1.SkipSweepError(err) {
			log.Printf("[WARN] Skipping SES Receipt Rule Sets sweep for %s: %s", region, err)
			return sweeperErrs.ErrorOrNil() // In case we have completed some pages, but had errors
		}
		if err != nil {
			sweeperErrs = multierror.Append(sweeperErrs, fmt.Errorf("retrieving SES Receipt Rule Sets: %w", err))
			return sweeperErrs
		}

		for _, ruleSet := range output.RuleSets {
			name := aws.StringValue(ruleSet.Name)

			log.Printf("[INFO] Deleting SES Receipt Rule Set: %s", name)
			_, err := conn.DeleteReceiptRuleSetWithContext(ctx, &ses.DeleteReceiptRuleSetInput{
				RuleSetName: aws.String(name),
			})
			if tfawserr.ErrCodeEquals(err, ses.ErrCodeRuleSetDoesNotExistException) {
				continue
			}
			if err != nil {
				sweeperErr := fmt.Errorf("deleting SES Receipt Rule Set (%s): %w", name, err)
				log.Printf("[ERROR] %s", sweeperErr)
				sweeperErrs = multierror.Append(sweeperErrs, sweeperErr)
				continue
			}
		}

		if aws.StringValue(output.NextToken) == "" {
			break
		}
		input.NextToken = output.NextToken
	}

	return sweeperErrs.ErrorOrNil()
}
