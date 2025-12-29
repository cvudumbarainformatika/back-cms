# KIRIM KE BACKEND — Spesifikasi API & Database CMS

Dokumen ini merangkum kebutuhan backend untuk CMS PDPI berdasarkan implementasi frontend yang sudah ada. Fokusnya pada endpoint API, aturan autentikasi/otorisasi (RBAC), skema database/relasi, format respons, serta aturan validasi/paginasi/soft delete.

## Ringkasan Fitur
- Konten publik: Berita, Agenda, Direktori, Profil Organisasi, Menu Navigasi, Homepage, Konten Dinamis
- Dashboard member/admin: Dokumen anggota, registrasi agenda, pengelolaan konten, pengelolaan menu
- Upload file: Gambar untuk konten (validasi tipe/ukuran, struktur direktori `public/uploads/YYYY/MM/`)
- RBAC: `public`, `member`, `admin_cabang`, `admin_wilayah`, `admin_pusat`

## Prinsip API
- Format respons konsisten: `{ success, data, message, pagination? }`
- Soft delete menggunakan kolom `deleted_at`; default GET tidak menampilkan data yang terhapus
- Paginasi standar: `page` (default 1), `limit` (default bervariasi per resource)
- Filter/Sort sesuai kebutuhan resource (lihat masing-masing endpoint)
- Slug unik dan dinormalisasi (lowercase, strip non-alnum, replace spasi dengan `-`, dedup `-`)
- Error standar: 400 (bad request), 401 (unauthorized), 403 (forbidden), 404 (not found), 409 (conflict), 500 (server error)

## Autentikasi & RBAC
- JWT Access + Refresh disarankan; kirim pada header `Authorization: Bearer <token>`
- Endpoint autentikasi:
  - `POST /api/auth/login`
  - `POST /api/auth/register`
  - `GET  /api/auth/me`
  - `POST /api/auth/logout`
  - `POST /api/auth/refresh`
- Klaim token memuat: `sub`, `email`, `role`, `organization_level`, `branch_id?`, `region_id?`
- Proteksi endpoint:
  - Publik (GET): berita, agenda, direktori, menus, profil, homepage, dynamic content
  - Member: dokumen milik sendiri, registrasi agenda
  - Admin (cabang/wilayah/pusat): CRUD konten (berita, agenda, direktori, pengurus, profil, homepage, menus), manajemen users, dokumen anggota (moderasi)

## Rekomendasi Skema Database (Relasional)
Gunakan MySQL. Nama kolom disarankan snake_case. Index sesuai kolom filter.

### users
- `id` (PK), `email` (unique), `password_hash`, `name`, `avatar_url?`
- `member_id?`, `organization_level?` (`pusat|wilayah|cabang|null`)
- `branch_id?`, `region_id?`, `role` (`public|member|admin_cabang|admin_wilayah|admin_pusat`)
- `created_at`, `updated_at`, `deleted_at?`

### sessions / user_tokens
- `id` (PK), `user_id` (FK -> users), `refresh_token_hash`, `expires_at`, `created_at`

### berita
- `id` (PK), `slug` (unique), `title`, `excerpt`, `content`, `image_url?`
- `category` (mis. `umum|kegiatan|ilmiah|prestasi|...`), `author`
- `status` (`draft|published`), `published_at?`, `views` (default 0)
- `created_at`, `updated_at`, `deleted_at?`
Index: `slug`, `category`, `author`, `published_at`, `status`

### berita_tags
- `id` (PK), `name` (unique)

### berita_tag_map (M2M)
- `berita_id` (FK), `tag_id` (FK)
- Unique composite `(berita_id, tag_id)`

### agenda
- `id` (PK), `slug` (unique), `title`, `description`, `type` (`webinar|workshop|seminar|kongres`)
- `date`, `end_date?`, `is_online`, `location?`, `skp?`, `quota?`, `registered?`
- `registration_url?`, `image_url?`, `fee?`
- `status` (`draft|published`), `created_at`, `updated_at`, `deleted_at?`
Index: `slug`, `type`, `date`, `status`

### agenda_registrations
- `id` (PK), `agenda_id` (FK -> agenda), `user_id` (FK -> users)
- `status` (`pending|confirmed|cancelled`), `registered_at`
- Unique composite `(agenda_id, user_id)`

### direktori
- `id` (PK), `name`, `type` (`rumah_sakit|klinik|institusi`)
- `address`, `phone?`, `email?`, `website?`
- `city`, `province`, `has_respirologist`
- `facilities` (JSON array)
- `created_at`, `updated_at`, `deleted_at?`
Index: `province`, `city`, `type`

### pengurus
- `id` (PK), `name`, `position`, `bidang`, `level` (`pusat|wilayah|cabang`)
- `periode`, `email?`, `created_at`, `updated_at`, `deleted_at?`
Index: `level`, `periode`

### documents
- `id` (PK), `user_id` (FK -> users)
- `name`, `type` (`STR|SIP|Serkom|Identitas|...`), `valid_until?` (date)
- `status` (`valid|expired|pending`), `file_url`, `created_at`, `updated_at`, `deleted_at?`
Index: `user_id`, `type`, `status`

### uploads
- `id` (PK), `uploader_id?` (FK -> users)
- `filename`, `original_name`, `mime_type`, `size`, `url`, `created_at`

### menus
- `id` (PK), `label`, `slug`, `to`, `icon?`, `parent_id?` (self FK)
- `position` (`header|sidebar|footer`), `order`, `is_active`, `is_fixed`, `roles` (JSON array)
- `created_at`, `updated_at`
Index: `position`, `parent_id`, `order`

### dynamic_contents
- `id` (PK), `slug` (unique), `title`, `description?`, `body?` (markdown), `html?`
- `date?`, `image?` (JSON), `authors?` (JSON array), `badge?` (JSON)
- `created_at`, `updated_at`

### homepage
- `id` (PK, single row), `hero` (JSON), `stats` (JSON array), `features` (JSON array), `seo` (JSON), `updated_at`

## Daftar Endpoint & Kontrak
Di bawah ini adalah endpoint yang diperlukan. Endpoint yang sudah ada di frontend ditandai (sudah ada).

### Auth
- `POST /api/auth/login`
  - Body: `{ email, password }`
  - Res: `{ success, data: { accessToken, refreshToken, user }, message }`
- `POST /api/auth/register`
  - Body: `{ name, email, password, phone?, category?, branchId? }`
  - Res: `{ success, data: { user }, message }`
- `GET /api/auth/me`
  - Header: `Authorization: Bearer <token>`
  - Res: `{ success, data: user }`
- `POST /api/auth/logout`
- `POST /api/auth/refresh`

### Users (Admin)
- `GET /api/users?role=&search=&page=&limit=`
- `GET /api/users/:id`
- `POST /api/users`
- `PUT /api/users/:id`
- `PATCH /api/users/:id` (ubah role/status)
- `DELETE /api/users/:id` (soft delete)

### Berita (sebagian sudah ada)
- (sudah ada) `GET /api/berita?category=&author=&month=YYYY-MM&sort=popular|&status=&search=&page=&limit=`
- (sudah ada) `GET /api/berita/:slug`
- `POST /api/berita` (admin)
  - Body minimal: `{ title, excerpt, content, image_url?, category, tags?: string[], author, status? }`
  - Aturan: generate `slug` unik; jika `status='published'` set `published_at` bila belum ada
- (sudah ada) `PUT /api/berita/:id` (update semua field + cek slug unik)
- (sudah ada) `PATCH /api/berita/:id` (ubah `status`, `published_at?`, `deleted_at?`)
- (sudah ada) `DELETE /api/berita/:id` (soft delete: set `deleted_at`)
- `GET /api/berita/categories` (list kategori unik)
- `GET /api/berita/tags` (list semua tag)

Contoh respons paginasi (GET /api/berita):
```json
{
  "success": true,
  "data": {
    "items": [ { "id": "1", "slug": "kongres-nasional-pdpi-2024", "title": "..." } ],
    "pagination": { "page": 1, "limit": 9, "total": 123, "totalPages": 14 }
  },
  "message": "Berita fetched successfully"
}
```

### Agenda (sebagian sudah ada)
- (sudah ada) `GET /api/agenda?upcoming=true|false&type=&status=&page=&limit=`
- (sudah ada) `GET /api/agenda/:slug`
- (sudah ada) `POST /api/agenda` (admin)
- (sudah ada) `PUT /api/agenda/:id`
- (sudah ada) `PATCH /api/agenda/:id`
- (sudah ada) `DELETE /api/agenda/:id`
- `GET /api/agenda/:id/registrations` (admin)
- `POST /api/agenda/:id/registrations` (member)
  - Body: `{ user_id }` atau infer dari token user
- `PATCH /api/agenda/:id/registrations/:regId` (admin; konfirmasi/batal)

### Direktori (sudah ada GET)
- (sudah ada) `GET /api/direktori?province=&type=&search=&page=&limit=`
- `GET /api/direktori/:id`
- `POST /api/direktori` (admin)
- `PUT /api/direktori/:id` (admin)
- `DELETE /api/direktori/:id` (admin)

### Pengurus (sudah ada GET)
- (sudah ada) `GET /api/pengurus?level=pusat|wilayah|cabang`
- `POST /api/pengurus` (admin)
- `PUT /api/pengurus/:id` (admin)
- `DELETE /api/pengurus/:id` (admin)

### Menu Navigasi (sebagian sudah ada)
- (sudah ada) `GET /api/menus` (gabungan semua menu)
- (sudah ada) `GET /api/menus/:position` (`header|sidebar|footer`)
- (sudah ada) `POST /api/menus` (admin; payload `{ position, menus: MenuItem[] }`)
  - Aturan: item `isFixed: true` tidak boleh dihapus/diubah `slug/to`; validasi hirarki (`parent_id`/children)

### Profil Organisasi (sudah ada GET)
- (sudah ada) `GET /api/profil`
- `POST /api/profil` (admin) — update `visiMisi`, `sejarah`, `adArt`

### Homepage (sebagian sudah ada)
- Tambah `GET /api/homepage` — ambil data beranda
- (sudah ada) `POST /api/homepage` (admin) — update partial `hero|stats|features|seo`

### Konten Dinamis
- Alias nama seragam: gunakan `content`
- `GET /api/content/:slug` — ambil konten berdasarkan slug (nested path didukung)
- `POST /api/content` (admin) — upsert berdasarkan `slug`
- `PUT /api/content/:id` (admin)
- `DELETE /api/content/:id` (admin)

### Dokumen Anggota & Upload (sudah ada upload)
- (sudah ada) `POST /api/upload` — multipart file; validasi tipe: `image/jpeg|jpg|png|webp|gif`; max size 5MB; simpan ke `public/uploads/YYYY/MM/`
- `GET /api/uploads` (admin) — daftar file yang pernah diupload
- `DELETE /api/uploads/:id` (admin) — hapus file (opsional: soft delete)
- `GET /api/documents` (member/admin)
  - Member: hanya dokumen miliknya; Admin: bisa filter `user_id`
  - Filter: `type`, `status`, `page`, `limit`
- `GET /api/documents/:id`
- `POST /api/documents`
  - Body: `{ user_id?, name, type, valid_until?, status, file_url }`
  - Jika member, `user_id` diambil dari token; jika admin, bisa set untuk user lain
- `PUT /api/documents/:id`
- `DELETE /api/documents/:id`

## Validasi & Normalisasi
- Slug: `toLowerCase()`, trim, hapus non `[a-z0-9\s-]`, ubah spasi ke `-`, dedup `-`, cek unik di resource yang sama terhadap entri non-deleted
- Status: hanya enum yang diizinkan (lihat masing-masing resource)
- Paginasi: angka positif; `limit` dengan batas maksimum wajar (mis. 100)
- Upload: tipe dan ukuran sesuai aturan; sanitize nama file saat simpan

## Contoh Header & RBAC
- Autentikasi:
  - `Authorization: Bearer <access_token>`
- RBAC (opsional header tambahan jika diperlukan logging/trace):
  - `X-User-Role: member|admin_cabang|admin_wilayah|admin_pusat`

## Contoh Payload & Respons

### POST /api/berita (admin)
Request:
```json
{
  "title": "Workshop Bronkoskopi Tingkat Lanjut",
  "excerpt": "Workshop hands-on bronkoskopi...",
  "content": "<p>Detail konten...</p>",
  "image_url": "/uploads/2025/01/1738130000_hash_img.jpg",
  "category": "kegiatan",
  "tags": ["workshop", "bronkoskopi"],
  "author": "Komite Pendidikan PDPI",
  "status": "draft"
}
```
Response:
```json
{ "success": true, "data": { "id": "...", "slug": "workshop-bronkoskopi-tingkat-lanjut" }, "message": "Berita berhasil dibuat" }
```

### GET /api/berita
Response (paginasi):
```json
{
  "success": true,
  "data": {
    "items": [ { "id": "1", "slug": "kongres-nasional-pdpi-2024", "title": "..." } ],
    "pagination": { "page": 1, "limit": 9, "total": 123, "totalPages": 14 }
  },
  "message": "Berita fetched successfully"
}
```

### PUT /api/agenda/:id (admin)
- Aturan: cek unik slug jika berubah; jangan ubah entri yang `deleted_at` terisi
Response:
```json
{ "success": true, "data": { "id": "...", "slug": "webinar-asma-jan2025", "updated_at": "..." } }
```

### POST /api/upload
Response:
```json
{
  "success": true,
  "data": {
    "url": "/uploads/2025/01/1738130000_hash_img.jpg",
    "filename": "1738130000_hash_img.jpg",
    "originalName": "img.jpg",
    "size": 123456,
    "type": "image/jpeg"
  },
  "message": "File uploaded successfully"
}
```

## Variabel Lingkungan (disarankan)
- `JWT_SECRET`, `JWT_EXPIRES_IN`, `REFRESH_EXPIRES_IN`
- `DB_URL` (PostgreSQL/MySQL)
- `APP_URL` (untuk membangun URL publik upload, jika perlu)
- `UPLOAD_DIR` (default `public/uploads`)

## Catatan Implementasi
- Ikuti format respons dan aturan filter/paginasi seperti yang digunakan di frontend saat ini
- Terapkan soft delete secara konsisten; GET default mengecualikan `deleted_at` != null
- Tambahkan index DB sesuai kolom filter untuk performa
- Pastikan RBAC memproteksi `POST/PUT/PATCH/DELETE` sesuai role; `GET` publik tidak memerlukan token
- Untuk `agenda_registrations`, validasi kapasitas (`quota`) dan cegah duplikasi per user
- Untuk `menus`, validasi item `isFixed: true` tidak boleh dihapus atau diubah field kunci

---
Dokumen ini bertujuan menjadi kontrak kerja awal backend. Jika ada preferensi teknologi (mis. ORM Prisma/TypeORM, DB PostgreSQL/MySQL), sertakan dalam implementasi dan sesuaikan detail kecil seperti tipe tanggal, penamaan kolom, serta mekanisme error handling.