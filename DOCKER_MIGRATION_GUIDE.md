# Docker Migration Guide

Karena Anda menggunakan Docker, migrations akan otomatis berjalan ketika container start. Berikut panduan lengkapnya:

---

## **STEP 1: Update .env File**

Copy dari `.env.example`:
```bash
cp .env.example .env
```

Update `.env` sesuai dengan `docker-compose.yml` Anda:

```env
# APPLICATION
APP_NAME=Go Gin Backend
APP_ENV=local
APP_PORT=8080

# DATABASE (sesuaikan dengan docker-compose.yml)
DB_CONNECTION=mysql
DB_HOST=mysql              # ← Service name dari docker-compose
DB_PORT=3306              # Internal port (bukan 33067)
DB_DATABASE=sasacms       # ← Dari MYSQL_DATABASE di docker-compose
DB_USERNAME=admin         # ← Dari MYSQL_USER di docker-compose
DB_PASSWORD=sasa0102      # ← Dari MYSQL_PASSWORD di docker-compose

# JWT
JWT_SECRET=your-secret-key-change-this
JWT_ACCESS_TOKEN_EXPIRATION=15
JWT_REFRESH_TOKEN_EXPIRATION=10080

# REDIS (sesuaikan dengan docker-compose)
REDIS_HOST=redis          # ← Service name dari docker-compose
REDIS_PORT=6379          # Internal port
REDIS_PASSWORD=
REDIS_DB=0

# Rate Limiting & CORS
RATE_LIMIT_ENABLED=true
RATE_LIMIT_MAX_REQUESTS=100
RATE_LIMIT_WINDOW_SECONDS=60
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
```

**PENTING:** 
- `DB_HOST=mysql` (bukan localhost) - karena dalam Docker, hostname adalah service name
- `DB_PORT=3306` (internal port) - bukan 33067 (external port)
- Database: `sasacms`, User: `admin`, Password: `sasa0102`

---

## **STEP 2: Start Docker Containers**

```bash
# Build dan start semua containers
docker-compose up -d

# Lihat logs (untuk verifikasi)
docker-compose logs -f
```

Containers yang akan jalan:
- `starter-mysql` - MySQL database
- `starter-redis` - Redis cache
- `starter-api` - Go application

---

## **STEP 3: Tunggu MySQL Siap (Health Check)**

docker-compose sudah punya health check, tapi tunggu sampe MySQL fully ready:

```bash
# Check status containers
docker-compose ps

# Seharusnya semua status "healthy" atau "Up"
```

Atau tunggu logs sampai ada output:
```
starter-api | Database connection established successfully
starter-api | ✓ All migrations completed successfully!
```

---

## **STEP 4: Verify Migrations Berhasil**

### Option A: Melalui Docker CLI

```bash
# Masuk ke MySQL container
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms -e "SHOW TABLES;"
```

Seharusnya output:
```
+---------------------------+
| Tables_in_sasacms        |
+---------------------------+
| agenda                   |
| agenda_registrations     |
| berita                   |
| berita_tag_map           |
| berita_tags              |
| direktori                |
| documents                |
| dynamic_contents         |
| homepage                 |
| menus                    |
| migrations               |
| pengurus                 |
| profil_organisasi        |
| uploads                  |
| users                    |
| user_sessions            |
+---------------------------+
```

### Option B: Melalui GUI Tool (MySQL Workbench / DBeaver)

Connect ke MySQL:
- Host: `localhost`
- Port: `33067` (dari docker-compose.yml)
- Username: `admin`
- Password: `sasa0102`
- Database: `sasacms`

---

## **STEP 5: Test API**

```bash
# Health check
curl http://localhost:8080/health

# Expected response:
# {"status":"ok","database":"connected"}
```

---

## **Troubleshooting**

### Error: "Can't connect to MySQL server"

Tunggu lebih lama (MySQL butuh waktu start):
```bash
docker-compose ps  # Check status
docker-compose logs mysql  # Lihat MySQL logs
```

### Error: "Unknown database 'sasacms'"

Database akan dibuat otomatis oleh docker-compose (MYSQL_DATABASE env var). Jika belum, buat manual:

```bash
docker exec -it starter-mysql mysql -u root -psasa0102 -e "CREATE DATABASE sasacms CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
```

### Error: "Migrations gagal"

Check logs aplikasi:
```bash
docker-compose logs app
```

Jika ada error, Anda bisa:
1. Stop containers: `docker-compose down`
2. Delete database: `docker volume rm starter_mysql_data`
3. Start ulang: `docker-compose up -d`

---

## **Common Commands**

```bash
# Start containers
docker-compose up -d

# Stop containers
docker-compose down

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f app
docker-compose logs -f mysql

# Access MySQL CLI
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms

# Access application container
docker exec -it starter-api sh

# Rebuild containers (jika ada code changes)
docker-compose up -d --build

# Remove everything (including volumes/data)
docker-compose down -v

# Restart containers
docker-compose restart
```

---

## **STEP 6 (Optional): Manual Migration if Auto-migrate Fails**

Jika auto-migrate tidak berjalan (misal ada error), Anda bisa manual:

```bash
# Akses MySQL container
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms

# Kemudian copy-paste isi setiap migration file dari database/migrations/
```

Atau gunakan bash script:
```bash
# Copy migrations ke container
docker cp database/migrations/. starter-mysql:/tmp/migrations

# Jalankan
docker exec starter-mysql bash -c 'for file in /tmp/migrations/00*.sql; do mysql -u admin -psasa0102 sasacms < $file; done'
```

---

## **Summary**

Dengan Docker, Anda tinggal:
1. ✅ Update `.env` dengan benar
2. ✅ Jalankan: `docker-compose up -d`
3. ✅ Tunggu ~30 detik sampe MySQL fully ready
4. ✅ Check logs: `docker-compose logs -f`
5. ✅ Verify tables: `docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms -e "SHOW TABLES;"`

**That's it!** Migrations otomatis jalan saat container start. 🎉

---

## **Next Steps**

Setelah database ready, Anda bisa:
- Test auth endpoints
- Implement berita endpoints
- Test dengan Postman/Thunder Client

Siap lanjut ke Task berikutnya? 🚀
