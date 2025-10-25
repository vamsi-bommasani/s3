output "s3_bucket_id" {
  description = "The name of the S3 bucket"
  value       = aws_s3_bucket.s3.id
}

output "s3_bucket_arn" {
  description = "The ARN of the S3 bucket"
  value       = aws_s3_bucket.s3.arn
}

output "s3_bucket_region" {
  description = "The AWS region where the bucket is created"
  value       = aws_s3_bucket.s3.region
}

output "s3_bucket_policy" {
  description = "The JSON-encoded bucket policy"
  value       = aws_s3_bucket_policy.s3.policy
}
