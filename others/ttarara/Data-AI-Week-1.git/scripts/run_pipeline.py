#!/usr/bin/env python
from __future__ import annotations
import os, sys
import argparse
from pathlib import Path


sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))

from src.cleaning import pipeline
from src.loaders import load_to_postgres

def main():
    ap = argparse.ArgumentParser(description="Clean and (optionally) load sales data to Postgres.")
    ap.add_argument("--input", required=True, type=Path, help="Path to raw CSV (e.g., data/raw/sales_transactions.csv)")
    ap.add_argument("--outdir", required=True, type=Path, help="Artifacts output directory (e.g., data/artifacts)")
    ap.add_argument("--no-db", action="store_true", help="Skip loading to Postgres even if POSTGRES_URL is set")
    args = ap.parse_args()

    clean_csv = pipeline(args.input, args.outdir)
    print(f"Clean CSV saved to: {clean_csv}")

    if not args.no_db:
        load_to_postgres(clean_csv)

if __name__ == "__main__":
    main()
