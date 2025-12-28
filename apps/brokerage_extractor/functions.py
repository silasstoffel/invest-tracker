import pdfplumber
import re
import requests
import os
import json

from datetime import datetime

from .symbol import search_symbol, search_symbol_type
from .types import Operation

def extract_operations_from_pdf(file_path: str) -> list[Operation]:
    operations: list[Operation] = []
    
    operation_total_value = None
    liquidation_value = None
    exchange_fees_value = None

    with pdfplumber.open(file_path) as pdf:
        for page in pdf.pages:
            text = page.extract_text()

            if text:
                # Getting operations
                operations = extract_operations(text)
                operation_total_value = extract_operation_total_value(text)
                liquidation_value = extract_liquidation_value(text)
                exchange_fees_value = extract_exchange_fees_value(text)
                costs = liquidation_value + exchange_fees_value
                split_operation_costs(operations, costs, operation_total_value)


    #is_unknown_item = is_unknown_value(operations)
    #if is_unknown_item == False:
    #    send_to_api(operations)
    return operations


def extract_numeric_value(text, pattern, group = 1):
    match = re.compile(pattern).search(text)
    if match:
        return float(match.group(group).replace(".", "").replace(",", "."))
    return None


def extract_liquidation_value(text):
    return extract_numeric_value(text, r"Taxa de liquidação\s+([\d,.]+)")


def extract_exchange_fees_value(text):
    return extract_numeric_value(text, r"Emolumentos\s+([\d,.]+)")


def extract_operation_total_value(text):
    return extract_numeric_value(text, r"Valor das operações\s+([\d,.]+)")


def extract_operation_type(text):
    b3 = "B3 RV LISTADO"
    bovespa = "1-BOVESPA"
    operation_type = (text.replace(b3, "").replace(bovespa, "")).strip()[:1]
    operation = { "C": "buy", "V": "sell" }

    return operation.get(operation_type, "UNKNOWN")


def extract_operations(text: str) -> list[Operation]:
    operations: list[Operation] = []
    txt = re.sub(r"\s+", " ", text)
    date_pattern = re.compile(r"\d{2}/\d{2}/\d{4}")
    pattern = re.compile(r"(B3 RV LISTADO|1-BOVESPA) (.*?) (\d+\s+[\d,.]+\s+[\d,.]+ [DC])")

    match_operation_date = date_pattern.search(text.replace("\n", " "))
    if match_operation_date:
        brazilian_date = match_operation_date.group(0)
        date = datetime.strptime(brazilian_date, "%d/%m/%Y")
        operation_date = date.strftime("%Y-%m-%d")

    matches = pattern.findall(txt)

    for match in matches:
        text_chunk = f"{match[0]} {match[1]}".strip()
        numbers_chunk = match[2].strip()
        raw_data = f"{text_chunk} {numbers_chunk}"
        values = numbers_chunk.split(" ")
        quantity = int(values[0])
        unit_price = float(values[1].replace(".", "").replace(",", "."))
        total_value = float(values[2].replace(".", "").replace(",", "."))

        operations.append({
            "raw_data": raw_data, 
            "operation_date": operation_date,
            "operation_type": extract_operation_type(text_chunk),
            "symbol": search_symbol(text_chunk, 'CLEAR'),
            "quantity": quantity,
            "unity_price": unit_price,
            "total_value": total_value,
            "costs": 0
        })

    return operations


def split_operation_costs(operations, costs, operation_total_value):
    for o in operations:
        if costs > 0:
            o["costs"] = round(costs * o["total_value"] / operation_total_value, 2)
    return operations


def is_unknown_value(operations):
    for o in operations:
        if o["operation_type"] == "UNKNOWN" or o["symbol"] == "UNKNOWN":
            return True
    
    return False


def send_to_api(operations: list[Operation]) -> None:
    if os.getenv("INTEGRATE_TO_API", "") != "true":
        print("Integration to API is disabled")
        return
    
    base_url = os.getenv("API_BASE_URL", "")
    url = f"{base_url}/investments/schedule"

    for operation in operations:
        symbol = operation['symbol']
        data = {
            "type": search_symbol_type(symbol), 
            "symbol": symbol,
            "quantity": operation['quantity'],
            "totalValue": operation['total_value'] + operation['costs'],
            "cost": operation['costs'],
            "operationType": operation['operation_type'],
            "operationDate": operation['operation_date'],
            "brokerage": "clear",
            "redemptionPolicyType" : "any_time",
            "note": "source: sinacor-brokerage-statements-extractor" 
        }
        
        response = requests.post(url, json=data, headers={"Content-Type": "application/json"})

        if response.status_code != 202:
            print("[%s] Failure request. URL: %s Status: %s - Content: %s" % (symbol, url, response.status_code, response.json()))
            print("[%s] Request body %s" % (symbol))
            print(json.dumps(data))
            continue
