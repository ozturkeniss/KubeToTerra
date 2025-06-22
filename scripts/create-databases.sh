#!/bin/bash

export PGPASSWORD=${DB_PASSWORD:-postgres}

# Database creation script for BWENG microservices
# This script creates the necessary databases for user and order services

set -e  # Exit on any error

# Database configuration
DB_HOST=${DB_HOST:-127.0.0.1}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-postgres}

# Database names
USER_DB=${USER_DB:-bweng_user_db}
ORDER_DB=${ORDER_DB:-bweng_order_db}

echo "🚀 Starting database creation for BWENG microservices..."
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USER"
echo "User Database: $USER_DB"
echo "Order Database: $ORDER_DB"
echo ""

# Function to create database
create_database() {
    local db_name=$1
    local service_name=$2
    
    echo "📦 Creating $service_name database: $db_name"
    
    # Check if database already exists
    if psql -p "$DB_PORT" -U "$DB_USER" -lqt | cut -d \| -f 1 | grep -qw "$db_name"; then
        echo "✅ Database '$db_name' already exists, skipping..."
    else
        # Create database
        createdb -p "$DB_PORT" -U "$DB_USER" "$db_name"
        echo "✅ Database '$db_name' created successfully!"
    fi
}

# Function to test database connection
test_connection() {
    echo "🔍 Testing database connection..."
    if psql -p "$DB_PORT" -U "$DB_USER" -c "SELECT version();" > /dev/null 2>&1; then
        echo "✅ Database connection successful!"
        return 0
    else
        echo "❌ Failed to connect to database!"
        echo "Please check your PostgreSQL connection settings."
        return 1
    fi
}

# Main execution
main() {
    # Test connection first
    if ! test_connection; then
        exit 1
    fi
    
    echo ""
    
    # Create databases
    create_database "$USER_DB" "User Service"
    create_database "$ORDER_DB" "Order Service"
    
    echo ""
    echo "🎉 All databases created successfully!"
    echo ""
    echo "📋 Summary:"
    echo "   - User Service DB: $USER_DB"
    echo "   - Order Service DB: $ORDER_DB"
    echo ""
    echo "🔧 Next steps:"
    echo "   1. Start your microservices"
    echo "   2. They will automatically run migrations"
    echo "   3. Or use docker-compose to start everything together"
}

# Run main function
main "$@" 