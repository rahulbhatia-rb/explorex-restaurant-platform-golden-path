variable "bucket_name" { type = string }

resource "aws_s3_bucket" "web" {
  bucket = var.bucket_name
}

resource "aws_s3_bucket_public_access_block" "web" {
  bucket                  = aws_s3_bucket.web.id
  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

# Production implementation would add CloudFront, Origin Access Control,
# ACM, Route 53 aliases, and immutable hashed asset cache policies.
