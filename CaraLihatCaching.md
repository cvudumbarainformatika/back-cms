# Cara Melihat Isi Caching Redis

Untuk melihat isi caching yang digunakan oleh sistem Juwita Tourism Information System, ikuti langkah-langkah berikut:

## 1. Menggunakan CLI Redis

Jika Redis berjalan di lokal mesin Anda, Anda bisa menggunakan CLI Redis untuk melihat isi cache:

### Akses Redis CLI
```bash
redis-cli
```

### Lihat semua key yang tersimpan
```bash
KEYS *
```

### Lihat key yang terkait dengan master data
```bash
KEYS master_*
```

### Lihat isi dari key tertentu
```bash
GET master_wisata_list
GET master_hotel_list
GET master_bulan_list
GET master_tahun_list
```

### Periksa TTL (Time To Live) dari key tertentu
```bash
TTL master_wisata_list
```

## 2. Menggunakan Docker (jika Redis berjalan di container)

### Temukan container Redis
```bash
docker ps | grep redis
```

### Jalankan Redis CLI di dalam container
```bash
docker exec -it <container_name> redis-cli
```

Lalu ikuti langkah-langkah di atas.

## 3. Melalui Aplikasi Web (opsional)

Jika Anda menggunakan tools seperti Redis Commander:

### Jalankan Redis Commander (jika belum)
```bash
docker run -d --name redis-commander \
  -p 8081:8081 \
  --env REDIS_HOSTS=local:redis:6379 \
  rediscommander/redis-commander:latest
```

Buka browser di `http://localhost:8081` untuk antarmuka web Redis.

## 4. Key Caching yang Digunakan di Sistem

Berikut adalah daftar key caching yang digunakan di sistem Juwita:

- `master_wisata_list`: Cache untuk endpoint `/api/v1/master/wisata`
- `master_hotel_list`: Cache untuk endpoint `/api/v1/master/hotels`
- `master_bulan_list`: Cache untuk endpoint `/api/v1/master/bulan`
- `master_tahun_list`: Cache untuk endpoint `/api/v1/master/tahun`

Semua cache ini memiliki TTL (Time To Live) 600 detik (10 menit).

## 5. Melihat Cache dalam Log Aplikasi

Anda juga bisa melihat penggunaan cache dari log aplikasi saat request masuk:
- Jika data diambil dari cache, akan langsung dikembalikan tanpa query ke database
- Jika cache kosong/expire, sistem akan mengambil dari database dan menyimpan ke cache

## 6. Membersihkan Cache

Untuk membersihkan semua cache:
```bash
FLUSHALL
```

Atau membersihkan key tertentu:
```bash
DEL master_wisata_list
```

## 7. Troubleshooting

Jika cache tidak berfungsi:
- Pastikan Redis server berjalan
- Periksa koneksi aplikasi ke Redis
- Verifikasi konfigurasi Redis di file `.env`
- Periksa apakah parameter query digunakan (karena caching dilewati saat ada parameter)