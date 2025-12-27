import logging
import boto3
import urllib.parse
import tempfile
import os

from .functions import extract_operations_from_pdf, is_unknown_value, send_to_api
from .types import Operation

s3 = boto3.client("s3")

def handler(event: dict, context:dict) -> None:
    logging.info("Hello from Brokerage Extractor!")

    record = event["Records"][0]
    bucket = record["s3"]["bucket"]["name"]
    key = urllib.parse.unquote_plus(
        record["s3"]["object"]["key"]
    )

    logging.info(f"Processing file from bucket: {bucket}, key: {key}")
    operations:list[Operation] = []

    with tempfile.NamedTemporaryFile() as tmp_file:
        s3.download_file(bucket, key, tmp_file.name)
        operations = extract_operations_from_pdf(tmp_file.name)

        logging.info(f"Extracted text {bucket}, key: {key}")

    is_unknown_item = is_unknown_value(operations)
    if is_unknown_item == False:
        logging.info(f"Sending operations to API, count: {len(operations)}")
        send_to_api(operations)
        logging.info(f"Operations sent to API successfully.")
    else:
        #TODO: notify telegram bot about unknown items
        logging.warning(f"Unknown items found in operations, skipping API send.")    
