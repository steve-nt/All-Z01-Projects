from __future__ import annotations
from pathlib import Path
import pandas as pd
import numpy as np
from .utils import round_if_near_int, clean_product_name

CRITICAL_COLS = ["price_per_unit", "quantity", "total_price"]

def load_csv(path: Path) -> pd.DataFrame:
    return pd.read_csv(path)

def ensure_columns(df: pd.DataFrame) -> pd.DataFrame:
    for col in CRITICAL_COLS:
        if col not in df.columns:
            df[col] = np.nan
    if "product_name" not in df.columns:
        candidates = [c for c in df.columns if "product" in c.lower() and "name" in c.lower()]
        if candidates:
            df = df.rename(columns={candidates[0]: "product_name"})
        else:
            df["product_name"] = ""
    if "transaction_date" not in df.columns:
        candidates = [c for c in df.columns if "date" in c.lower()]
        if candidates:
            df = df.rename(columns={candidates[0]: "transaction_date"})
        else:
            df["transaction_date"] = pd.NaT
    return df

def coerce_numeric(df: pd.DataFrame) -> pd.DataFrame:
    for col in CRITICAL_COLS:
        df[col] = pd.to_numeric(df[col], errors="coerce")
    return df

def handle_missing(df: pd.DataFrame, logs_dir: Path) -> pd.DataFrame:
    logs_dir.mkdir(parents=True, exist_ok=True)
    mp = df["price_per_unit"].isna()
    mq = df["quantity"].isna()
    mt = df["total_price"].isna()
    missing_count = mp.astype(int) + mq.astype(int) + mt.astype(int)
    rows_two_plus = df[missing_count >= 2].copy()
    if not rows_two_plus.empty:
        rows_two_plus.to_csv(logs_dir / "dropped_missing.csv", index=False)
        df = df[missing_count < 2].copy()

    mp = df["price_per_unit"].isna()
    mq = df["quantity"].isna()
    mt = df["total_price"].isna()

    idx = df.index[mt & (~mp) & (~mq)]
    df.loc[idx, "total_price"] = df.loc[idx, "price_per_unit"] * df.loc[idx, "quantity"]

    idx = df.index[mp & (~mt) & (~mq) & (df["quantity"] != 0)]
    df.loc[idx, "price_per_unit"] = df.loc[idx, "total_price"] / df.loc[idx, "quantity"]

    idx = df.index[mq & (~mt) & (~mp) & (df["price_per_unit"] != 0)]
    df.loc[idx, "quantity"] = df.loc[idx, "total_price"] / df.loc[idx, "price_per_unit"]

    df["quantity"] = df["quantity"].apply(round_if_near_int)
    return df

def normalize_products(df: pd.DataFrame) -> pd.DataFrame:
    df["product_name"] = df["product_name"].apply(clean_product_name)
    return df

def standardize_dates(df: pd.DataFrame, logs_dir: Path) -> pd.DataFrame:
    raw = df["transaction_date"].copy()
    parsed = pd.to_datetime(raw, errors="coerce", infer_datetime_format=True, dayfirst=False)
    invalid_mask = parsed.isna() & raw.notna()
    invalid_rows = df[invalid_mask].copy()
    if not invalid_rows.empty:
        invalid_rows.to_csv(logs_dir / "invalid_dates.csv", index=False)
    df = df[~invalid_mask].copy()
    parsed = parsed[~invalid_mask]
    df["transaction_date"] = parsed.dt.strftime("%Y-%m-%d")
    return df

def drop_duplicates(df: pd.DataFrame) -> pd.DataFrame:
    return df.drop_duplicates()

def iqr_bounds(series: pd.Series):
    q1 = series.quantile(0.25)
    q3 = series.quantile(0.75)
    iqr = q3 - q1
    return q1 - 1.5 * iqr, q3 + 1.5 * iqr

def handle_outliers(df: pd.DataFrame, logs_dir: Path) -> pd.DataFrame:
    removed_log = []
    for col in ["quantity", "total_price"]:
        if col not in df.columns:
            continue
        s = pd.to_numeric(df[col], errors="coerce")
        lower, upper = iqr_bounds(s.dropna())
        mask_out = (s < lower) | (s > upper)
        if mask_out.any():
            part = df[mask_out].copy()
            if not part.empty:
                part["outlier_column"] = col
                removed_log.append(part)
            df = df[~mask_out].copy()
    if removed_log:
        outliers_removed = pd.concat(removed_log, axis=0)
        outliers_removed.to_csv(logs_dir / "outliers_removed.csv", index=False)
    return df

def pipeline(path: Path, outdir: Path) -> Path:
    outdir.mkdir(parents=True, exist_ok=True)
    logs_dir = outdir / "logs"
    logs_dir.mkdir(parents=True, exist_ok=True)

    df = load_csv(path)
    df = ensure_columns(df)
    df = coerce_numeric(df)
    df = handle_missing(df, logs_dir)
    df = normalize_products(df)
    df = standardize_dates(df, logs_dir)
    df = drop_duplicates(df)
    df = handle_outliers(df, logs_dir)

    clean_path = outdir / "clean_sales.csv"
    df.to_csv(clean_path, index=False)
    return clean_path
