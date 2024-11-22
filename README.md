
# Expense Tracker API

## Persyaratan
- [Go](https://go.dev/) versi 1.23.3 atau lebih baru
- [Git](https://git-scm.com/) (untuk mengelola repositori)
- Database (contoh: PostgreSQL)
- Environment variables (lihat file `.env.example`)

## Instalasi
Ikuti langkah-langkah berikut untuk menjalankan proyek ini:

1. Clone repositori ini:
   ```bash
   git clone <URL_REPOSITORI>
   cd <NAMA_FOLDER_PROYEK>
   ```

2. Rename file `.env.example` menjadi `.env`:
   ```bash
   mv .env.example .env
   ```

3. Edit file `.env` dan tambahkan nilai konfigurasi:
   ```plaintext
   JWT_SECRET=           # Secret untuk JWT
   DB_HOST=              # Host database (contoh: localhost)
   DB_PORT=              # Port database (contoh: 5432)
   DB_USER=              # Username database
   DB_PASSWORD=          # Password database
   DB_NAME=              # Nama database
   DB_SSLMODE=           # Mode SSL (contoh: disable)
   FRONTEND_URL=         # URL frontend (contoh: http://localhost:3000)
   ```

4. Instal dependensi yang dibutuhkan:
   ```bash
   go mod tidy
   ```

5. Jalankan aplikasi:
   ```bash
   go run main.go
   ```
