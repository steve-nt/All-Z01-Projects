#!/usr/bin/env python
from __future__ import annotations
import argparse
from pathlib import Path
import pandas as pd

def safe_len(path: Path) -> int:
    return int(pd.read_csv(path).shape[0]) if path.exists() else 0

def main():
    ap = argparse.ArgumentParser(description="Show quick audit summary from logs directory.")
    ap.add_argument("--logs", required=True, type=Path, help="Path to logs directory")
    args = ap.parse_args()

    dropped_missing = args.logs / "dropped_missing.csv"
    invalid_dates   = args.logs / "invalid_dates.csv"
    outliers_removed= args.logs / "outliers_removed.csv"

    summary = {
        "dropped_missing": safe_len(dropped_missing),
        "invalid_dates": safe_len(invalid_dates),
        "outliers_removed": safe_len(outliers_removed),
    }
    print("Audit summary:", summary)

if __name__ == "__main__":
    main()
