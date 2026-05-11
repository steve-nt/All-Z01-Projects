from __future__ import annotations
import re
import pandas as pd

def round_if_near_int(x, tol: float = 1e-6):
    if pd.isna(x):
        return x
    if abs(x - round(x)) < tol:
        return int(round(x))
    return x

def clean_product_name(s: str) -> str:
    if pd.isna(s):
        return ""
    s = str(s).lower().strip()
    s = re.sub(r"\s+", " ", s)
    if re.search(r"\busb[-\s]?c\b", s) and re.search(r"\bcable\b", s):
        return "usb-c cable"
    return s
