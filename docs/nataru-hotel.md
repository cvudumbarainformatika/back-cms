# Nataru Hotel - Create Endpoint

Endpoint: POST /api/v1/nataru-hotel/create
Auth: JWT (protected route). idUser diambil dari JWT.

Perilaku umum
- Upsert: jika key upsert cocok, data di-update; jika tidak, dibuat record baru.
- Response: 200 OK untuk update, 201 Created untuk create baru.

Event: Idul Fitri
- tanggal: opsional; jika ada, diabaikan untuk upsert.
- alias: wajib.
- idHotel: wajib (camelCase: idHotel; integer). Kirim angka, bukan objek.
- Key upsert: (event, alias, periode, idHotel).

Contoh payload (Idul Fitri)
{
  "event": "Idul Fitri",
  "alias": "H+3",
  "periode": 2025,
  "idHotel": 32,
  "jml_kmr": 100,
  "jml_tt": 150,
  "jml_karyawan": 40,
  "wisman": 10,
  "wisnus": 60,
  "jml_kmr_tersedia": 90,
  "kmr_terjual_wisman": 5,
  "kmr_terjual_wisnus": 20,
  "jml_tt_tersedia": 120,
  "tt_terjual_wisman": 4,
  "tt_terjual_wisnus": 16
}

Event: selain "Idul Fitri" (mis. "Nataru")
- tanggal: wajib.
- idHotel: wajib (integer).
- Key upsert: (tanggal, event, periode, idHotel).

Contoh payload (Nataru)
{
  "event": "Nataru",
  "tanggal": "2025-12-27",
  "alias": "H+1",
  "periode": 2025,
  "idHotel": 32,
  "jml_kmr": 120,
  "jml_tt": 200,
  "jml_karyawan": 50,
  "wisman": 30,
  "wisnus": 140,
  "jml_kmr_tersedia": 100,
  "kmr_terjual_wisman": 10,
  "kmr_terjual_wisnus": 30,
  "jml_tt_tersedia": 160,
  "tt_terjual_wisman": 8,
  "tt_terjual_wisnus": 25
}

Catatan
- Field yang tidak dikirim saat update tidak akan di-null-kan.
- Get-list memiliki pilihan filter periode, event, dan id_hotel serta mengembalikan nameHotel via LEFT JOIN.
- Disarankan unique index: (event, alias, periode, idHotel) dan (tanggal, event, periode, idHotel) pada tabel nataru_hotel.
