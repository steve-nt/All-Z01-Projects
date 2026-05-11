# Cosmofone Telecom — Data Analytics Project

## Project Overview

This repository contains a complete data analytics solution for Cosmofone, a telecommunications provider. It includes schema definitions, data-loading scripts, data-quality checks, analytical views, and example queries suitable for business stakeholders and technical reviewers.

Key outcomes:
- All datasets ingested into PostgreSQL and validated
- Star schema implemented (fact tables + customer dimension)
- Data-quality automation and remediation scripts included
- Analytical views for customer summary, monthly KPIs, and churn risk indicators

## Audit Summary (All requirements met)

Data Ingestion
- All source files loaded into PostgreSQL:
  - `customers.csv` (100+ records)
  - `usage.csv` (350+ records)
  - `billing.csv` (350+ records)
  - `support_tickets.csv` (150+ records)

Schema Design
- Star schema with central fact tables (`usage`, `billing`, `support_tickets`) and a `customers` dimension
- Proper normalization and key relationships defined

## Project Structure

The following reflects the repository layout in this workspace:

```
data-ai-week-2/
├── README.md
├── setup_postgresql.sh             # PostgreSQL setup script
├── cleaned data/
│   ├── billing_clean.csv
│   ├── customers_clean.csv
│   ├── support_tickets_clean.csv
│   └── usage_clean.csv
├── db/
│   ├── analytical_queries.sql
│   ├── data_quality_check.sql
│   ├── load_data.sql
│   ├── telecom_schema_and_views.sql
│   └── validate_setup.sql
├── original_project_files/
│   ├── billing.csv
│   ├── customers.csv
│   ├── support_tickets.csv
│   └── usage.csv
├── views/
│   ├── churn_risk_indicators_high_tickets_missed_payments.csv
│   ├── customer_summary.csv
│   └── monthly_kpis.csv
└── tableau_dashboards/
    ├── Customer Insights.png
    ├── Executive Overview.png
    └── Plan Performance - Average Usage and Billing by Plan Type.png
```

## Database Schema Design

Star schema used in this project:

Fact tables:
- `usage` — customer usage metrics (call minutes, data MB, SMS)
- `billing` — billing and payment records
- `support_tickets` — customer support interactions

Dimension table:
- `customers` — customer demographics and plan details

Key relationships:
- `customers.customer_id` → `usage.customer_id` (1:N)
- `customers.customer_id` → `billing.customer_id` (1:N)
- `customers.customer_id` → `support_tickets.customer_id` (1:N)

## Setup

Prerequisites
- Ubuntu / Linux (WSL supported)
- PostgreSQL 12+ (the included setup script will install/configure PostgreSQL)

Quick setup (recommended):
1. Make the setup script executable:

```bash
chmod +x setup_postgresql.sh
```

2. Run the setup script:

```bash
./setup_postgresql.sh
```

Manual setup (alternative):

1. Install PostgreSQL:

```bash
sudo apt update
sudo apt install -y postgresql postgresql-contrib
```

2. Create the database and user:

```bash
sudo -u postgres psql -c "CREATE DATABASE cosmofone;"
sudo -u postgres psql -c "CREATE USER gina WITH PASSWORD 'ginapass';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE cosmofone TO gina;"
```

3. Create schema and load data:

```bash
psql -h localhost -U gina -d cosmofone -f db/telecom_schema_and_views.sql
psql -h localhost -U gina -d cosmofone -f db/load_data.sql
```

## Data Quality & Cleaning

The project includes automated checks and remediation in `db/data_quality_check.sql` and manual commands used during analysis.

Summary of checks performed:
- Record counts and non-null critical fields
- Duplicate detection
- Reasonableness of numeric ranges
- Categorical standardization (regions, payment status)

Examples of fixes applied

```sql
-- Standardize region values
UPDATE customers SET region = 'North' WHERE region = 'north';

-- Standardize payment status
UPDATE billing SET payment_status = 'Paid' WHERE payment_status = 'paid';
```

Unusual values
- A small number of customers have unusually high usage values (e.g., 9999 minutes). These were reviewed and retained as indicators of unlimited plans.

Automated check (run the SQL script):

```bash
psql -h localhost -U gina -d cosmofone -f db/data_quality_check.sql
```

Manual validation queries used during development:

```bash
psql -h localhost -U gina -d cosmofone -c "SELECT 'customers' as table, COUNT(*) as total, COUNT(customer_id) as non_null_ids, COUNT(name) as non_null_names FROM customers;"

psql -h localhost -U gina -d cosmofone -c "SELECT customer_id, COUNT(*) as count FROM customers GROUP BY customer_id HAVING COUNT(*) > 1;"

psql -h localhost -U gina -d cosmofone -c "SELECT MIN(amount) as min_amount, MAX(amount) as max_amount, COUNT(*) as total_bills FROM billing;"

psql -h localhost -U gina -d cosmofone -c "SELECT DISTINCT region FROM customers ORDER BY region;"
psql -h localhost -U gina -d cosmofone -c "SELECT DISTINCT payment_status FROM billing;"
```

Final data quality status
- Cleaned and validated — ready for analytics

## Analytical Views

Views included (in `db/telecom_schema_and_views.sql`):

1. `customer_summary` — combines customer details with aggregated usage and billing statistics
2. `monthly_kpis` — monthly aggregates for active users, average usage, and revenue
3. `churn_risk_indicators` — flags customers at risk of churn based on support activity and unpaid bills

## Analytical Queries

A set of example analytical queries is provided in `db/analytical_queries.sql`. These include revenue trends, ARPU calculations, and churn analyses used in dashboard prototypes.

## Connection Details (for development)

- Host: `localhost`
- Port: `5432`
- Database: `cosmofone`
- Username: `gina`
- Password: `ginapass`

## Dashboards (Tableau)

Dashboard mockups and configuration notes are in `tableau_dashboards/`. Key dashboards:
- Executive Overview — high-level KPIs and revenue trends
- Customer Insights — segmentation, churn analysis, support performance
- Plan Performance — usage and revenue breakdown by plan

## Validation Checklist

- [x] Data ingestion completed
- [x] Star schema implemented
- [x] Data quality checks and cleaning automated
- [x] Analytical views created and validated
- [x] Sample queries for BI tools provided

## Support

If you have questions, inspect the queries in `db/analytical_queries.sql` or contact the data team! We would love to hear your feedback!:)

 ## 👩‍💻 Authors

For questions or issues, please contact us: [Georgia Marouli](https://discordapp.com/users/1277216244910522371) - [Andriana Stas](https://discordapp.com/users/780150798927134740)

> © 2025 Georgia Marouli and Andriana Stas for Zone01Athens Projects
            