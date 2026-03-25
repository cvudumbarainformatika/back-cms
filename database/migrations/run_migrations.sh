#!/bin/sh

# Database Migration Runner
# This script runs all SQL migration files in order (POSIX sh compliant for Alpine)

# Load environment variables if not already in environment
if [ -f .env ]; then
    # Simple .env loader for sh
    while read -r line || [ -n "$line" ]; do
        case "$line" in
            \#*) ;;
            "") ;;
            *) export "$line" ;;
        esac
    done < .env
fi

# Use defaults if env vars not set
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-3306}
DB_USERNAME=${DB_USERNAME:-root}
DB_PASSWORD=${DB_PASSWORD:-secret}
DB_DATABASE=${DB_DATABASE:-go_backend_db}

# Colors for output (using printf)
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

printf "${YELLOW}========================================${NC}\n"
printf "${YELLOW}Database Migration Runner (sh)${NC}\n"
printf "${YELLOW}========================================${NC}\n\n"

printf "Database: %s\n" "$DB_DATABASE"
printf "Host: %s:%s\n" "$DB_HOST" "$DB_PORT"
printf "User: %s\n\n" "$DB_USERNAME"

# Check if password is set and build mysql command
if [ -z "$DB_PASSWORD" ]; then
    MYSQL_CMD="mysql --skip-ssl -h $DB_HOST -P $DB_PORT -u $DB_USERNAME $DB_DATABASE"
else
    MYSQL_CMD="mysql --skip-ssl -h $DB_HOST -P $DB_PORT -u $DB_USERNAME -p$DB_PASSWORD $DB_DATABASE"
fi

# Check if mysql command exists
if ! command -v mysql > /dev/null 2>&1; then
    printf "${RED}❌ Error: 'mysql' command not found!${NC}\n"
    printf "Please install mysql-client in your environment.\n"
    printf "Inside Docker Alpine, run: apk add mysql-client\n"
    exit 1
fi

# Test connection
printf "${YELLOW}Testing database connection...${NC}\n"
if ! $MYSQL_CMD -e "SELECT 1;" > /dev/null 2>&1; then
    printf "${RED}❌ Failed to connect to database!${NC}\n"
    printf "Please check your database credentials in .env file\n"
    exit 1
fi
printf "${GREEN}✓ Database connection successful${NC}\n\n"

# Run migrations
MIGRATION_DIR="database/migrations"

printf "${YELLOW}Running migrations...${NC}\n\n"

FAILED=0
# Find all .sql files, sort them, and loop
FILES=$(ls $MIGRATION_DIR/*.sql | sort)
for MIGRATION_FILE in $FILES; do
    migration=$(basename "$MIGRATION_FILE")
    
    if [ ! -f "$MIGRATION_FILE" ]; then
        continue
    fi
    
    printf "Running: %s ... " "$migration"
    
    if $MYSQL_CMD < "$MIGRATION_FILE" > /dev/null 2>&1; then
        printf "${GREEN}✓${NC}\n"
    else
        printf "${RED}✗${NC}\n"
        printf "${RED}   Error running migration: %s${NC}\n" "$migration"
        FAILED=$((FAILED + 1))
    fi
done

printf "\n${YELLOW}========================================${NC}\n"

if [ $FAILED -eq 0 ]; then
    printf "${GREEN}✓ All migrations completed successfully!${NC}\n\n"
    printf "Tables created:\n"
    $MYSQL_CMD -e "SHOW TABLES;" | tail -n +2 | sed 's/^/  - /'
    printf "\n"
    exit 0
else
    printf "${RED}✗ %d migration(s) failed!${NC}\n\n" "$FAILED"
    exit 1
fi
