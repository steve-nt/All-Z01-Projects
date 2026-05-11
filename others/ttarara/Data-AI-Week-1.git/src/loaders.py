from __future__ import annotations
import os
from pathlib import Path
import pandas as pd
from sqlalchemy import create_engine
from dotenv import load_dotenv

def load_to_postgres(clean_csv: Path, table_name: str = "sales") -> bool:
    load_dotenv()
    postgres_url = os.getenv("POSTGRES_URL")
    if not postgres_url:
        print("POSTGRES_URL not set; skipping DB load.")
        return False
    df = pd.read_csv(clean_csv)
    engine = create_engine(postgres_url, future=True)
    with engine.begin() as conn:
        df.to_sql(table_name, con=conn, if_exists="replace", index=False)
    print(f"Loaded {len(df)} rows into Postgres table '{table_name}'.")
    return True
