# 🧹 Data & AI – The Little eShop of Horrors

A data-cleaning and ingestion pipeline built with **Python, Pandas, and PostgreSQL**.  
The goal is to rescue *Little eShop’s* chaotic sales data from CSV nightmares — missing values, inconsistent formats, duplicates, outliers, and messy product names.

---

## 📊 Project Overview

On your first day as a junior data specialist, your manager hands you a `sales_transactions.csv` file full of:
- 🕳️ Missing prices and quantities  
- 🔁 Duplicated rows  
- 🧩 Mixed product name formats (“USB-C cable”, “ usb-c Cable”, “USBC Cable”, etc.)  
- 🗓️ Inconsistent or invalid dates  
- 💥 Unrealistic outliers (e.g., 1,500 headphones in one order)  

Your mission is to **clean, normalize, and load** this dataset into a **PostgreSQL database**.

---

## ⚙️ Features

✅ Detect & clean missing values  
✅ Normalize product names (lowercase, trimmed, standardized)  
✅ Standardize all dates to `YYYY-MM-DD` format  
✅ Remove duplicates  
✅ Detect and handle outliers in `quantity` and `total_price`  
✅ Export cleaned dataset to CSV  
✅ Optional automatic load to PostgreSQL (`sales` table)  
✅ Detailed logs of invalid / removed entries  
✅ Fully automatable with one command (`run.sh`)

---

## 🧱 Project Structure

```
eshop_data_pipeline/
│
├── data/
│   ├── raw/                  # original dataset (sales_transactions.csv)
│   └── artifacts/            # cleaned data & logs
│       ├── clean_sales.csv
│       └── logs/
│           ├── invalid_dates.csv
│           ├── outliers_removed.csv
│           └── clean_audit.log
│
├── src/                      # main cleaning logic
│   ├── cleaning.py
│   ├── loaders.py
│   └── utils.py
│
├── scripts/                  # execution entry points
│   ├── run_pipeline.py
│   └── quick_audit.py
│
├── notebooks/                # optional Jupyter exploration
│   └── week1_sales_ingestion_notebook.ipynb
│
├── .env.example              # PostgreSQL connection template
├── run.sh                    # auto-run script
├── requirements.txt
├── README.md
└── Makefile (optional)
```

---

## 🚀 Setup Instructions

### 1️⃣ Clone the repo
```bash
git clone https://platform.zone01.gr/git/ttarara/Data-and-AI-Week1.git
cd eshop_data_pipeline
```

### 2️⃣ Create a virtual environment
```bash
python3 -m venv .venv
source .venv/bin/activate
```

### 3️⃣ Install dependencies
```bash
pip install --upgrade pip
pip install -r requirements.txt
```

### 4️⃣ (Optional) Set up PostgreSQL with Docker
```bash
docker run -d --name eshop-pg   -e POSTGRES_USER=postgres   -e POSTGRES_PASSWORD=pass123   -e POSTGRES_DB=eshopdb   -p 5432:5432 postgres:16
```

### 5️⃣ Configure the environment
```bash
cp .env.example .env
# Edit .env to include:
POSTGRES_URL=postgresql+psycopg2://postgres:pass123@localhost:5432/eshopdb
```

---

## 🧩 Running the Pipeline

### Option A – Full automated mode
```bash
chmod +x run.sh
./run.sh
```

### Option B – Manual mode
```bash
# Clean and load data
python3 scripts/run_pipeline.py --input data/raw/sales_transactions.csv --outdir data/artifacts

# Optional: skip DB load
python3 scripts/run_pipeline.py --input data/raw/sales_transactions.csv --outdir data/artifacts --no-db

# Audit summary
python3 scripts/quick_audit.py --logs data/artifacts/logs
```

---

## 🗃️ Database Verification

Check the loaded data directly in PostgreSQL:

```bash
psql "postgresql://postgres:pass123@localhost:5432/eshopdb" -c "SELECT COUNT(*) FROM sales;"
psql "postgresql://postgres:pass123@localhost:5432/eshopdb" -c "SELECT * FROM sales LIMIT 5;"
```

Example Output:

```
 count
-------
   995
(1 row)
```

---

## 📂 Output Files

After running the pipeline, you’ll find:

| File | Description |
|------|--------------|
| `data/artifacts/clean_sales.csv` | Final cleaned dataset |
| `data/artifacts/logs/invalid_dates.csv` | Entries with invalid dates |
| `data/artifacts/logs/outliers_removed.csv` | Detected and removed outliers |
| `data/artifacts/logs/clean_audit.log` | Summary of cleaning steps |

---

## 🧪 Example Audit Summary

```bash
Audit summary: {'dropped_missing': 0, 'invalid_dates': 3, 'outliers_removed': 5}
```

---

## 🧰 Technologies Used

| Tool | Purpose |
|------|----------|
| 🐍 Python 3.12 | Data processing |
| 🧮 Pandas | Data cleaning and transformation |
| 🗃️ PostgreSQL | Data storage |
| ⚡ SQLAlchemy | Database ORM |
| 🧰 psycopg2 | Postgres driver |
| 🧩 dotenv | Environment management |
| 🐳 Docker | Containerized Postgres instance |

---

## 📘 Author

**Theoharoula Tarara**  
👩‍💻 Junior Data Specialist – Zone01 / Data & AI Piscine  
📧 Contact: [ttarara@zone01.gr](mailto:ttarara@zone01.gr)

---

## 🏁 Final Notes

- To rerun the cleaning anytime:  
  ```bash
  ./run.sh
  ```
- To stop your local database:
  ```bash
  docker stop eshop-pg
  ```
- To remove everything (clean slate):
  ```bash
  docker rm -f eshop-pg
  rm -rf data/artifacts
  ```

> 💡 *All requirements of the Piscine Data & AI Week 1 project are implemented: data validation, normalization, logging, and database ingestion. Ready for audit submission.*