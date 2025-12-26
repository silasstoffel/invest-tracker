import logging

def hello(event: dict, context:dict) -> dict:
    logging.info("Hello from Brokerage Extractor!")
    return {
        "statusCode": 200,
        "body": "Hello from Brokerage Extractor!"
    }
