-- Cosmofone Telecom Analytics - Additional Queries
-- These queries support dashboard creation and business insights

-- ==============================================
-- EXECUTIVE OVERVIEW QUERIES
-- ==============================================

-- Total Active Customers
SELECT COUNT(DISTINCT customer_id) as total_active_customers
FROM customers;

-- Monthly Revenue Trends (Last 12 months)
SELECT 
    DATE_TRUNC('month', invoice_date) as month,
    SUM(amount) as monthly_revenue,
    COUNT(DISTINCT customer_id) as paying_customers,
    ROUND(SUM(amount) / COUNT(DISTINCT customer_id), 2) as avg_revenue_per_customer
FROM billing
WHERE payment_status = 'Paid'
    AND invoice_date >= CURRENT_DATE - INTERVAL '12 months'
GROUP BY month
ORDER BY month;

-- Top 5 Regions by ARPU (Average Revenue Per User)
SELECT 
    c.region,
    COUNT(DISTINCT c.customer_id) as total_customers,
    SUM(b.amount) as total_revenue,
    ROUND(SUM(b.amount) / COUNT(DISTINCT c.customer_id), 2) as arpu
FROM customers c
JOIN billing b ON c.customer_id = b.customer_id
WHERE b.payment_status = 'Paid'
GROUP BY c.region
ORDER BY arpu DESC
LIMIT 5;

-- Customer Growth by Month
SELECT 
    DATE_TRUNC('month', signup_date) as month,
    COUNT(*) as new_customers
FROM customers
WHERE signup_date >= CURRENT_DATE - INTERVAL '12 months'
GROUP BY month
ORDER BY month;

-- ==============================================
-- CUSTOMER INSIGHTS QUERIES
-- ==============================================

-- Customer Segmentation by Usage Behavior
SELECT 
    c.plan_type,
    CASE 
        WHEN AVG(u.call_minutes) > 1000 THEN 'Heavy Users'
        WHEN AVG(u.call_minutes) > 500 THEN 'Moderate Users'
        ELSE 'Light Users'
    END as usage_segment,
    COUNT(DISTINCT c.customer_id) as customer_count,
    ROUND(AVG(u.call_minutes), 2) as avg_call_minutes,
    ROUND(AVG(u.data_usage_gb), 2) as avg_data_usage,
    ROUND(AVG(u.num_sms), 0) as avg_sms
FROM customers c
JOIN usage u ON c.customer_id = u.customer_id
GROUP BY c.plan_type, 
    CASE 
        WHEN AVG(u.call_minutes) > 1000 THEN 'Heavy Users'
        WHEN AVG(u.call_minutes) > 500 THEN 'Moderate Users'
        ELSE 'Light Users'
    END
ORDER BY c.plan_type, customer_count DESC;

-- Support Ticket Resolution Performance
SELECT 
    issue_type,
    COUNT(*) as total_tickets,
    ROUND(AVG(resolution_time_hrs), 2) as avg_resolution_hours,
    COUNT(CASE WHEN resolution_time_hrs <= 8 THEN 1 END) as resolved_within_8hrs,
    ROUND(COUNT(CASE WHEN resolution_time_hrs <= 8 THEN 1 END) * 100.0 / COUNT(*), 2) as resolution_rate_8hrs
FROM support_tickets
WHERE resolution_time_hrs IS NOT NULL
GROUP BY issue_type
ORDER BY avg_resolution_hours;

-- Payment Behavior Analysis
SELECT 
    c.plan_type,
    COUNT(DISTINCT c.customer_id) as total_customers,
    COUNT(DISTINCT CASE WHEN b.payment_status = 'Paid' THEN c.customer_id END) as paying_customers,
    COUNT(DISTINCT CASE WHEN b.payment_status = 'Unpaid' THEN c.customer_id END) as customers_with_unpaid,
    ROUND(COUNT(DISTINCT CASE WHEN b.payment_status = 'Paid' THEN c.customer_id END) * 100.0 / 
          COUNT(DISTINCT c.customer_id), 2) as payment_rate
FROM customers c
LEFT JOIN billing b ON c.customer_id = b.customer_id
GROUP BY c.plan_type;

-- ==============================================
-- PLAN PERFORMANCE QUERIES
-- ==============================================

-- Average Usage and Billing by Plan Type
SELECT 
    c.plan_type,
    COUNT(DISTINCT c.customer_id) as customer_count,
    ROUND(AVG(u.call_minutes), 2) as avg_call_minutes,
    ROUND(AVG(u.data_usage_gb), 2) as avg_data_usage_gb,
    ROUND(AVG(u.num_sms), 0) as avg_sms_count,
    ROUND(AVG(b.amount), 2) as avg_bill_amount,
    SUM(b.amount) as total_revenue
FROM customers c
LEFT JOIN usage u ON c.customer_id = u.customer_id
LEFT JOIN billing b ON c.customer_id = b.customer_id AND b.payment_status = 'Paid'
GROUP BY c.plan_type
ORDER BY total_revenue DESC;

-- Plan Conversion Analysis (Prepaid to Postpaid)
WITH customer_plan_history AS (
    SELECT 
        customer_id,
        plan_type,
        signup_date,
        ROW_NUMBER() OVER (PARTITION BY customer_id ORDER BY signup_date) as plan_sequence
    FROM customers
)
SELECT 
    'Prepaid to Postpaid Conversion' as metric,
    COUNT(CASE WHEN plan_sequence = 1 AND plan_type = 'Prepaid' THEN 1 END) as initial_prepaid,
    COUNT(CASE WHEN plan_sequence > 1 AND plan_type = 'Postpaid' THEN 1 END) as converted_to_postpaid,
    ROUND(COUNT(CASE WHEN plan_sequence > 1 AND plan_type = 'Postpaid' THEN 1 END) * 100.0 / 
          COUNT(CASE WHEN plan_sequence = 1 AND plan_type = 'Prepaid' THEN 1 END), 2) as conversion_rate
FROM customer_plan_history;

-- Usage Pattern Comparison by Plan
SELECT 
    c.plan_type,
    CASE 
        WHEN u.usage_month >= CURRENT_DATE - INTERVAL '30 days' THEN 'Last 30 Days'
        WHEN u.usage_month >= CURRENT_DATE - INTERVAL '60 days' THEN '31-60 Days Ago'
        ELSE 'Older'
    END as time_period,
    COUNT(DISTINCT u.customer_id) as active_users,
    ROUND(AVG(u.call_minutes), 2) as avg_call_minutes,
    ROUND(AVG(u.data_usage_gb), 2) as avg_data_usage
FROM customers c
JOIN usage u ON c.customer_id = u.customer_id
WHERE u.usage_month >= CURRENT_DATE - INTERVAL '90 days'
GROUP BY c.plan_type, 
    CASE 
        WHEN u.usage_month >= CURRENT_DATE - INTERVAL '30 days' THEN 'Last 30 Days'
        WHEN u.usage_month >= CURRENT_DATE - INTERVAL '60 days' THEN '31-60 Days Ago'
        ELSE 'Older'
    END
ORDER BY c.plan_type, time_period;

-- ==============================================
-- CHURN ANALYSIS QUERIES
-- ==============================================

-- Detailed Churn Risk Analysis
SELECT 
    c.customer_id,
    c.name,
    c.region,
    c.plan_type,
    cs.recent_tickets,
    cs.unpaid_bills,
    cs.total_billed,
    CASE 
        WHEN cs.recent_tickets > 3 OR cs.unpaid_bills > 2 THEN 'High Risk'
        WHEN cs.recent_tickets > 1 OR cs.unpaid_bills > 0 THEN 'Medium Risk'
        ELSE 'Low Risk'
    END as risk_level,
    ROUND(cs.avg_call_minutes, 2) as avg_call_minutes,
    ROUND(cs.avg_data_usage_gb, 2) as avg_data_usage
FROM customers c
JOIN customer_summary cs ON c.customer_id = cs.customer_id
WHERE cs.recent_tickets > 0 OR cs.unpaid_bills > 0
ORDER BY cs.recent_tickets DESC, cs.unpaid_bills DESC;

-- Churn Risk by Region and Plan
SELECT 
    c.region,
    c.plan_type,
    COUNT(*) as total_customers,
    COUNT(CASE WHEN cri.recent_tickets > 2 OR cri.unpaid_bills > 1 THEN 1 END) as at_risk_customers,
    ROUND(COUNT(CASE WHEN cri.recent_tickets > 2 OR cri.unpaid_bills > 1 THEN 1 END) * 100.0 / 
          COUNT(*), 2) as churn_risk_percentage
FROM customers c
LEFT JOIN churn_risk_indicators cri ON c.customer_id = cri.customer_id
GROUP BY c.region, c.plan_type
ORDER BY churn_risk_percentage DESC;

-- ==============================================
-- REVENUE ANALYSIS QUERIES
-- ==============================================

-- Revenue by Region and Plan Type
SELECT 
    c.region,
    c.plan_type,
    COUNT(DISTINCT c.customer_id) as customers,
    SUM(b.amount) as total_revenue,
    ROUND(SUM(b.amount) / COUNT(DISTINCT c.customer_id), 2) as revenue_per_customer,
    ROUND(SUM(b.amount) * 100.0 / SUM(SUM(b.amount)) OVER (), 2) as revenue_percentage
FROM customers c
JOIN billing b ON c.customer_id = b.customer_id
WHERE b.payment_status = 'Paid'
GROUP BY c.region, c.plan_type
ORDER BY total_revenue DESC;

-- Monthly Revenue Breakdown
SELECT 
    DATE_TRUNC('month', b.invoice_date) as month,
    c.plan_type,
    SUM(b.amount) as revenue,
    COUNT(DISTINCT b.customer_id) as paying_customers
FROM billing b
JOIN customers c ON b.customer_id = c.customer_id
WHERE b.payment_status = 'Paid'
    AND b.invoice_date >= CURRENT_DATE - INTERVAL '12 months'
GROUP BY month, c.plan_type
ORDER BY month, revenue DESC;

-- ==============================================
-- OPERATIONAL METRICS
-- ==============================================

-- Support Ticket Trends
SELECT 
    DATE_TRUNC('month', ticket_date) as month,
    issue_type,
    COUNT(*) as ticket_count,
    ROUND(AVG(resolution_time_hrs), 2) as avg_resolution_hours
FROM support_tickets
WHERE ticket_date >= CURRENT_DATE - INTERVAL '12 months'
GROUP BY month, issue_type
ORDER BY month, ticket_count DESC;

-- Customer Lifetime Value Analysis
WITH customer_metrics AS (
    SELECT 
        c.customer_id,
        c.plan_type,
        c.signup_date,
        SUM(b.amount) as total_revenue,
        COUNT(DISTINCT u.usage_month) as active_months,
        COUNT(s.ticket_date) as support_tickets
    FROM customers c
    LEFT JOIN billing b ON c.customer_id = b.customer_id AND b.payment_status = 'Paid'
    LEFT JOIN usage u ON c.customer_id = u.customer_id
    LEFT JOIN support_tickets s ON c.customer_id = s.customer_id
    GROUP BY c.customer_id, c.plan_type, c.signup_date
)
SELECT 
    plan_type,
    COUNT(*) as customers,
    ROUND(AVG(total_revenue), 2) as avg_lifetime_value,
    ROUND(AVG(active_months), 1) as avg_active_months,
    ROUND(AVG(support_tickets), 1) as avg_support_tickets,
    ROUND(AVG(total_revenue / NULLIF(active_months, 0)), 2) as avg_monthly_value
FROM customer_metrics
GROUP BY plan_type
ORDER BY avg_lifetime_value DESC;
