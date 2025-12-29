# Implementation Tasks - CMS Backend API

Status: Tracking implementation progress berdasarkan `apiyangharusdibuattask.md`

---

## 🔐 Authentication & Authorization

- [x] **Task 1:** Setup Auth Endpoints (login, register, me, logout, refresh)
- [ ] **Task 15:** Implement RBAC middleware and role-based access control
- [ ] **Task 17:** Implement consistent error handling and response format

## 👥 User Management

- [ ] **Task 2:** Implement User Management Endpoints (GET, POST, PUT, PATCH, DELETE)

## 📰 Berita / News

- [ ] **Task 3:** Implement Berita POST endpoint (create news with slug generation)
- [ ] **Task 4:** Implement Berita GET categories and tags endpoints

## 📅 Agenda

- [ ] **Task 5:** Implement Agenda Registrations endpoints (GET, POST, PATCH)

## 📂 Direktori

- [ ] **Task 6:** Implement Direktori CRUD endpoints (GET by ID, POST, PUT, DELETE)

## 👔 Pengurus / Leadership

- [ ] **Task 7:** Implement Pengurus CRUD endpoints (POST, PUT, DELETE)

## 🧭 Menu & Navigation

- [ ] **Task 8:** Implement Menu Navigasi POST endpoint with validation

## 🏢 Organization Profile & Homepage

- [ ] **Task 9:** Implement Profil Organisasi POST endpoint
- [ ] **Task 10:** Implement Homepage GET endpoint

## 📝 Dynamic Content

- [ ] **Task 11:** Implement Dynamic Content endpoints (GET, POST, PUT, DELETE)

## 📄 Documents & Uploads

- [ ] **Task 12:** Implement Documents CRUD endpoints (GET, POST, PUT, DELETE)
- [ ] **Task 13:** Implement Uploads GET and DELETE endpoints

## 🗄️ Infrastructure & Utilities

- [x] **Task 14:** Setup database migrations and models for all tables
- [x] **Task 16:** Setup slug normalization utility function

---

## ✅ Additional Tasks Completed

- [x] **Seeder Task:** Create Admin/Root users and initial menus
  - 4 admin users created (admin_pusat, admin_wilayah, admin_cabang, member)
  - 4 navigation menus seeded (Beranda, Berita, Agenda, Direktori)
  - Auto-seeding on application startup

## Notes
- Mark tasks as `[x]` ketika sudah completed
- Urutan bisa disesuaikan sesuai prioritas development
- Cross-check dengan `apiyangharusdibuattask.md` untuk requirements lengkap

---

## 🐳 Database Migration (Docker)

**Status:** ✅ READY TO RUN

Untuk menjalankan migrations dengan Docker:

```bash
# 1. Setup .env
cp .env.example .env

# 2. Edit .env dengan nilai ini:
# DB_HOST=mysql (BUKAN localhost!)
# DB_PORT=3306 (BUKAN 33067!)
# DB_DATABASE=sasacms
# DB_USERNAME=admin
# DB_PASSWORD=sasa0102
# REDIS_HOST=redis

# 3. Start Docker
docker-compose up -d

# 4. Check logs
docker-compose logs -f

# 5. Verify tables (setelah ~30 detik)
docker exec -it starter-mysql mysql -u admin -psasa0102 sasacms -e "SHOW TABLES;"
```

Lihat dokumentasi lengkap di:
- `DOCKER_QUICK_START.md` - Quick reference
- `DOCKER_MIGRATION_GUIDE.md` - Detail & troubleshooting
