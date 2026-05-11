import pandas as pd
from sqlalchemy import create_engine
import numpy as np

# -----------------------------
# Step 1: Load the CSV
# -----------------------------
df = pd.read_csv('sales_transactions.csv')
print("Initial data shape:", df.shape)

# -----------------------------
# Step 2: Normalize product names
# -----------------------------
before_unique_products = df['product_name'].nunique()

df['product_name'] = df['product_name'].astype(str).str.strip().str.lower().str.title()

# Fix common known typos or variations
df['product_name'] = df['product_name'].replace({
    'Usbc Cable': 'USB-C Cable',
    'Usb-C Cable': 'USB-C Cable',
    'Usb Cable': 'USB-C Cable',
    'Hdmi Cable': 'HDMI Cable',
    'Head Phones': 'Headphones'
})

after_unique_products = df['product_name'].nunique()
print(f"Product name normalization done. Unique products before: {before_unique_products}, after: {after_unique_products}")

# -----------------------------
# Step 3: Standardize dates
# -----------------------------
before_invalid_dates = df['transaction_date'].isna().sum()
df['transaction_date'] = pd.to_datetime(df['transaction_date'], errors='coerce')
after_invalid_dates = df['transaction_date'].isna().sum()
df = df[df['transaction_date'].notna()]
df['transaction_date'] = df['transaction_date'].dt.strftime('%Y-%m-%d')
print(f"Standardized dates. Invalid dates removed: {before_invalid_dates - after_invalid_dates}")

# -----------------------------
# Step 4: Detect & Impute Missing Values
# -----------------------------
missing_before = df[['quantity', 'price_per_unit', 'total_price']].isna().sum()

# Impute logically using relationships
df['quantity'] = df['quantity'].fillna(df['total_price'] / df['price_per_unit'])
df['price_per_unit'] = df['price_per_unit'].fillna(df['total_price'] / df['quantity'])
df['total_price'] = df['total_price'].fillna(df['price_per_unit'] * df['quantity'])

# Fill any leftovers with median
df['quantity'] = df['quantity'].fillna(df['quantity'].median())
df['price_per_unit'] = df['price_per_unit'].fillna(df['price_per_unit'].median())
df['total_price'] = df['total_price'].fillna(df['total_price'].median())

missing_after = df[['quantity', 'price_per_unit', 'total_price']].isna().sum()
print("Missing values imputed:")
print("Before:\n", missing_before)
print("After:\n", missing_after)

# -----------------------------
# Step 5: Remove Duplicates
# -----------------------------
before_duplicates = df.shape[0]
df = df.drop_duplicates()
after_duplicates = df.shape[0]
print(f"Duplicates removed: {before_duplicates - after_duplicates}")

# -----------------------------
# Step 6: Handle Outliers
# -----------------------------
before_outliers = df.shape[0]

# Rule-based filtering
df = df[(df['quantity'] > 0) & (df['quantity'] <= 1000)]
df = df[(df['price_per_unit'] > 0) & (df['total_price'] > 0)]

# Optional capping at 99th percentile
q_upper = df['quantity'].quantile(0.99)
tp_upper = df['total_price'].quantile(0.99)
df.loc[df['quantity'] > q_upper, 'quantity'] = q_upper
df.loc[df['total_price'] > tp_upper, 'total_price'] = tp_upper

after_outliers = df.shape[0]
print(f"Outlier handling done. Rows before: {before_outliers}, after: {after_outliers}")

# -----------------------------
# Step 6b: Ensure quantity is integer and format price columns
# -----------------------------
# Convert quantity to integer (rounding if needed)

df['quantity'] = df['quantity'].round().astype(int)

# Format price_per_unit and total_price to two decimals
for col in ['price_per_unit', 'total_price']:
    df[col] = df[col].apply(lambda x: f"{x:.2f}")

# -----------------------------
# Step 7: Save to PostgreSQL
# -----------------------------
engine = create_engine('postgresql://andy:andypass@localhost:5432/sales_db')

df.to_sql('sales_transactions', con=engine, schema='andy_schema', if_exists='replace', index=False)
print("Data saved to PostgreSQL. Final shape:", df.shape)

# -----------------------------
# Step 8: Save cleaned CSV
# -----------------------------
df.to_csv('sales_transactions_cleaned.csv', index=False)
print("Cleaned CSV saved as sales_transactions_cleaned.csv")
