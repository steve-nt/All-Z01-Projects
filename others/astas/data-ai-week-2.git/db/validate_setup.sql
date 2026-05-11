-- Cosmofone Telecom Analytics - Validation Script
-- This script validates that all components are working correctly

\echo '=============================================='
\echo 'COSMOFONE TELECOM ANALYTICS - VALIDATION REPORT'
\echo '=============================================='

-- Check database connection
\echo '1. DATABASE CONNECTION:'
SELECT 'Connected to database: ' || current_database() as status;

-- Verify all tables exist and have data
\echo ''
\echo '2. TABLE VALIDATION:'

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

-- Verify views exist and return data
\echo ''
\echo '3. VIEW VALIDATION:'

SELECT 
    'customer_summary' as view_name,
    COUNT(*) as record_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM customer_summary
UNION ALL
SELECT 
    'monthly_kpis' as view_name,
    COUNT(*) as record_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM monthly_kpis
UNION ALL
SELECT 
    'churn_risk_indicators' as view_name,
    COUNT(*) as record_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM churn_risk_indicators;

-- Check data quality
\echo ''
\echo '4. DATA QUALITY CHECKS:'

-- Check for NULL customer IDs
SELECT 
    'NULL customer IDs in usage' as check_name,
    COUNT(*) as issue_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM usage 
WHERE customer_id IS NULL
UNION ALL
SELECT 
    'NULL customer IDs in billing' as check_name,
    COUNT(*) as issue_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM billing 
WHERE customer_id IS NULL
UNION ALL
SELECT 
    'NULL customer IDs in support_tickets' as check_name,
    COUNT(*) as issue_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM support_tickets 
WHERE customer_id IS NULL;

-- Check for invalid payment statuses
SELECT 
    'Invalid payment statuses' as check_name,
    COUNT(*) as issue_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM billing 
WHERE payment_status NOT IN ('Paid', 'paid', 'Unpaid', 'unpaid');

-- Check for negative amounts
SELECT 
    'Negative billing amounts' as check_name,
    COUNT(*) as issue_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM billing 
WHERE amount < 0;

-- Verify foreign key relationships
\echo ''
\echo '5. REFERENTIAL INTEGRITY CHECKS:'

-- Check for orphaned usage records
SELECT 
    'Orphaned usage records' as check_name,
    COUNT(*) as issue_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM usage u
LEFT JOIN customers c ON u.customer_id = c.customer_id
WHERE c.customer_id IS NULL
UNION ALL
-- Check for orphaned billing records
SELECT 
    'Orphaned billing records' as check_name,
    COUNT(*) as issue_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM billing b
LEFT JOIN customers c ON b.customer_id = c.customer_id
WHERE c.customer_id IS NULL
UNION ALL
-- Check for orphaned support tickets
SELECT 
    'Orphaned support tickets' as check_name,
    COUNT(*) as issue_count,
    CASE WHEN COUNT(*) = 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM support_tickets s
LEFT JOIN customers c ON s.customer_id = c.customer_id
WHERE c.customer_id IS NULL;

-- Test analytical queries
\echo ''
\echo '6. ANALYTICAL QUERY TESTS:'

-- Test customer summary view
SELECT 
    'Customer Summary View' as test_name,
    COUNT(*) as result_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM customer_summary
WHERE avg_call_minutes IS NOT NULL;

-- Test monthly KPIs view
SELECT 
    'Monthly KPIs View' as test_name,
    COUNT(*) as result_count,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM monthly_kpis
WHERE active_users > 0;

-- Test churn risk indicators
SELECT 
    'Churn Risk Indicators' as test_name,
    COUNT(*) as result_count,
    CASE WHEN COUNT(*) >= 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM churn_risk_indicators;

-- Business metrics validation
\echo ''
\echo '7. BUSINESS METRICS VALIDATION:'

-- Total customers
SELECT 
    'Total Customers' as metric,
    COUNT(*) as value,
    CASE WHEN COUNT(*) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM customers;

-- Total revenue
SELECT 
    'Total Revenue (Paid)' as metric,
    ROUND(SUM(amount), 2) as value,
    CASE WHEN SUM(amount) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM billing
WHERE payment_status IN ('Paid', 'paid');

-- Average ARPU
SELECT 
    'Average ARPU' as metric,
    ROUND(SUM(b.amount) / COUNT(DISTINCT c.customer_id), 2) as value,
    CASE WHEN SUM(b.amount) / COUNT(DISTINCT c.customer_id) > 0 THEN '✅ PASS' ELSE '❌ FAIL' END as status
FROM customers c
JOIN billing b ON c.customer_id = b.customer_id
WHERE b.payment_status IN ('Paid', 'paid');

-- Sample data preview
\echo ''
\echo '8. SAMPLE DATA PREVIEW:'

\echo 'Sample Customers:'
SELECT customer_id, name, region, plan_type, signup_date 
FROM customers 
LIMIT 3;

\echo 'Sample Usage Data:'
SELECT customer_id, usage_month, call_minutes, data_usage_gb, num_sms 
FROM usage 
LIMIT 3;

\echo 'Sample Billing Data:'
SELECT customer_id, invoice_date, amount, payment_status 
FROM billing 
LIMIT 3;

\echo 'Sample Support Tickets:'
SELECT customer_id, ticket_date, issue_type, resolution_time_hrs 
FROM support_tickets 
LIMIT 3;

\echo ''
\echo '=============================================='
\echo 'VALIDATION COMPLETE'
\echo '=============================================='
\echo 'If all tests show ✅ PASS, your setup is ready!'
\echo 'You can now connect BI tools and create dashboards.'
\echo '=============================================='
