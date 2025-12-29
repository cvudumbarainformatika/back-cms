# Docker Quick Start - For Your Setup

**TL;DR - Langsung eksekusi ini:**

```bash
# 1. Setup .env
cp .env.example .env

# Edit .env, gunakan nilai ini:
# DB_HOST=mysql
# DB_PORT=3306
# DB_DATABASE=sasacms
# DB_USERNAME=admin
# DB_PASSWORD=sasa0102
# REDIS_HOST=redis

# 2. Start Docker
docker-compose up -d

# 3. Wait ~30 seconds, then check
docker-compose logs -f

# 4. Verify tables created
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms -e "SHOW TABLES;"

# 5. Test API
curl http://localhost:8080/health
```

---

## **Container Info**

| Service | Host | Port (External) | Port (Internal) | Credentials |
|---------|------|-----------------|-----------------|-------------|
| MySQL | localhost | **33067** | 3306 | User: `admin` / Pass: `sasa0102` |
| Redis | localhost | 6379 | 6379 | No auth |
| API | localhost | **8080** | 8080 | - |

**Important:** Di `.env`, gunakan port INTERNAL (3306), bukan external (33067)!

---

## **Essential Commands**

```bash
# Start
docker-compose up -d

# Logs
docker-compose logs -f app          # App logs
docker-compose logs -f mysql        # MySQL logs

# Stop
docker-compose down

# Reset everything
docker-compose down -v              # Delete volumes too

# Access MySQL
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms

# Access app container
docker exec -it starter-api sh
```

---

## **Verify Setup**

```bash
# Check all tables created
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms -e "SHOW TABLES;"

# Check tables count (should be 16)
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms -e "SHOW TABLES;" | wc -l

# Check specific table structure
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms -e "DESC users;"
```

---

## **If Something Goes Wrong**

```bash
# View error logs
docker-compose logs app

# Restart containers
docker-compose restart

# Full reset (WARNING: deletes data)
docker-compose down -v
docker-compose up -d
```

---

## **Connect via GUI Tools**

### MySQL Workbench / DBeaver / TablePlus

**Connection Settings:**
```
Host: localhost
Port: 33067
Username: admin
Password: sasa0102
Database: sasacms
```

### phpMyAdmin (Optional - add to docker-compose if needed)

```yaml
phpmyadmin:
  image: phpmyadmin:latest
  ports:
    - "8081:80"
  environment:
    PMA_HOST: mysql
    PMA_USER: admin
    PMA_PASSWORD: sasa0102
  depends_on:
    - mysql
  networks:
    - starter-network
```

Then access at: `http://localhost:8081`

---

Done! Now you can start development. 🚀
