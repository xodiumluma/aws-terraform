// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package kafka_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/kafka"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
)

func TestAccKafkaClusterDataSource_basic(t *testing.T) {
	ctx := acctest.Context(t)
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	dataSourceName := "data.aws_msk_cluster.test"
	resourceName := "aws_msk_cluster.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t); testAccPreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, kafka.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckClusterDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccClusterDataSourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "arn", resourceName, "arn"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bootstrap_brokers", resourceName, "bootstrap_brokers"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bootstrap_brokers_public_sasl_iam", resourceName, "bootstrap_brokers_public_sasl_iam"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bootstrap_brokers_public_sasl_scram", resourceName, "bootstrap_brokers_public_sasl_scram"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bootstrap_brokers_public_tls", resourceName, "bootstrap_brokers_public_tls"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bootstrap_brokers_sasl_iam", resourceName, "bootstrap_brokers_sasl_iam"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bootstrap_brokers_sasl_scram", resourceName, "bootstrap_brokers_sasl_scram"),
					resource.TestCheckResourceAttrPair(dataSourceName, "bootstrap_brokers_tls", resourceName, "bootstrap_brokers_tls"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cluster_name", resourceName, "cluster_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "cluster_uuid", resourceName, "cluster_uuid"),
					resource.TestCheckResourceAttrPair(dataSourceName, "kafka_version", resourceName, "kafka_version"),
					resource.TestCheckResourceAttrPair(dataSourceName, "number_of_broker_nodes", resourceName, "number_of_broker_nodes"),
					resource.TestCheckResourceAttrPair(dataSourceName, "tags.%", resourceName, "tags.%"),
					resource.TestCheckResourceAttrPair(dataSourceName, "zookeeper_connect_string", resourceName, "zookeeper_connect_string"),
					resource.TestCheckResourceAttrPair(dataSourceName, "zookeeper_connect_string_tls", resourceName, "zookeeper_connect_string_tls"),
				),
			},
		},
	})
}

func testAccClusterDataSourceConfig_basic(rName string) string {
	return acctest.ConfigCompose(testAccClusterConfig_base(rName), fmt.Sprintf(`
resource "aws_msk_cluster" "test" {
  cluster_name           = %[1]q
  kafka_version          = "2.8.1"
  number_of_broker_nodes = 3

  broker_node_group_info {
    client_subnets  = aws_subnet.test[*].id
    instance_type   = "kafka.t3.small"
    security_groups = [aws_security_group.test.id]

    storage_info {
      ebs_storage_info {
        volume_size = 10
      }
    }
  }

  tags = {
    Name = %[1]q
  }
}

data "aws_msk_cluster" "test" {
  cluster_name = aws_msk_cluster.test.cluster_name
}
`, rName))
}
