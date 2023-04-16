package elasticbeanstalk_test

import (
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/service/elasticbeanstalk"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccElasticBeanstalkSolutionStackDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	dataSourceName := "data.aws_elastic_beanstalk_solution_stack.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, elasticbeanstalk.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccSolutionStackDataSourceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(dataSourceName, "name", regexp.MustCompile("^64bit Amazon Linux (.*) running Python (.*)$")),
				),
			},
		},
	})
}

const testAccSolutionStackDataSourceConfig_basic = `
data "aws_elastic_beanstalk_solution_stack" "test" {
  most_recent = true

  # e.g. "64bit Amazon Linux 2018.03 v2.10.14 running Python 3.6"
  name_regex = "^64bit Amazon Linux (.*) running Python (.*)$"
}
`
