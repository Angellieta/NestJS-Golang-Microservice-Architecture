# Arsitektur Microservice dengan NestJS & Golang

Ini adalah project aplikasi microservice yang dibangun berdasarkan tantangan _fullstack developer_, yang mendemonstrasikan arsitektur _event-driven_, _caching_, dan komunikasi antar service.

## Tujuan Proyek

Proyek ini bertujuan untuk mendemonstrasikan implementasi arsitektur microservice modern yang tangguh dan skalabel. Fokus utamanya adalah untuk menunjukkan:

- **Komunikasi Asynchronous**: Komunikasi antar service yang terpisah (_decoupled_) menggunakan _message broker_.
- **Arsitektur Bersih (Clean Architecture)**: Pemisahan yang jelas antara lapisan logika, data, dan presentasi.
- **Caching**: Pemanfaatan _caching layer_ untuk meningkatkan performa dan mengurangi beban database.
- **Poliglot**: Penggunaan dua bahasa pemrograman berbeda (NestJS & Go) dalam satu sistem yang terintegrasi.

## Arsitektur

Aplikasi ini terdiri dari dua microservice utama:

- **Product Service** (`product-service`): Dibangun dengan **NestJS (TypeScript)**. Bertanggung jawab untuk mengelola data produk (CRUD), stok, dan mendengarkan event dari `order-service`.
- **Order Service** (`order-service`): Dibangun dengan **Golang**. Bertanggung jawab untuk mengelola pesanan dan memvalidasi keberadaan produk sebelum membuat pesanan.

Kedua service berkomunikasi secara asynchronous menggunakan **RabbitMQ** sebagai _message broker_. **Redis** digunakan sebagai _caching layer_ untuk mempercepat pengambilan data, dan **PostgreSQL** sebagai database utama. Semua layanan dijalankan dalam container menggunakan **Docker Compose**.

## Teknologi yang Digunakan

- **Backend**: NestJS (TypeScript), Golang
- **Database**: PostgreSQL
- **Caching**: Redis
- **Message Broker**: RabbitMQ
- **Containerization**: Docker & Docker Compose
- **Load Testing**: k6

## Cara Menjalankan Project Secara Lokal

### Prasyarat

Download Docker Desktop (Windows 10/11)

- [Docker](https://www.docker.com/products/docker-desktop/)
- [Docker Compose](https://docs.docker.com/compose/install/) (biasanya sudah termasuk dalam Docker Desktop)

### Langkah-langkah Menjalankan

1.  Clone repository
2.  Buka terminal di direktori utama project.
3.  Jalankan perintah berikut untuk membangun dan memulai semua service:
    ```bash
    docker-compose up --build
    ```
4.  Aplikasi akan siap dalam beberapa saat.
    - `product-service` berjalan di `http://localhost:3000`
    - `order-service` berjalan di `http://localhost:8080`
    - `RabbitMQ Management UI` berjalan di `http://localhost:15672` (login: guest/guest)

## Contoh Request API (PowerShell)

### 1. Membuat Produk Baru

```powershell
Invoke-WebRequest -Uri http://localhost:3000/products -Method POST -ContentType "application/json" -Body '{"name": "Laptop Pro", "price": 2000, "qty": 25}'
```

### 2. Mengambil Produk Berdasarkan ID (Cached)

Ganti <PRODUCT_ID> dengan ID dari hasil pembuatan produk

```powershell
Invoke-WebRequest -Uri http://localhost:3000/products/<PRODUCT_ID>
```

### 3. Membuat Order Baru (dengan Validasi Produk)

Ganti <PRODUCT_ID> dengan ID yang valid

```powershell
Invoke-WebRequest -Uri http://localhost:8080/orders -Method POST -ContentType "application/json" -Body '{"productId": "<PRODUCT_ID>", "price": 2000, "qty": 3}'
```

### 4. Mengambil Order Berdasarkan Product ID (Cached)

Ganti <PRODUCT_ID> dengan ID yang valid

```powershell
Invoke-WebRequest -Uri http://localhost:8080/orders/product/<PRODUCT_ID>
```

## Alur Kerja Event-Driven

1. Client membuat order baru melalui _POST /orders_ di _order-service_.
2. _order-service_ melakukan panggilan API sinkron ke _product-service_ untuk memvalidasi productId.
3. Jika valid, _order-service_ menyimpan order ke PostgreSQL dan menerbitkan event _order.created_ ke RabbitMQ.
4. _product-service_ mendengarkan event _order.created_, menerima pesan, dan mengurangi stok produk yang sesuai di databasenya.
