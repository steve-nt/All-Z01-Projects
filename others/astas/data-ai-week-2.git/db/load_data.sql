-- Data Loading Script for Cosmofone Telecom Data
-- This script loads all CSV data into PostgreSQL tables

-- Create temp table first, then insert unique records
CREATE TEMP TABLE customers_temp (
    customer_id VARCHAR,
    name VARCHAR,
    region VARCHAR,
    age NUMERIC,
    plan_type VARCHAR,
    signup_date DATE
);

-- Load customers data into temp table
\copy customers_temp (customer_id, name, region, age, plan_type, signup_date) FROM 'customers.csv' WITH (FORMAT csv, HEADER true, NULL '');

-- Insert unique customers only
INSERT INTO customers (customer_id, name, region, age, plan_type, signup_date)
SELECT DISTINCT customer_id, name, region, age, plan_type, signup_date
FROM customers_temp
ON CONFLICT (customer_id) DO NOTHING;

-- Load usage data
\copy usage (customer_id, usage_month, call_minutes, data_usage_gb, num_sms) FROM 'usage.csv' WITH (FORMAT csv, HEADER true, NULL '');

-- Load billing data
\copy billing (customer_id, invoice_date, amount, payment_status) FROM 'billing.csv' WITH (FORMAT csv, HEADER true, NULL '');

-- Load support tickets data
\copy support_tickets (customer_id, ticket_date, issue_type, resolution_time_hrs) FROM 'support_tickets.csv' WITH (FORMAT csv, HEADER true, NULL '');

-- Update any NULL values in amount column to 0
UPDATE billing SET amount = 0 WHERE amount IS NULL;

-- Update any NULL values in resolution_time_hrs to 0
UPDATE support_tickets SET resolution_time_hrs = 0 WHERE resolution_time_hrs IS NULL;

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_usage_customer_id ON usage(customer_id);
CREATE INDEX IF NOT EXISTS idx_usage_month ON usage(usage_month);
CREATE INDEX IF NOT EXISTS idx_billing_customer_id ON billing(customer_id);
CREATE INDEX IF NOT EXISTS idx_billing_invoice_date ON billing(invoice_date);
CREATE INDEX IF NOT EXISTS idx_support_tickets_customer_id ON support_tickets(customer_id);
CREATE INDEX IF NOT EXISTS idx_support_tickets_ticket_date ON support_tickets(ticket_date);

-- Verify data loading
SELECT 'Customers loaded:' as table_name, COUNT(*) as record_count FROM customers
UNION ALL
SELECT 'Usage records loaded:', COUNT(*) FROM usage
UNION ALL
SELECT 'Billing records loaded:', COUNT(*) FROM billing
UNION ALL
SELECT 'Support tickets loaded:', COUNT(*) FROM support_tickets;

-- Display sample data
\echo 'Sample customer data:'
SELECT * FROM customers LIMIT 5;

\echo 'Sample usage data:'
SELECT * FROM usage LIMIT 5;

\echo 'Sample billing data:'
SELECT * FROM billing LIMIT 5;

\echo 'Sample support tickets data:'
SELECT * FROM support_tickets LIMIT 5;
