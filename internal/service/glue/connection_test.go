// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package glue_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/glue"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfglue "github.com/hashicorp/terraform-provider-aws/internal/service/glue"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func TestAccGlueConnection_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	jdbcConnectionUrl := fmt.Sprintf("jdbc:mysql://%s/testdatabase", acctest.RandomDomainName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_required(rName, jdbcConnectionUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					acctest.CheckResourceAttrRegionalARN(resourceName, "arn", "glue", fmt.Sprintf("connection/%s", rName)),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.JDBC_CONNECTION_URL", jdbcConnectionUrl),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.PASSWORD", "testpassword"),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.USERNAME", "testusername"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "physical_connection_requirements.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlueConnection_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	jdbcConnectionUrl := fmt.Sprintf("jdbc:mysql://%s/testdatabase", acctest.RandomDomainName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_tags1(rName, jdbcConnectionUrl, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConnectionConfig_tags2(rName, jdbcConnectionUrl, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccConnectionConfig_tags1(rName, jdbcConnectionUrl, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccGlueConnection_mongoDB(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	connectionUrl := fmt.Sprintf("mongodb://%s:27017/testdatabase", acctest.RandomDomainName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_mongoDB(rName, connectionUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.CONNECTION_URL", connectionUrl),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.USERNAME", "testusername"),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.PASSWORD", "testpassword"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "MONGODB"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "physical_connection_requirements.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlueConnection_kafka(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	bootstrapServers := fmt.Sprintf("%s:9094,%s:9094", acctest.RandomDomainName(), acctest.RandomDomainName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_kafka(rName, bootstrapServers),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.KAFKA_BOOTSTRAP_SERVERS", bootstrapServers),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "KAFKA"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "physical_connection_requirements.#", "0"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlueConnection_network(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_network(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.%", "0"),
					resource.TestCheckResourceAttr(resourceName, "connection_type", "NETWORK"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "physical_connection_requirements.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "physical_connection_requirements.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "physical_connection_requirements.0.security_group_id_list.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "physical_connection_requirements.0.subnet_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlueConnection_description(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	jdbcConnectionUrl := fmt.Sprintf("jdbc:mysql://%s/testdatabase", acctest.RandomDomainName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_description(rName, jdbcConnectionUrl, "First Description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "description", "First Description"),
				),
			},
			{
				Config: testAccConnectionConfig_description(rName, jdbcConnectionUrl, "Second Description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "description", "Second Description"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlueConnection_matchCriteria(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	jdbcConnectionUrl := fmt.Sprintf("jdbc:mysql://%s/testdatabase", acctest.RandomDomainName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_matchCriteriaFirst(rName, jdbcConnectionUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.#", "4"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.0", "criteria1"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.1", "criteria2"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.2", "criteria3"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.3", "criteria4"),
				),
			},
			{
				Config: testAccConnectionConfig_matchCriteriaSecond(rName, jdbcConnectionUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.0", "criteria1"),
				),
			},
			{
				Config: testAccConnectionConfig_matchCriteriaThird(rName, jdbcConnectionUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.#", "3"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.0", "criteria2"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.1", "criteria3"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.2", "criteria4"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlueConnection_physicalConnectionRequirements(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_physicalRequirements(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					resource.TestCheckResourceAttr(resourceName, "connection_properties.%", "3"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_properties.JDBC_CONNECTION_URL"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_properties.PASSWORD"),
					resource.TestCheckResourceAttrSet(resourceName, "connection_properties.USERNAME"),
					resource.TestCheckResourceAttr(resourceName, "match_criteria.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "physical_connection_requirements.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "physical_connection_requirements.0.availability_zone"),
					resource.TestCheckResourceAttr(resourceName, "physical_connection_requirements.0.security_group_id_list.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "physical_connection_requirements.0.subnet_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGlueConnection_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var connection glue.Connection

	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_glue_connection.test"

	jdbcConnectionUrl := fmt.Sprintf("jdbc:mysql://%s/testdatabase", acctest.RandomDomainName())

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, glue.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckConnectionDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccConnectionConfig_required(rName, jdbcConnectionUrl),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckConnectionExists(ctx, resourceName, &connection),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfglue.ResourceConnection(), resourceName),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfglue.ResourceConnection(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckConnectionExists(ctx context.Context, resourceName string, connection *glue.Connection) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Glue Connection ID is set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).GlueConn(ctx)
		catalogID, connectionName, err := tfglue.DecodeConnectionID(rs.Primary.ID)
		if err != nil {
			return err
		}

		output, err := tfglue.FindConnectionByName(ctx, conn, connectionName, catalogID)

		if err != nil {
			return err
		}

		*connection = *output

		return nil
	}
}

func testAccCheckConnectionDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_glue_connection" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).GlueConn(ctx)
			catalogID, connectionName, err := tfglue.DecodeConnectionID(rs.Primary.ID)
			if err != nil {
				return err
			}

			_, err = tfglue.FindConnectionByName(ctx, conn, connectionName, catalogID)

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return fmt.Errorf("Glue Connection %s still exists", rs.Primary.ID)
		}

		return nil
	}
}

func testAccConnectionConfig_description(rName, jdbcConnectionUrl, description string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name        = %[1]q
  description = %[2]q

  connection_properties = {
    JDBC_CONNECTION_URL = %[3]q
    PASSWORD            = "testpassword"
    USERNAME            = "testusername"
  }

}
`, rName, description, jdbcConnectionUrl)
}

func testAccConnectionConfig_matchCriteriaFirst(rName, jdbcConnectionUrl string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name = %[1]q

  match_criteria = ["criteria1", "criteria2", "criteria3", "criteria4"]

  connection_properties = {
    JDBC_CONNECTION_URL = %[2]q
    PASSWORD            = "testpassword"
    USERNAME            = "testusername"
  }
}
`, rName, jdbcConnectionUrl)
}

func testAccConnectionConfig_matchCriteriaSecond(rName, jdbcConnectionUrl string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name = %[1]q

  match_criteria = ["criteria1"]

  connection_properties = {
    JDBC_CONNECTION_URL = %[2]q
    PASSWORD            = "testpassword"
    USERNAME            = "testusername"
  }
}
`, rName, jdbcConnectionUrl)
}

func testAccConnectionConfig_matchCriteriaThird(rName, jdbcConnectionUrl string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name = "%s"

  match_criteria = ["criteria2", "criteria3", "criteria4"]

  connection_properties = {
    JDBC_CONNECTION_URL = %[2]q
    PASSWORD            = "testpassword"
    USERNAME            = "testusername"
  }
}
`, rName, jdbcConnectionUrl)
}

func testAccConnectionConfig_physicalRequirements(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-glue-connection-base"
  }
}

resource "aws_security_group" "test" {
  name   = "%[1]s"
  vpc_id = aws_vpc.test.id

  ingress {
    from_port = 1
    protocol  = "tcp"
    self      = true
    to_port   = 65535
  }
}

resource "aws_subnet" "test" {
  count = 2

  availability_zone = data.aws_availability_zones.available.names[count.index]
  cidr_block        = "10.0.${count.index}.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "terraform-testacc-glue-connection-base"
  }
}

resource "aws_db_subnet_group" "test" {
  name       = "%[1]s"
  subnet_ids = aws_subnet.test[*].id
}

data "aws_rds_engine_version" "default" {
  engine = "aurora-mysql"
}

data "aws_rds_orderable_db_instance" "test" {
  engine                     = data.aws_rds_engine_version.default.engine
  engine_version             = data.aws_rds_engine_version.default.version
  preferred_instance_classes = ["db.t3.small", "db.t3.medium", "db.t3.large"]
}

resource "aws_rds_cluster" "test" {
  cluster_identifier     = "%[1]s"
  database_name          = "gluedatabase"
  db_subnet_group_name   = aws_db_subnet_group.test.name
  engine                 = data.aws_rds_orderable_db_instance.test.engine
  engine_version         = data.aws_rds_orderable_db_instance.test.engine_version
  master_password        = "gluepassword"
  master_username        = "glueusername"
  skip_final_snapshot    = true
  vpc_security_group_ids = [aws_security_group.test.id]
}

resource "aws_rds_cluster_instance" "test" {
  identifier         = "%[1]s"
  cluster_identifier = aws_rds_cluster.test.id
  engine             = data.aws_rds_orderable_db_instance.test.engine
  engine_version     = data.aws_rds_orderable_db_instance.test.engine_version
  instance_class     = data.aws_rds_orderable_db_instance.test.instance_class
}

resource "aws_glue_connection" "test" {
  connection_properties = {
    JDBC_CONNECTION_URL = "jdbc:mysql://${aws_rds_cluster.test.endpoint}/${aws_rds_cluster.test.database_name}"
    PASSWORD            = aws_rds_cluster.test.master_password
    USERNAME            = aws_rds_cluster.test.master_username
  }

  name = "%[1]s"

  physical_connection_requirements {
    availability_zone      = aws_subnet.test[0].availability_zone
    security_group_id_list = [aws_security_group.test.id]
    subnet_id              = aws_subnet.test[0].id
  }
}
`, rName)
}

func testAccConnectionConfig_required(rName, jdbcConnectionUrl string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name = %[1]q

  connection_properties = {
    JDBC_CONNECTION_URL = %[2]q
    PASSWORD            = "testpassword"
    USERNAME            = "testusername"
  }
}
`, rName, jdbcConnectionUrl)
}

func testAccConnectionConfig_tags1(rName, jdbcConnectionUrl, tagKey1, tagValue1 string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name = %[1]q

  connection_properties = {
    JDBC_CONNECTION_URL = %[2]q
    PASSWORD            = "testpassword"
    USERNAME            = "testusername"
  }

  tags = {
    %[3]q = %[4]q
  }
}
`, rName, jdbcConnectionUrl, tagKey1, tagValue1)
}

func testAccConnectionConfig_tags2(rName, jdbcConnectionUrl, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name = %[1]q

  connection_properties = {
    JDBC_CONNECTION_URL = %[2]q
    PASSWORD            = "testpassword"
    USERNAME            = "testusername"
  }

  tags = {
    %[3]q = %[4]q
    %[5]q = %[6]q
  }
}
`, rName, jdbcConnectionUrl, tagKey1, tagValue1, tagKey2, tagValue2)
}

func testAccConnectionConfig_mongoDB(rName, connectionUrl string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name = %[1]q

  connection_type = "MONGODB"
  connection_properties = {
    CONNECTION_URL = %[2]q
    PASSWORD       = "testpassword"
    USERNAME       = "testusername"
  }
}
`, rName, connectionUrl)
}

func testAccConnectionConfig_kafka(rName, bootstrapServers string) string {
	return fmt.Sprintf(`
resource "aws_glue_connection" "test" {
  name = %[1]q

  connection_type = "KAFKA"
  connection_properties = {
    KAFKA_BOOTSTRAP_SERVERS = %[2]q
  }
}
`, rName, bootstrapServers)
}

func testAccConnectionConfig_network(rName string) string {
	return fmt.Sprintf(`
data "aws_availability_zones" "available" {
  state = "available"

  filter {
    name   = "opt-in-status"
    values = ["opt-in-not-required"]
  }
}

resource "aws_vpc" "test" {
  cidr_block = "10.0.0.0/16"

  tags = {
    Name = "terraform-testacc-glue-connection-network"
  }
}

resource "aws_subnet" "test" {
  availability_zone = data.aws_availability_zones.available.names[0]
  cidr_block        = "10.0.0.0/24"
  vpc_id            = aws_vpc.test.id

  tags = {
    Name = "terraform-testacc-glue-connection-network"
  }
}

resource "aws_security_group" "test" {
  name   = "%[1]s"
  vpc_id = aws_vpc.test.id

  ingress {
    protocol  = "tcp"
    self      = true
    from_port = 1
    to_port   = 65535
  }
}

resource "aws_glue_connection" "test" {
  connection_type = "NETWORK"
  name            = "%[1]s"

  physical_connection_requirements {
    availability_zone      = aws_subnet.test.availability_zone
    security_group_id_list = [aws_security_group.test.id]
    subnet_id              = aws_subnet.test.id
  }
}
`, rName)
}
