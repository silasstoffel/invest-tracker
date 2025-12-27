from typing import TypedDict
from datetime import date

class Operation(TypedDict):
    raw_data: str
    operation_date: str
    operation_type: str
    symbol: str
    quantity: int
    unity_price: float
    total_value: float
    costs: float

    def to_dict(self) -> dict:
        return {
            "raw_data": self.raw_data,
            "operation_date": self.operation_date,
            "operation_type": self.operation_type,
            "symbol": self.symbol,
            "quantity": self.quantity,
            "unity_price": self.unity_price,
            "total_value": self.total_value,
            "costs": self.costs,
        }

