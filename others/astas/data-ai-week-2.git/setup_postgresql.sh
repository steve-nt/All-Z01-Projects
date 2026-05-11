#!/bin/bash

# Cosmofone Telecom Data Analytics Setup Script
# This script sets up PostgreSQL and loads the telecom data

echo "🚀 Setting up Cosmofone Telecom Data Analytics Environment..."

# Check if PostgreSQL is installed
if ! command -v psql &> /dev/null; then
    echo "📦 Installing PostgreSQL..."
    sudo apt update
    sudo apt install -y postgresql postgresql-contrib
    sudo systemctl start postgresql
    sudo systemctl enable postgresql
else
    echo "✅ PostgreSQL is already installed"
fi

# Create database and user
echo "🗄️ Creating database and user..."
sudo -u postgres psql -c "CREATE DATABASE cosmofone;"
sudo -u postgres psql -c "CREATE USER gina WITH PASSWORD 'ginapass';"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE cosmofone TO gina;"
sudo -u postgres psql -c "ALTER USER gina CREATEDB;"
# Grant ownership of public schema to gina
sudo -u postgres psql -d cosmofone -c "ALTER SCHEMA public OWNER TO gina;"
sudo -u postgres psql -d cosmofone -c "GRANT ALL ON SCHEMA public TO gina;"

cd db
echo "📊 Creating schema and tables..."
psql -h localhost -U gina -d cosmofone -f telecom_schema_and_views.sql

echo "📈 Loading data..."
psql -h localhost -U gina -d cosmofone -f load_data.sql

echo "📈 Commiting data quality checks..."
psql -h localhost -U gina -d cosmofone -f data_quality_check.sql

echo "✅ Setup complete!"
echo ""
echo "🔗 Connection details:"
echo "   Host: localhost"
echo "   Port: 5432"
echo "   Database: cosmofone"
echo "   Username: gina"
echo "   Password: ginapass"
echo ""
echo "📊 You can now connect to the database and start building dashboards!"