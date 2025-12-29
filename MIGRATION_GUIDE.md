# Database Migration Guide

Database Anda masih kosong dan perlu diisi dengan tables sesuai schema. Ada beberapa cara untuk menjalankan migrasi:

## Step 1: Setup Database Connection

Pastikan `.env` file sudah dikonfigurasi dengan benar:

```env
DB_CONNECTION=mysql
DB_HOST=localhost
DB_PORT=3306
DB_DATABASE=go_backend_db      # ← Buat database ini dulu jika belum ada
DB_USERNAME=root
DB_PASSWORD=secret
```

### Buat Database (jika belum ada)

Gunakan MySQL CLI atau GUI tool (phpMyAdmin, MySQL Workbench):

```sql
CREATE DATABASE go_backend_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

---

## Method 1: Menggunakan MySQL Command Line (Recommended untuk Development)

### Prerequisites:
- MySQL CLI sudah terinstall
- Database sudah dibuat

### Steps:

1. Navigasi ke project root:
```bash
cd /path/to/project
```

2. Jalankan semua migration files secara berurutan:
```bash
mysql -u root -p go_backend_db < database/migrations/001_create_users_table.sql
mysql -u root -p go_backend_db < database/migrations/002_create_berita_table.sql
mysql -u root -p go_backend_db < database/migrations/003_create_berita_tags_table.sql
mysql -u root -p go_backend_db < database/migrations/004_create_agenda_table.sql
mysql -u root -p go_backend_db < database/migrations/005_create_direktori_table.sql
mysql -u root -p go_backend_db < database/migrations/006_create_pengurus_table.sql
mysql -u root -p go_backend_db < database/migrations/007_create_documents_table.sql
mysql -u root -p go_backend_db < database/migrations/008_create_menus_table.sql
mysql -u root -p go_backend_db < database/migrations/009_create_dynamic_contents_table.sql
mysql -u root -p go_backend_db < database/migrations/010_create_homepage_table.sql
mysql -u root -p go_backend_db < database/migrations/011_create_user_sessions_table.sql
```

Atau gunakan bash script:
```bash
bash database/migrations/run_migrations.sh
```

---

## Method 2: Menggunakan GUI Tool (phpMyAdmin/MySQL Workbench)

1. Login ke MySQL GUI tool
2. Buat database baru: `go_backend_db`
3. Buka setiap file SQL dari `database/migrations/` folder
4. Copy-paste isi file ke query editor
5. Jalankan query

---

## Method 3: Menggunakan Go Migration Tool (Coming Soon)

Saya akan membuat automatic migration tool di Go. Untuk sekarang gunakan Method 1 atau 2.

---

## Verify Migration

Setelah menjalankan semua migration, verify tables sudah terbuat:

```bash
mysql -u root -p go_backend_db -e "SHOW TABLES;"
```

Expected output:
```
+-----------------------------+
| Tables_in_go_backend_db     |
+-----------------------------+
| agenda                      |
| agenda_registrations        |
| berita                      |
| berita_tag_map              |
| berita_tags                 |
| direktori                   |
| documents                   |
| dynamic_contents            |
| homepage                    |
| menus                       |
| pengurus                    |
| profil_organisasi           |
| uploads                     |
| users                       |
| user_sessions               |
+-----------------------------+
```

---

## Troubleshooting

### Error: "database does not exist"
Pastikan database sudah dibuat:
```sql
CREATE DATABASE go_backend_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### Error: "Access denied for user"
Periksa credential di `.env` file, pastikan sesuai dengan MySQL user Anda.

### Error: "Table already exists"
Jika ingin reset, drop semua tables dulu:
```sql
DROP DATABASE go_backend_db;
CREATE DATABASE go_backend_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

---

## Next Steps

Setelah migrasi selesai:
1. Start application: `go run main.go`
2. Test dengan health check endpoint: `GET http://localhost:8080/health`
3. Mulai menggunakan endpoints (Auth, Berita, dll)
