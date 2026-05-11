#!/usr/bin/env bash
set -euo pipefail

# 1) Ενεργοποίηση venv
source .venv/bin/activate

# 2) Βάλε .env αν δεν υπάρχει
if [ ! -f .env ]; then
  echo 'POSTGRES_URL=postgresql+psycopg2://postgres:pass123@localhost:5432/eshopdb' > .env
fi

# 3) Τρέξε το pipeline (θα φορτώσει και σε Postgres αν το .env είναι σωστό)
python3 scripts/run_pipeline.py --input data/raw/sales_transactions.csv --outdir data/artifacts

# 4) Γρήγορο audit των logs
python3 scripts/quick_audit.py --logs data/artifacts/logs

# 5) Μερικά checks στη βάση
psql "postgresql://postgres:pass123@localhost:5432/eshopdb" -c "SELECT COUNT(*) FROM sales;"
psql "postgresql://postgres:pass123@localhost:5432/eshopdb" -c "SELECT MIN(transaction_date), MAX(transaction_date) FROM sales;"
psql "postgresql://postgres:pass123@localhost:5432/eshopdb" -c "SELECT * FROM sales LIMIT 5;"
