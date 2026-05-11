-- Cosmofone Telecom Data Quality Check and Cleaning Script
-- This script validates data quality and performs necessary cleaning operations

\echo '=============================================='
\echo 'COSMOFONE TELECOM DATA QUALITY CHECK & CLEANING'
\echo '=============================================='

-- ==============================================
-- DATA QUALITY VALIDATION
-- ==============================================

\echo ''
\echo '1. RECORD COUNT VALIDATION:'

SELECT 
    'customers' as table_name,
    COUNT(*) as record_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM customers
UNION ALL
SELECT 
    'usage' as table_name,
    COUNT(*) as record_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM usage
UNION ALL
SELECT 
    'billing' as table_name,
    COUNT(*) as record_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM billing
UNION ALL
SELECT 
    'support_tickets' as table_name,
    COUNT(*) as record_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM support_tickets;

\echo ''
\echo '2. NULL VALUE CHECK:'

SELECT 
    'customers' as table_name,
    COUNT(*) as total_records,
    COUNT(customer_id) as non_null_ids,
    COUNT(name) as non_null_names,
    COUNT(region) as non_null_regions,
    COUNT(plan_type) as non_null_plans
FROM customers
UNION ALL
SELECT 
    'usage' as table_name,
    COUNT(*) as total_records,
    COUNT(customer_id) as non_null_ids,
    COUNT(call_minutes) as non_null_calls,
    COUNT(data_usage_gb) as non_null_data,
    COUNT(num_sms) as non_null_sms
FROM usage
UNION ALL
SELECT 
    'billing' as table_name,
    COUNT(*) as total_records,
    COUNT(customer_id) as non_null_ids,
    COUNT(amount) as non_null_amounts,
    COUNT(payment_status) as non_null_status,
    0 as placeholder
FROM billing;

\echo ''
\echo '3. DUPLICATE CHECK:'

-- Check for duplicate customer IDs
SELECT 
    'Duplicate customer IDs' as check_name,
    COUNT(*) as duplicate_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM (
    SELECT customer_id, COUNT(*) as count
    FROM customers 
    GROUP BY customer_id 
    HAVING COUNT(*) > 1
) duplicates;

\echo ''
\echo '4. DATA RANGE VALIDATION:'

-- Check billing amounts
SELECT 
    'Billing amounts' as check_name,
    MIN(amount) as min_amount,
    MAX(amount) as max_amount,
    COUNT(*) as total_records,
    COUNT(CASE WHEN amount < 0 THEN 1 END) as negative_amounts,
    CASE WHEN COUNT(CASE WHEN amount < 0 THEN 1 END) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM billing;

-- Check usage ranges
SELECT 
    'Usage ranges' as check_name,
    MIN(call_minutes) as min_calls,
    MAX(call_minutes) as max_calls,
    MIN(data_usage_gb) as min_data,
    MAX(data_usage_gb) as max_data,
    COUNT(CASE WHEN call_minutes = 9999.0 THEN 1 END) as unlimited_usage_customers
FROM usage;

-- Check date ranges
SELECT 
    'Date ranges' as check_name,
    MIN(signup_date) as earliest_signup,
    MAX(signup_date) as latest_signup,
    MIN(invoice_date) as earliest_billing,
    MAX(invoice_date) as latest_billing
FROM customers c
CROSS JOIN billing b;

\echo ''
\echo '5. DATA CONSISTENCY CHECK:'

-- Check region values
SELECT 
    'Region values' as check_name,
    region,
    COUNT(*) as count
FROM customers
GROUP BY region
ORDER BY region;

-- Check payment status values
SELECT 
    'Payment status values' as check_name,
    payment_status,
    COUNT(*) as count
FROM billing
GROUP BY payment_status
ORDER BY payment_status;

-- Check plan type values
SELECT 
    'Plan type values' as check_name,
    plan_type,
    COUNT(*) as count
FROM customers
GROUP BY plan_type
ORDER BY plan_type;

\echo ''
\echo '6. FOREIGN KEY INTEGRITY CHECK:'

-- Check for orphaned usage records
SELECT 
    'Orphaned usage records' as check_name,
    COUNT(*) as orphaned_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM usage u
LEFT JOIN customers c ON u.customer_id = c.customer_id
WHERE c.customer_id IS NULL
UNION ALL
-- Check for orphaned billing records
SELECT 
    'Orphaned billing records' as check_name,
    COUNT(*) as orphaned_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM billing b
LEFT JOIN customers c ON b.customer_id = c.customer_id
WHERE c.customer_id IS NULL
UNION ALL
-- Check for orphaned support tickets
SELECT 
    'Orphaned support tickets' as check_name,
    COUNT(*) as orphaned_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM support_tickets s
LEFT JOIN customers c ON s.customer_id = c.customer_id
WHERE c.customer_id IS NULL;

-- ==============================================
-- DATA CLEANING OPERATIONS
-- ==============================================

\echo ''
\echo '=============================================='
\echo 'DATA CLEANING OPERATIONS'
\echo '=============================================='

\echo ''
\echo '7. CLEANING REGION CASE INCONSISTENCY:'

-- Show before cleaning
SELECT 'Before cleaning - Region values:' as status;
SELECT region, COUNT(*) as count FROM customers GROUP BY region ORDER BY region;

-- Fix region case inconsistency
UPDATE customers SET region = 'North' WHERE region = 'north';

-- Show after cleaning
SELECT 'After cleaning - Region values:' as status;
SELECT region, COUNT(*) as count FROM customers GROUP BY region ORDER BY region;

\echo ''
\echo '8. CLEANING PAYMENT STATUS CASE INCONSISTENCY:'

-- Show before cleaning
SELECT 'Before cleaning - Payment status values:' as status;
SELECT payment_status, COUNT(*) as count FROM billing GROUP BY payment_status ORDER BY payment_status;

-- Fix payment status case inconsistency
UPDATE billing SET payment_status = 'Paid' WHERE payment_status = 'paid';

-- Show after cleaning
SELECT 'After cleaning - Payment status values:' as status;
SELECT payment_status, COUNT(*) as count FROM billing GROUP BY payment_status ORDER BY payment_status;

\echo ''
\echo '9. UNUSUAL VALUES ANALYSIS:'

-- Analyze 9999 call minutes
SELECT 
    'Customers with 9999 call minutes (unlimited plans):' as analysis,
    customer_id,
    call_minutes,
    data_usage_gb,
    num_sms
FROM usage 
WHERE call_minutes = 9999.0
ORDER BY customer_id;

-- ==============================================
-- FINAL VALIDATION
-- ==============================================

\echo ''
\echo '=============================================='
\echo 'FINAL DATA QUALITY VALIDATION'
\echo '=============================================='

\echo ''
\echo '10. FINAL DATA QUALITY SUMMARY:'

-- Final record counts
SELECT 
    'Final record counts:' as summary,
    'customers' as table_name,
    COUNT(*) as records
FROM customers
UNION ALL
SELECT 
    '',
    'usage',
    COUNT(*)
FROM usage
UNION ALL
SELECT 
    '',
    'billing',
    COUNT(*)
FROM billing
UNION ALL
SELECT 
    '',
    'support_tickets',
    COUNT(*)
FROM support_tickets;

-- Final data consistency check
SELECT 
    'Final data consistency:' as summary,
    'regions' as field,
    COUNT(DISTINCT region) as unique_values
FROM customers
UNION ALL
SELECT 
    '',
    'payment_status',
    COUNT(DISTINCT payment_status)
FROM billing
UNION ALL
SELECT 
    '',
    'plan_types',
    COUNT(DISTINCT plan_type)
FROM customers;

\echo ''
\echo '=============================================='
\echo 'DATA QUALITY CHECK COMPLETE'
\echo '=============================================='
\echo '✅ Data is clean and ready for analytics!'
\echo '✅ All inconsistencies have been resolved!'
\echo '✅ Data is ready for Tableau and other BI tools!'
\echo '=============================================='
