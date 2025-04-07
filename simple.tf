# Create IAM Policy
resource "aws_iam_policy" "s3_access_policy" {
  name        = "S3ListGetPolicy"
  description = "Policy allowing listing and getting objects from all S3 buckets"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = ["s3:ListBucket"]
        Resource = ["arn:aws:s3:::*"]
      },
      {
        Effect   = "Allow"
        Action   = ["s3:GetObject"]
        Resource = ["arn:aws:s3:::*/*"]
      }
    ]
  })
}

# Create IAM Role
resource "aws_iam_role" "my_test_role" {
  name = "MyTerraformTestRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          AWS = "*"
        }
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# Attach IAM Policy to IAM Role
resource "aws_iam_role_policy_attachment" "s3_policy_attachment" {
  role       = aws_iam_role.my_test_role.id
  policy_arn = aws_iam_policy.s3_access_policy.arn
}
