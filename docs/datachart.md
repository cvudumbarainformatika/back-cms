# Query SQL untuk Chart di Sistem Informasi Juwita Malam

File ini berisi semua query SQL yang digunakan untuk mengambil data untuk pembuatan chart di sistem informasi Juwita Malam (admin/index.php).

## Query Pertama: Wisatawan Mancanegara (Wisman) - Pie Chart

### Query untuk Pie Chart Wisatawan Mancanegara (qWisman):

```sql
SELECT dh.tahun, n.namaNegara AS negara, SUM(dh.total) AS total
FROM data_hotel dh
JOIN hotels h ON dh.idHotel = h.idHotel
JOIN negara n ON dh.negaraProvinsi = n.idNegara
WHERE dh.tahun = '$tahun' AND dh.kategori = 'Wisman' AND dh.status = '1'
GROUP BY dh.tahun, n.namaNegara
ORDER BY dh.tahun ASC, n.namaNegara
```

## Query Kedua: Wisatawan Nusantara (Wisnus) - Pie Chart

### Query untuk Pie Chart Wisatawan Nusantara (qWisnus):

```sql
SELECT dh.tahun, p.namaProvinsi AS provinsi, SUM(dh.total) AS total
FROM data_hotel dh
JOIN hotels h ON dh.idHotel = h.idHotel
JOIN provinsi p ON dh.negaraProvinsi = p.idProvinsi
WHERE dh.tahun = '$tahun' AND dh.kategori = 'Wisnus' AND dh.status = '1'
GROUP BY dh.tahun, p.namaProvinsi
ORDER BY dh.tahun ASC, p.namaProvinsi
```

## Query Ketiga: Statistik Kunjungan Hotel per Bulan (Wisman & Wisnus)

Menggabungkan pengambilan data Wisman dan Wisnus per bulan menjadi satu query untuk efisiensi.

### Query untuk Bar Chart Kunjungan Hotel - Admin:

```sql
SELECT
    bulan,
    SUM(CASE WHEN kategori = 'Wisman' THEN total ELSE 0 END) as wisman,
    SUM(CASE WHEN kategori = 'Wisnus' THEN total ELSE 0 END) as wisnus
FROM data_hotel
WHERE tahun = '$tahun' AND status = '1'
GROUP BY bulan
```

### Query untuk Bar Chart Kunjungan Hotel - Hotel (Specific User):

```sql
SELECT
    bulan,
    SUM(CASE WHEN kategori = 'Wisman' THEN total ELSE 0 END) as wisman,
    SUM(CASE WHEN kategori = 'Wisnus' THEN total ELSE 0 END) as wisnus
FROM data_hotel
WHERE idHotel = '$login' AND tahun = '$tahun' AND status = '1'
GROUP BY bulan
```

## Query Keempat: Statistik Kunjungan Destinasi Wisata per Bulan

Menggabungkan pengambilan data Wisman dan Wisnus di destinasi wisata.

### Query untuk Bar Chart Destinasi Wisata - Admin:

```sql
SELECT
    bulan,
    SUM(wisman) as wisman,
    SUM(wisnus) as wisnus
FROM data_wisata
WHERE tahun = '$tahun'
GROUP BY bulan
```

### Query untuk Bar Chart Destinasi Wisata - Wisata (Specific User):

```sql
SELECT
    bulan,
    SUM(wisman) as wisman,
    SUM(wisnus) as wisnus
FROM data_wisata
WHERE idWisata = '$login' AND tahun = '$tahun'
GROUP BY bulan
```

## Ringkasan Tabel yang Digunakan

1. **data_hotel**: Tabel utama untuk data kunjungan wisatawan ke hotel

   - Kolom penting: tahun, kategori, idHotel, negaraProvinsi, total, status, bulan

2. **data_wisata**: Tabel untuk data kunjungan wisatawan ke destinasi wisata

   - Kolom penting: tahun, wisman, wisnus, idWisata, bulan

3. **hotels**: Tabel informasi hotel

   - Kolom penting: idHotel, nameHotel

4. **negara**: Tabel informasi negara

   - Kolom penting: idNegara, namaNegara

5. **provinsi**: Tabel informasi provinsi

   - Kolom penting: idProvinsi, namaProvinsi

6. **wisata**: Tabel informasi destinasi wisata
   - Kolom penting: idWisata, nameWisata

## Variabel Dinamis dalam Query

- `$tahun`: Tahun yang sedang dipilih/diset di sesi
- `$login`: ID hotel atau destinasi wisata berdasarkan login pengguna
- `$akses`: Hak akses pengguna (Administrator, Hotel, Wisata)

## Format Respons API dari Backend

Untuk mengimplementasikan chart-chart di atas di frontend Vue, berikut adalah format respons API yang diharapkan dari backend Go:

### 1. Endpoint: GET /api/v1/charts/wisman-by-country

Respons untuk pie chart wisatawan mancanegara berdasarkan negara asal:

```json
{
  "success": true,
  "message": "Data wisatawan mancanegara berhasil diambil",
  "data": {
    "wismanByCountry": [
      { "country": "Indonesia", "value": 12345 },
      { "country": "Malaysia", "value": 9876 },
      { "country": "Singapura", "value": 8765 },
      { "country": "Australia", "value": 7654 }
    ],
    "year": 2024
  }
}
```

### 2. Endpoint: GET /api/v1/charts/wisnus-by-province

Respons untuk pie chart wisatawan nusantara berdasarkan provinsi asal:

```json
{
  "success": true,
  "message": "Data wisatawan nusantara berdasarkan provinsi berhasil diambil",
  "data": {
    "wisnusByProvince": [
      { "province": "Jawa Timur", "value": 25000 },
      { "province": "Jawa Tengah", "value": 18000 },
      { "province": "Jawa Barat", "value": 15000 },
      { "province": "DKI Jakarta", "value": 12000 },
      { "province": "Bali", "value": 10000 }
    ],
    "year": 2024
  }
}
```

### 3. Endpoint: GET /api/v1/charts/hotel-visitors-by-month

**[OPTIMIZED]** Menggabungkan data Wisman dan Wisnus per bulan dalam satu respons.
Respons untuk bar chart kunjungan hotel:

```json
{
  "success": true,
  "message": "Data kunjungan hotel per bulan berhasil diambil",
  "data": {
    "visitorsByMonth": [
      { "month": "Januari", "wisman": 1000, "wisnus": 5000 },
      { "month": "Februari", "wisman": 1200, "wisnus": 5500 },
      { "month": "Maret", "wisman": 1500, "wisnus": 6000 },
      { "month": "April", "wisman": 1800, "wisnus": 6500 },
      { "month": "Mei", "wisman": 2000, "wisnus": 7000 },
      { "month": "Juni", "wisman": 2200, "wisnus": 7500 },
      { "month": "Juli", "wisman": 2500, "wisnus": 8000 },
      { "month": "Agustus", "wisman": 2800, "wisnus": 8500 },
      { "month": "September", "wisman": 3000, "wisnus": 9000 },
      { "month": "Oktober", "wisman": 2700, "wisnus": 8800 },
      { "month": "November", "wisman": 2400, "wisnus": 8200 },
      { "month": "Desember", "wisman": 2100, "wisnus": 7800 }
    ],
    "year": 2024
  }
}
```

### 4. Endpoint: GET /api/v1/charts/tourism-visitors-by-month

**[OPTIMIZED]** Menggabungkan data Wisman dan Wisnus di destinasi wisata per bulan.
Respons untuk chart kunjungan destinasi wisata:

```json
{
  "success": true,
  "message": "Data kunjungan destinasi wisata berhasil diambil",
  "data": {
    "visitorsByMonth": [
      { "month": "Januari", "wisman": 800, "wisnus": 4200 },
      { "month": "Februari", "wisman": 950, "wisnus": 4600 },
      { "month": "Maret", "wisman": 1200, "wisnus": 5000 },
      { "month": "April", "wisman": 1400, "wisnus": 5500 },
      { "month": "Mei", "wisman": 1600, "wisnus": 6000 },
      { "month": "Juni", "wisman": 1800, "wisnus": 6200 },
      { "month": "Juli", "wisman": 2000, "wisnus": 6800 },
      { "month": "Agustus", "wisman": 2200, "wisnus": 7200 },
      { "month": "September", "wisman": 2500, "wisnus": 7500 },
      { "month": "Oktober", "wisman": 2300, "wisnus": 7300 },
      { "month": "November", "wisman": 2100, "wisnus": 6900 },
      { "month": "Desember", "wisman": 1900, "wisnus": 6500 }
    ],
    "year": 2024
  }
}
```

### Parameter Umum untuk Semua Endpoint

- `year`: Tahun untuk filter data (opsional, default: tahun saat ini)
- `access_token`: Token otentikasi untuk akses (diperlukan untuk endpoint yang dilindungi)

### Error Response Umum

```json
{
  "success": false,
  "message": "Deskripsi kesalahan",
  "data": null
}
```
