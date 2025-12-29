# Alur Kerja Deploy dan Update Aplikasi di Produksi

Dokumen ini menjelaskan langkah-langkah untuk melakukan deployment awal dan pembaruan (update) kode aplikasi di lingkungan server produksi.

## Prasyarat

- Server produksi sudah terinstall `git` dan `docker-compose`.
- Akses ke repository Git proyek.

## 1. Deployment Awal (Initial Setup)

Langkah-langkah ini hanya dilakukan sekali saat pertama kali men-deploy aplikasi di server baru.

1.  **Clone Repository**
    Clone kode sumber dari repository Git Anda.
    ```bash
    git clone [URL-repo-anda]
    cd [nama-folder-proyek]
    ```

2.  **Buat File Environment**
    Salin file contoh `.env.production.example` (jika ada) atau buat file `.env.production` secara manual.
    ```bash
    cp .env.example .env.production
    ```

3.  **Isi Konfigurasi Environment**
    Buka dan edit file `.env.production` untuk mengisi semua nilai yang dibutuhkan, seperti password database, JWT secret, dan lain-lain.
    ```bash
    nano .env.production
    ```

4.  **Build dan Jalankan Layanan**
    Gunakan `docker-compose` untuk build image dan menjalankan semua layanan (aplikasi, database, redis) di background (`-d`). Flag `--env-file` memastikan semua variabel dimuat dengan benar.
    ```bash
    docker-compose -f docker-compose.prod.yml --env-file .env.production up -d --build
    ```
    Flag `--build` memastikan image aplikasi dibangun saat pertama kali.

## 2. Pembaruan Kode (Update)

Lakukan langkah-langkah berikut setiap kali ada perubahan kode yang ingin di-deploy ke server.

1.  **Tarik Perubahan Terbaru**
    Masuk ke direktori proyek di server dan lakukan `git pull` untuk mendapatkan kode terbaru.
    ```bash
    cd /path/to/[nama-folder-proyek]
    git pull
    ```

2.  **Build Ulang dan Restart Layanan**
    Jalankan kembali perintah `up` dengan flag `--build` dan `--env-file`. Docker Compose akan cukup pintar untuk:
    - Membangun ulang image aplikasi dengan kode terbaru.
    - Menggunakan *cache* sehingga proses build tetap cepat jika dependensi tidak berubah.
    - Secara otomatis mematikan kontainer lama dan menjalankan yang baru dengan image yang sudah diperbarui.
    ```bash
    docker-compose -f docker-compose.prod.yml --env-file .env.production up -d --build
    ```

## 3. Perintah Berguna Lainnya

-   **Melihat Log Aplikasi:**
    ```bash
    docker-compose -f docker-compose.prod.yml --env-file .env.production logs -f app
    ```
-   **Menghentikan Semua Layanan:**
    ```bash
    docker-compose -f docker-compose.prod.yml --env-file .env.production down
    ```
-   **Membersihkan Image yang Tidak Terpakai (Opsional):**
    Setelah beberapa kali update, mungkin ada image-image lama yang tidak terpakai. Anda bisa membersihkannya dengan perintah ini.
    ```bash
    docker image prune
    ```
