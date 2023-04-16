# EKS
# * IAM Role - allow EKS to manage other AWS
# * EC2 Security Group - permit networking with EKS cluster
# * EKS Cluster

resource "aws_iam_role" "demo-cluster" {
  name = "terraform-eks-demo-cluster"
}