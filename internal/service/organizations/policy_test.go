package organizations_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/service/organizations"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tforganizations "github.com/hashicorp/terraform-provider-aws/internal/service/organizations"
)

func testAccPolicy_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var policy organizations.Policy
	content1 := `{"Version": "2012-10-17", "Statement": { "Effect": "Allow", "Action": "*", "Resource": "*"}}`
	content2 := `{"Version": "2012-10-17", "Statement": { "Effect": "Allow", "Action": "s3:*", "Resource": "*"}}`
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_required(rName, content1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					acctest.MatchResourceAttrGlobalARN(resourceName, "arn", "organizations", regexp.MustCompile("policy/o-.+/service_control_policy/p-.+$")),
					resource.TestCheckResourceAttr(resourceName, "content", content1),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "type", organizations.PolicyTypeServiceControlPolicy),
				),
			},
			{
				Config: testAccPolicyConfig_required(rName, content2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "content", content2),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_destroy"},
			},
		},
	})
}

// Reference: https://github.com/hashicorp/terraform-provider-aws/issues/5073
func testAccPolicy_concurrent(t *testing.T) {
	ctx := acctest.Context(t)
	var policy1, policy2, policy3, policy4, policy5 organizations.Policy
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName1 := "aws_organizations_policy.test1"
	resourceName2 := "aws_organizations_policy.test2"
	resourceName3 := "aws_organizations_policy.test3"
	resourceName4 := "aws_organizations_policy.test4"
	resourceName5 := "aws_organizations_policy.test5"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_concurrent(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName1, &policy1),
					testAccCheckPolicyExists(ctx, resourceName2, &policy2),
					testAccCheckPolicyExists(ctx, resourceName3, &policy3),
					testAccCheckPolicyExists(ctx, resourceName4, &policy4),
					testAccCheckPolicyExists(ctx, resourceName5, &policy5),
				),
			},
		},
	})
}

func testAccPolicy_description(t *testing.T) {
	ctx := acctest.Context(t)
	var policy organizations.Policy
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_description(rName, "description1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "description", "description1"),
				),
			},
			{
				Config: testAccPolicyConfig_description(rName, "description2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_destroy"},
			},
		},
	})
}

func testAccPolicy_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var p1, p2, p3, p4 organizations.Policy
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_tagA(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &p1),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.TerraformProviderAwsTest", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.Alpha", "1"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_destroy"},
			},
			{
				Config: testAccPolicyConfig_tagB(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &p2),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.TerraformProviderAwsTest", "true"),
					resource.TestCheckResourceAttr(resourceName, "tags.Beta", "1"),
				),
			},
			{
				Config: testAccPolicyConfig_tagC(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &p3),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.TerraformProviderAwsTest", "true"),
				),
			},
			{
				Config: testAccPolicyConfig_noTag(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &p4),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
		},
	})
}

func testAccPolicy_skipDestroy(t *testing.T) {
	ctx := acctest.Context(t)
	var policy organizations.Policy
	content := `{"Version": "2012-10-17", "Statement": { "Effect": "Allow", "Action": "*", "Resource": "*"}}`
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyNoDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_skipDestroy(rName, content),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					acctest.MatchResourceAttrGlobalARN(resourceName, "arn", "organizations", regexp.MustCompile("policy/o-.+/service_control_policy/p-.+$")),
					resource.TestCheckResourceAttr(resourceName, "content", content),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "skip_destroy", "true"),
					resource.TestCheckResourceAttr(resourceName, "type", organizations.PolicyTypeServiceControlPolicy),
				),
			},
		},
	})
}

func testAccPolicy_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var p organizations.Policy
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_description(rName, ""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &p),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tforganizations.ResourcePolicy(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccPolicy_type_AI_OPT_OUT(t *testing.T) {
	ctx := acctest.Context(t)
	var policy organizations.Policy
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"
	// Reference: https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_ai-opt-out_syntax.html
	AiOptOutPolicyContent := `{ "services": { "rekognition": { "opt_out_policy": { "@@assign": "optOut" } }, "lex": { "opt_out_policy": { "@@assign": "optIn" } } } }`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_type(rName, AiOptOutPolicyContent, organizations.PolicyTypeAiservicesOptOutPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "type", organizations.PolicyTypeAiservicesOptOutPolicy),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_destroy"},
			},
		},
	})
}

func testAccPolicy_type_Backup(t *testing.T) {
	ctx := acctest.Context(t)
	var policy organizations.Policy
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"
	// Reference: https://docs.aws.amazon.com/organizations/latest/userguide/orgs_manage_policies_backup_syntax.html
	backupPolicyContent := fmt.Sprintf(`{
   "plans":{
      "PII_Backup_Plan":{
         "regions":{
            "@@assign":[
               "%[1]s",
               "%[2]s"
            ]
         },
         "rules":{
            "Hourly":{
               "schedule_expression":{
                  "@@assign":"cron(0 5/1 ? * * *)"
               },
               "start_backup_window_minutes":{
                  "@@assign":"480"
               },
               "complete_backup_window_minutes":{
                  "@@assign":"10080"
               },
               "lifecycle":{
                  "move_to_cold_storage_after_days":{
                     "@@assign":"180"
                  },
                  "delete_after_days":{
                     "@@assign":"270"
                  }
               },
               "target_backup_vault_name":{
                  "@@assign":"FortKnox"
               },
               "copy_actions":{
                  "arn:%[3]s:backup:%[1]s:$account:backup-vault:secondary_vault":{
                     "target_backup_vault_arn":{
                        "@@assign":"arn:%[3]s:backup:%[1]s:$account:backup-vault:secondary_vault"
                     },
                     "lifecycle":{
                        "delete_after_days":{
                           "@@assign":"100"
                        },
                        "move_to_cold_storage_after_days":{
                           "@@assign":"10"
                        }
                     }
                  }
               }
            }
         },
         "selections":{
            "tags":{
               "datatype":{
                  "iam_role_arn":{
                     "@@assign":"arn:%[3]s:iam::$account:role/MyIamRole"
                  },
                  "tag_key":{
                     "@@assign":"dataType"
                  },
                  "tag_value":{
                     "@@assign":[
                        "PII",
                        "RED"
                     ]
                  }
               }
            }
         }
      }
   }
}`, acctest.AlternateRegion(), acctest.Region(), acctest.Partition())

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_type(rName, backupPolicyContent, organizations.PolicyTypeBackupPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "type", organizations.PolicyTypeBackupPolicy),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_destroy"},
			},
		},
	})
}

func testAccPolicy_type_SCP(t *testing.T) {
	ctx := acctest.Context(t)
	var policy organizations.Policy
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"
	serviceControlPolicyContent := `{"Version": "2012-10-17", "Statement": { "Effect": "Allow", "Action": "*", "Resource": "*"}}`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_type(rName, serviceControlPolicyContent, organizations.PolicyTypeServiceControlPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "type", organizations.PolicyTypeServiceControlPolicy),
				),
			},
			{
				Config: testAccPolicyConfig_required(rName, serviceControlPolicyContent),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "type", organizations.PolicyTypeServiceControlPolicy),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_destroy"},
			},
		},
	})
}

func testAccPolicy_type_Tag(t *testing.T) {
	ctx := acctest.Context(t)
	var policy organizations.Policy
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_organizations_policy.test"
	tagPolicyContent := `{ "tags": { "Product": { "tag_key": { "@@assign": "Product" }, "enforced_for": { "@@assign": [ "ec2:instance" ] } } } }`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_type(rName, tagPolicyContent, organizations.PolicyTypeTagPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPolicyExists(ctx, resourceName, &policy),
					resource.TestCheckResourceAttr(resourceName, "type", organizations.PolicyTypeTagPolicy),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"skip_destroy"},
			},
		},
	})
}

func testAccPolicy_importManagedPolicy(t *testing.T) {
	ctx := acctest.Context(t)
	resourceName := "aws_organizations_policy.test"

	resourceID := "p-FullAWSAccess"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); acctest.PreCheckOrganizationsAccount(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, organizations.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPolicyDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyConfig_managedSetup,
			},
			{
				Config:        testAccPolicyConfig_managed,
				ResourceName:  resourceName,
				ImportStateId: resourceID,
				ImportState:   true,
				ExpectError:   regexp.MustCompile(regexp.QuoteMeta(fmt.Sprintf("AWS-managed Organizations policy (%s) cannot be imported.", resourceID))),
			},
		},
	})
}

func testAccCheckPolicyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).OrganizationsConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_organizations_policy" {
				continue
			}

			input := &organizations.DescribePolicyInput{
				PolicyId: &rs.Primary.ID,
			}

			resp, err := conn.DescribePolicyWithContext(ctx, input)

			if tfawserr.ErrCodeEquals(err, organizations.ErrCodeAWSOrganizationsNotInUseException) {
				continue
			}

			if tfawserr.ErrCodeEquals(err, organizations.ErrCodePolicyNotFoundException) {
				continue
			}

			if err != nil {
				return err
			}

			if resp != nil && resp.Policy != nil {
				return fmt.Errorf("Policy %q still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

// testAccCheckPolicyNoDestroy is a variant of the CheckDestroy function to be used when
// skip_destroy is true and the policy should still exist after destroy completes
func testAccCheckPolicyNoDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).OrganizationsConn()

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_organizations_policy" {
				continue
			}

			input := &organizations.DescribePolicyInput{
				PolicyId: &rs.Primary.ID,
			}

			_, err := conn.DescribePolicyWithContext(ctx, input)
			if tfawserr.ErrCodeEquals(err, organizations.ErrCodeAWSOrganizationsNotInUseException) {
				// The organization was destroyed, so we can safely assume the policy
				// skipped during destruction was as well
				continue
			}
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccCheckPolicyExists(ctx context.Context, resourceName string, policy *organizations.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).OrganizationsConn()
		input := &organizations.DescribePolicyInput{
			PolicyId: &rs.Primary.ID,
		}

		resp, err := conn.DescribePolicyWithContext(ctx, input)

		if err != nil {
			return err
		}

		if resp == nil || resp.Policy == nil {
			return fmt.Errorf("Policy %q does not exist", rs.Primary.ID)
		}

		*policy = *resp.Policy

		return nil
	}
}

func testAccPolicyConfig_description(rName, description string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "*",
    "Resource": "*"
  }
}
EOF

  description = "%s"
  name        = "%s"

  depends_on = [aws_organizations_organization.test]
}
`, description, rName)
}

func testAccPolicyConfig_tagA(rName string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "*",
    "Resource": "*"
  }
}
EOF

  name = "%s"

  depends_on = [aws_organizations_organization.test]

  tags = {
    TerraformProviderAwsTest = true
    Alpha                    = 1
  }
}
`, rName)
}

func testAccPolicyConfig_tagB(rName string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "*",
    "Resource": "*"
  }
}
EOF

  name = "%s"

  depends_on = [aws_organizations_organization.test]

  tags = {
    TerraformProviderAwsTest = true
    Beta                     = 1
  }
}
`, rName)
}

func testAccPolicyConfig_tagC(rName string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "*",
    "Resource": "*"
  }
}
EOF

  name = "%s"

  depends_on = [aws_organizations_organization.test]

  tags = {
    TerraformProviderAwsTest = true
  }
}
`, rName)
}

func testAccPolicyConfig_noTag(rName string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Allow",
    "Action": "*",
    "Resource": "*"
  }
}
EOF

  name = "%s"

  depends_on = [aws_organizations_organization.test]
}
`, rName)
}

func testAccPolicyConfig_required(rName, content string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test" {
  content = %s
  name    = "%s"

  depends_on = [aws_organizations_organization.test]
}
`, strconv.Quote(content), rName)
}

func testAccPolicyConfig_concurrent(rName string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test1" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Deny",
    "Action": "cloudtrail:StopLogging",
    "Resource": "*"
  }
}
EOF

  name = "%[1]s1"

  depends_on = [aws_organizations_organization.test]
}

resource "aws_organizations_policy" "test2" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Deny",
    "Action": "ec2:DeleteFlowLogs",
    "Resource": "*"
  }
}
EOF

  name = "%[1]s2"

  depends_on = [aws_organizations_organization.test]
}

resource "aws_organizations_policy" "test3" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Deny",
    "Action": "logs:DeleteLogGroup",
    "Resource": "*"
  }
}
EOF

  name = "%[1]s3"

  depends_on = [aws_organizations_organization.test]
}

resource "aws_organizations_policy" "test4" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Deny",
    "Action": "config:DeleteConfigRule",
    "Resource": "*"
  }
}
EOF

  name = "%[1]s4"

  depends_on = [aws_organizations_organization.test]
}

resource "aws_organizations_policy" "test5" {
  content = <<EOF
{
  "Version": "2012-10-17",
  "Statement": {
    "Effect": "Deny",
    "Action": "iam:DeleteRolePermissionsBoundary",
    "Resource": "*"
  }
}
EOF

  name = "%[1]s5"

  depends_on = [aws_organizations_organization.test]
}
`, rName)
}

func testAccPolicyConfig_type(rName, content, policyType string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test" {
  content = %s
  name    = "%s"
  type    = "%s"

  depends_on = [aws_organizations_organization.test]
}
`, strconv.Quote(content), rName, policyType)
}

func testAccPolicyConfig_skipDestroy(rName, content string) string {
	return fmt.Sprintf(`
resource "aws_organizations_organization" "test" {}

resource "aws_organizations_policy" "test" {
  content = %s
  name    = "%s"

  skip_destroy = true

  depends_on = [aws_organizations_organization.test]
}
`, strconv.Quote(content), rName)
}

const testAccPolicyConfig_managedSetup = `
resource "aws_organizations_organization" "test" {
  enabled_policy_types = ["SERVICE_CONTROL_POLICY"]
}
`

const testAccPolicyConfig_managed = `
resource "aws_organizations_organization" "test" {
  enabled_policy_types = ["SERVICE_CONTROL_POLICY"]
}

resource "aws_organizations_policy" "test" {
  name = "FullAWSAccess"
}
`
