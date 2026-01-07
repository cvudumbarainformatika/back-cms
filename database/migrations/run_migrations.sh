#!/bin/bash

# Database Migration Runner
# This script runs all SQL migration files in order

# Load environment variables from .env if it exists
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Use defaults if env vars not set
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USERNAME=${DB_USERNAME:-root}
DB_PASSWORD=${DB_PASSWORD:-secret}
DB_DATABASE=${DB_DATABASE:-go_backend_db}

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}Database Migration Runner${NC}"
echo -e "${YELLOW}========================================${NC}"
echo ""
echo "Database: $DB_DATABASE"
echo "Host: $DB_HOST:$DB_PORT"
echo "User: $DB_USERNAME"
echo ""

# Check if password is set and build mysql command
if [ -z "$DB_PASSWORD" ]; then
    MYSQL_CMD="mysql -h $DB_HOST -P $DB_PORT -u $DB_USERNAME $DB_DATABASE"
else
    MYSQL_CMD="mysql -h $DB_HOST -P $DB_PORT -u $DB_USERNAME -p$DB_PASSWORD $DB_DATABASE"
fi

# Test connection
echo -e "${YELLOW}Testing database connection...${NC}"
if ! $MYSQL_CMD -e "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${RED}❌ Failed to connect to database!${NC}"
    echo "Please check your database credentials in .env file"
    exit 1
fi
echo -e "${GREEN}✓ Database connection successful${NC}"
echo ""

# Run migrations
MIGRATION_DIR="database/migrations"
MIGRATIONS=(
    "001_create_users_table.sql"
    "002_create_berita_table.sql"
    "003_create_berita_tags_table.sql"
    "004_create_agenda_table.sql"
    "005_create_direktori_table.sql"
    "006_create_pengurus_table.sql"
    "007_create_documents_table.sql"
    "008_create_menus_table.sql"
    "009_create_dynamic_contents_table.sql"
    "010_create_homepage_table.sql"
    "011_create_user_sessions_table.sql"
)

echo -e "${YELLOW}Running migrations...${NC}"
echo ""

FAILED=0
for migration in "${MIGRATIONS[@]}"; do
    MIGRATION_FILE="$MIGRATION_DIR/$migration"
    
    if [ ! -f "$MIGRATION_FILE" ]; then
        echo -e "${RED}❌ Migration file not found: $MIGRATION_FILE${NC}"
        FAILED=$((FAILED+1))
        continue
    fi
    
    echo -n "Running: $migration ... "
    
    if $MYSQL_CMD < "$MIGRATION_FILE" > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
    else
        echo -e "${RED}✗${NC}"
        echo -e "${RED}   Error running migration: $migration${NC}"
        FAILED=$((FAILED+1))
    fi
done

echo ""
echo -e "${YELLOW}========================================${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}✓ All migrations completed successfully!${NC}"
    echo ""
    echo "Tables created:"
    $MYSQL_CMD -e "SHOW TABLES;" | tail -n +2 | sed 's/^/  - /'
    echo ""
    echo "You can now start the application with: go run main.go"
    exit 0
else
    echo -e "${RED}✗ $FAILED migration(s) failed!${NC}"
    echo ""
    echo "Please check the errors above and try again."
    exit 1
fi
