# Nataru Kereta Api - Create Endpoint

Endpoint: POST /api/v1/nataru-kereta-api/create
Auth: JWT (protected route). idUser diambil dari JWT.

Perilaku umum
- Upsert: jika key upsert cocok, data di-update; jika tidak, dibuat record baru.
- Response: 200 OK untuk update, 201 Created untuk create baru.

Event: Idul Fitri
- tanggal: opsional; jika ada, diabaikan untuk upsert.
- alias: wajib.
- Key upsert: (event, alias, periode).

Contoh payload (Idul Fitri)
{
  "event": "Idul Fitri",
  "alias": "H+1",
  "periode": 2025,
  "naik": 1200,
  "turun": 1100
}

Event: selain "Idul Fitri" (mis. "Nataru")
- tanggal: wajib.
- Key upsert: (tanggal, event, periode).

Contoh payload (Nataru)
{
  "event": "Nataru",
  "tanggal": "2025-12-27",
  "alias": "H+1",
  "periode": 2025,
  "naik": 1000,
  "turun": 900
}

Catatan
- Field yang tidak dikirim saat update tidak akan di-null-kan (hanya overwrite field yang dikirim).
- Disarankan unique index: (event, alias, periode) dan (tanggal, event, periode) pada tabel nataru_keretaapi.
