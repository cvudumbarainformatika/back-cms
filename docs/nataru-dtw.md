# Nataru DTW - Create Endpoint

Endpoint: POST /api/v1/nataru-dtw/create
Auth: JWT (protected route). idUser diambil dari JWT.

Perilaku umum
- Upsert: jika key upsert cocok, data di-update; jika tidak, dibuat record baru.
- Response: 200 OK untuk update, 201 Created untuk create baru.

Event: Idul Fitri
- tanggal: opsional; jika ada, diabaikan untuk upsert.
- alias: wajib.
- idWisata: wajib (camelCase: idWisata; integer).
- Key upsert: (event, alias, periode, idWisata).

Contoh payload (Idul Fitri)
{
  "event": "Idul Fitri",
  "alias": "H+2",
  "periode": 2025,
  "idWisata": 45,
  "wisman": 30,
  "wisnus": 120
}

Event: selain "Idul Fitri" (mis. "Nataru")
- tanggal: wajib.
- idWisata: wajib.
- Key upsert: (tanggal, event, periode, idWisata).

Contoh payload (Nataru)
{
  "event": "Nataru",
  "tanggal": "2025-12-27",
  "alias": "H+1",
  "periode": 2025,
  "idWisata": 45,
  "wisman": 10,
  "wisnus": 200
}

Catatan
- Field yang tidak dikirim saat update tidak akan di-null-kan.
- Get-list memiliki pilihan filter periode, event, dan id_wisata serta mengembalikan nameWisata via LEFT JOIN.
- Disarankan unique index: (event, alias, periode, idWisata) dan (tanggal, event, periode, idWisata) pada tabel nataru_dtw.
