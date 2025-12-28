#!/bin/bash

# How use: ./scripts/statement-upload.sh "./XPINC_NOTA_NEGOCIACAO_B3_5_8_2025.pdf"

# Upload PDF file to S3 bucket using AWS CLI

set -euo pipefail

# Configuration
PDF_FILE="${1:?Error: PDF file path required as argument 1}"
S3_BUCKET="invest-track-statements-dev"
AWS_REGION="us-east-1"

# Validate file exists
if [[ ! -f "$PDF_FILE" ]]; then
    echo "Error: File '$PDF_FILE' not found"
    exit 1
fi

# Validate file is PDF
if [[ "$PDF_FILE" != *.pdf ]]; then
    echo "Error: File must be a PDF"
    exit 1
fi

# Extract filename
FILENAME=$(basename "$PDF_FILE")

# Upload to S3
echo "Uploading $FILENAME to s3://$S3_BUCKET"
aws s3 cp "$PDF_FILE" "s3://$S3_BUCKET/$FILENAME" \
    --region "$AWS_REGION" \
    --storage-class STANDARD

echo "Upload completed successfully"
