# 📊 Sales Transactions Data Cleaning & PostgreSQL Loading

This project cleans, normalizes, and prepares a sales transactions dataset and saves it to PostgreSQL and a cleaned CSV.

---

## **Project Steps**

### 1. Load the CSV
- Read the `sales_transactions.csv` file using `pandas`.
- Print the initial shape of the dataset.

### 2. Normalize Product Names
- Convert all product names to title case and strip whitespace.
- Correct common typos and variations:
  - `"Usbc Cable"`, `"Usb-C Cable"`, `"Usb Cable"` → `"USB-C Cable"`
  - `"Hdmi Cable"` → `"HDMI Cable"`
  - `"Head Phones"` → `"Headphones"`
- Print the number of unique products before and after normalization.

### 3. Standardize Dates
- Convert `transaction_date` to `YYYY-MM-DD` format.
- Remove invalid or missing dates.
- Print the number of invalid dates removed.

### 4. Detect & Impute Missing Values
- Columns: `quantity`, `price_per_unit`, `total_price`.
- Impute missing values using logical relationships:
  - `quantity = total_price / price_per_unit`
  - `price_per_unit = total_price / quantity`
  - `total_price = price_per_unit * quantity`
- Fill any remaining missing values with the median.
- Print missing values before and after imputation.

### 5. Remove Duplicates
- Drop duplicate rows from the dataset.
- Print the number of duplicates removed.

### 6. Handle Outliers
- Filter:
  - `quantity > 0` and `quantity <= 1000`
  - `price_per_unit > 0` and `total_price > 0`
- Cap `quantity` and `total_price` at the 99th percentile.
- Print row count before and after outlier handling.

### 6b. Format Decimals
- Columns: `price_per_unit`, `total_price`.
- Ensure two decimal places for all numeric values:
  - `15.5` → `15.50`
  - `15.55` → `15.55`
- Does not round numbers; preserves original precision.

### 7. Save to PostgreSQL
- Connect to PostgreSQL using SQLAlchemy:
  ```python
  engine = create_engine('postgresql://andy:andypass@localhost:5432/sales_db')

## Running the Script

```bash
python3 script.py
```

## Redirect logs to a file

```bash
python3 script.py > script.log 2>&1
```


## 👩‍💻 Authors

For questions or issues, please contact us: [Georgia Marouli](https://discordapp.com/users/1277216244910522371) - [Andriana Stas](https://discordapp.com/users/780150798927134740)

> © 2025 Georgia Marouli and Andriana Stas for Zone01Athens Projects