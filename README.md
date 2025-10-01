# Arsitektur Microservice dengan Nest JS dan Golang

Project ini adalah aplikasi microservice yang dibangun berdasarkan tantangan sebagai seorang fullstack developer dengan mendemonstrasikan arsitektur event-driven, caching, dan komunikasi antar service.

## 1. Ikhtisar Proyek

Proyek ini bertujuan untuk mengimplementasikan arsitektur microservice modern yang tangguh dan skalabel. Fokus utamanya adalah untuk menunjukkan:

- Komunikasi Asynchronous: Komunikasi antar service yang terpisah (decoupled) menggunakan message broker.
- Arsitektur Bersih (Clean Architecture): Pemisahan yang jelas antara lapisan logika, data, dan presentasi.
- Caching: Pemanfaatan caching layer untuk meningkatkan performa dan mengurangi beban database.
- Poliglot: Penggunaan dua bahasa pemrograman berbeda (NestJS & Go) dalam satu sistem yang terintegrasi.
- Backend for Frontend (BFF): Sebuah _single entrypoint_ (`localhost:5000`) yang menyederhanakan interaksi dari sisi klien dan menggabungkan data dari berbagai layanan.
- Validasi & Error Handling: _Middleware_ validasi otomatis di setiap layanan untuk memastikan integritas data yang masuk.
- Pelacakan Terdistribusi: Implementasi **Correlation ID** (`x-correlation-id`) yang mengalir melalui semua layanan untuk kemudahan _debugging_ dan pelacakan.

### Teknologi yang Digunakan

- Backend: NestJS (TypeScript), Golang, NodeJS (Express)
- Database: PostgreSQL
- Caching: Redis
- Message Broker: RabbitMQ
- Containerization: Docker & Docker Compose
- Load Testing: Jest (NestJS), Go Testing (Go), k6 (Performance)

## 2. Struktur Proyek

Proyek ini menggunakan pendekatan arsitektur berlapis (layered architecture) di dalam setiap microservice untuk memastikan pemisahan tanggung jawab (separation of concerns) yang jelas:

```
.
├── product-service
│     └── src
│            ├── products
│            │     ├── dto
│            │     │     └── create-product.dto.ts
│            │     ├── entities
│            │     │     └── product.entity.ts
│            │     ├── products.controller.spec.ts
│            │     ├── products.controller.ts
│            │     ├── products.service.spec.ts
│            │     ├── products.service.ts
│            │     └── products.module.ts
│            ├── app.controller.ts
│            ├── app.module.ts
│            ├── app.service.ts
│            ├── main.ts
│            ├── Dockerfile
│            └── . . .
├── order-service
│     ├── cmd
│     │     └── server
│     │            └── main.go
│     ├── internal
│     │     ├── handler
│     │     │     └── http
│     │     │            └── order_handler.go
│     │     ├── model
│     │     │     └── order.go
│     │     ├── repository
│     │     │     └── order_repository.go
│     │     └── service
│     │            ├── order_service_test.go
│     │            └── order_service.go
│     ├── pkg
│     │     ├── database
│     │     │     └── postgres.go
│     │     ├── rabbitmq
│     │     │     └── publisher.go
│     │     └── redis
│     │            └── client.go
│     ├── Dockerfile
│     ├── go.mod
│     └── go.sum
├── bff-service
│     ├── index.js
│     ├── Dockerfile
│     └── . . .
├── bff-service
│     ├── index.js
│     ├── Dockerfile
│     └── . . .
├── performance-tests
│     └── create-order.js
└── docker-compose.yml
```

### Arsitektur berlapis (Layered Architecture)

- `product-service`: Folder microservice yang bertanggung jawab atas semua hal yang berkaitan dengan produk (NestJS)

  - `src`: Berisi kode sumber aplikasi dengan arsitektur berlapis (layered architecture) terstruktur
    - `products`: Sebuah "modul" NestJS yang mengelompokkan semua logika terkait fitur produk.
      - `dto`: Folder yang mendefinisikan "bentuk" data yang diharapkan dari body request saat membuat produk baru yang berguna untuk validasi dan type safety.
      - `entities`: Folder yang merepresentasikan tabel `products` di database serta menghubungkan properti kelas ke kolom tabel.
      - `products.controller.ts`: Pintu gerbang untuk Menangani routing HTTP dan menerima `Correlation ID` dari header.
      - `products.controller.ts`: Pintu gerbang untuk Menangani routing HTTP dan menerima `Correlation ID` dari header.
      - `products.controller.spec.ts`: FIle unit test untuk `ProductsController`
      - `products.service.ts`: Berisi otak atau logika bisnis serta fungsi untuk mengurangi stok yang dipicu oleh event RabbitMQ dan menggunakan Correlation ID dalam log.
      - `products.service.ts`: Berisi otak atau logika bisnis serta fungsi untuk mengurangi stok yang dipicu oleh event RabbitMQ dan menggunakan Correlation ID dalam log.
      - `products.service.spec.ts`: File unit test untuk `ProductsService`.
      - `products.module.ts`: Berfungsi sebagai penyatu antara `controller`, `service`, `entity` produk, serta mendaftarkannya ke aplikasi NestJS.
    - `app.module.ts`: Modul utama aplikasi yang mengimpor dan mengonfigurasi semua modul lain, termasuk `ProductsModule`, koneksi database (TypeOrmModule), koneksi RabbitMQ, dan koneksi Redis (CacheModule).
    - `main.ts`: Entrypoint yang membuat aplikasi NestJS dan memulai server web untuk mendengarkan request.
    - `Dockerfile`: Resep untuk membangun image Docker dari aplikasi NestJS agar bisa dijalankan di dalam container.

- `order-service`: Folder microservice yang bertanggung jawab atas semua hal yang berkaitan dengan pesanan (Golang)
  - `cmd/server/main.go`: Entrypoint yang melakukan semua inisialisasi koneksi ke database PostgreSQL, Redis, dan RabbitMQ. Selain itu, melakukan dependency injection (menyambungkan `repository`, `service`, dan `handler`) dan memulai server web Go.
  - `internal`: Berisi kode inti aplikasi (Golang)
    - `handler`: Folder yang menangani request HTTP, membaca `Correlation ID`, memvalidasi DTO, dan memanggil service.
    - `model`: Folder yang berisikan file `order.go` yang merupakan sebuah struct Go yang mendefinisikan struktur data order, termasuk `CorrelationID` untuk event payload.
    - `handler`: Folder yang menangani request HTTP, membaca `Correlation ID`, memvalidasi DTO, dan memanggil service.
    - `model`: Folder yang berisikan file `order.go` yang merupakan sebuah struct Go yang mendefinisikan struktur data order, termasuk `CorrelationID` untuk event payload.
    - `repository`: Folder yang berisikan file `order_repository.go` memiliki fungsi seperti `CreateOrder` dan `GetOrdersByProductID` yang bertugas berkomunikasi langsung dengan database.
    - `service`: Folder yang berisikan 2 file. `order_service.go` berisi logika bisnis seperti memanggil `product-service` untuk validasi `Correlation ID`, mengimplementasikan caching di Redis, dan memanggil publisher RabbitMQ. `order_service_test.go` merupakan unit test untuk `OrderService`.
    - `service`: Folder yang berisikan 2 file. `order_service.go` berisi logika bisnis seperti memanggil `product-service` untuk validasi `Correlation ID`, mengimplementasikan caching di Redis, dan memanggil publisher RabbitMQ. `order_service_test.go` merupakan unit test untuk `OrderService`.
  - `pkg`: Berisi kode paket yang dapat digunakan kembali.
    - `database`: Helper untuk membuat koneksi ke PostgreSQL.
    - `rabbitmq`: Helper untuk koneksi dan mengirim pesan ke RabbitMQ.
    - `redis`: Helper untuk membuat koneksi ke Redis.
  - `Dockerfile`: Resep untuk membangun aplikasi Go menjadi binary kecil ke dalam image Docker yang efisien.
  - `go.mod` & `go.sum` : Berisi dependensi atau library eksternal yang digunakan oleh project Go.
- `bff-service`: Folder yang bertindak sebagai lapisan untuk menangani permintaan dari client.

  - `index.js`: Kode server Express.js yang bertanggungjawab untuk membuat serta meneruskan `Correlation ID`, menerima semua permintaan masuk dari penggunaa, meneruskan proxy permintaan tersebut ke `product-service` atau `order-service`, menggabungkan data dari beberapa layanan (`order-summary`).
  - `Dockerfile`: Resep untuk membangun image Docker untuk layanan BFF.

- `performance-test`: Berisikan file `create-order.js` yang menampung skrip tes performa yang ditulis dalam JavaScript untuk dijalankan oleh k6.
- `docker-compose.yml`: File utama yang mendefinisikan semua layanan (`product-service`, `order-service`, `postgres`, `redis`, `rabbitmq`), mengonfigurasinya, dan memberitahu Docker bagaimana cara menjalankan semuanya secara bersamaan sebagai satu sistem yang utuh.
  </br>
  </br>

```
                                      +--------------------------------+
                                      |         USER / CLIENT          |
                                      +--------------------------------+
                                                   |
                                                   | (Request with/without X-Correlation-ID)
                                                   v
+--------------------------------------------------------------------------------------------------+
|                                    BFF SERVICE (Node.js/Express)                                 |
|                                         (Port: 5000)                                             |
| 1. Middleware:                                                                                   |
|    - Menerima Request.                                                                           |
|    - Membuat X-Correlation-ID jika tidak ada.                                                    |
| 2. Meneruskan Request (Proxy) dengan menyertakan X-Correlation-ID di header.                     |
+--------------------------------------------------------------------------------------------------+
             |                                                                 |
(API Call with X-Correlation-ID)                                  (API Call with X-Correlation-ID)
             |                                                                 |
             v                                                                 v
+--------------------------------+                              +------------------------------------+
|    ORDER SERVICE (Go)          |--(API Call for Validation)-->|      PRODUCT SERVICE (NestJS)      |
|  (Internal Port: 8080)         |   (with X-Correlation-ID)    |       (Internal Port: 3000)        |
|--------------------------------|                              |------------------------------------|
| 3. Middleware Validasi:        |                              | 7. Menerima request validasi.      |
|    - Memeriksa body request.   |                              | 8. Log dengan Correlation ID.      |
| 4. Handler:                    |                              | 9. Cek data di DB/Redis & balas.   |
|    - Menerima Correlation ID.  |                              |------------------------------------|
|    - Log dengan Correlation ID.|                              | 12. Menerima Event dari RabbitMQ.  |
| 5. Service:                    |                              | 13. Log dengan Correlation ID.     |
|    - Menyimpan order ke DB.    |                              | 14. Kurangi stok produk di DB.     |
|    - Menerbitkan Event ke      |                              +------------------------------------+
|      RabbitMQ dengan           |                                                ^
|      Correlation ID di         |                                                |
|      dalam payload.            |                                                |
+--------------------------------+                                                |
             |                                                                    | (Event dikonsumsi)
             | (Event diterbitkan)                                                |
             v                                                                    |
+--------------------------------+                                                |
|       RABBITMQ BROKER          |------------------------------------------------+
|      (orders_exchange)         |
+--------------------------------+

```

## 3. Prasyarat Instalasi

#### Pastikan telah menginstal **Git** dan **Docker Desktop** di sistem.

</br>

<details>
<summary><strong>Instruksi untuk Windows 10/11</strong></summary>

1.  **Instal Git**:

    - Unduh dan instal dari [git-scm.com](https://git-scm.com/downloads).

2.  **Instal Docker Desktop**:
    - Docker di Windows memerlukan WSL 2. Buka **PowerShell sebagai Administrator** dan jalankan:
      ```powershell
      wsl --install
      ```
    - **Restart komputer Anda** setelah instalasi WSL selesai.
    - Unduh dan instal [**Docker Desktop for Windows**](https://docs.docker.com/desktop/install/windows-install/).
3.  **Instal k6** (opsional): - Cara termudah adalah menggunakan Windows Package Manager (winget) yang sudah terpasang di Windows modern.
`powershell
        winget install k6.k6
        ` - Jika Anda tidak memiliki winget, Anda bisa mengunduh dan menjalankan installer `.msi` resmi dari [halaman rilis k6](https://github.com/grafana/k6/releases/latest) di GitHub.
</details>

</br>

<details>
<summary><strong>Instruksi untuk macOS</strong></summary>

1.  **Instal Git**:

    - Cara termudah adalah dengan menginstal Xcode Command Line Tools. Buka Terminal dan jalankan:
      ```bash
      xcode-select --install
      ```
    - Alternatif lain, gunakan [Homebrew](https://brew.sh/): `brew install git`.

2.  **Instal Docker Desktop**:

    - Unduh dan instal [**Docker Desktop for Mac**](https://docs.docker.com/desktop/install/mac-install/). Pastikan Anda memilih versi yang benar (Apple Silicon atau Intel Chip).

3.  **Instal k6** (opsional):
    ```bash
    brew install k6
    ```

</details>

</br>

<details>
<summary><strong>Instruksi untuk Linux (Ubuntu/Debian)</strong></summary>

1.  **Instal Git**:

    - Buka terminal dan jalankan:
      ```bash
      sudo apt update
      sudo apt install git
      ```

2.  **Instal Docker Engine & Compose Plugin**:

    - Prosesnya sedikit lebih manual. Ikuti panduan resmi untuk hasil terbaik:
      1.  [**Instal Docker Engine**](https://docs.docker.com/engine/install/ubuntu/) (ikuti metode "Install using the `apt` repository").
      2.  Setelah itu, instal Compose Plugin:
          ```bash
          sudo apt-get install docker-compose-plugin
          ```

3.  **instal k6** (opsional): - Anda bisa menginstalnya menggunakan `apt`.
`bash
       sudo gpg -k
       sudo gpg --no-default-keyring --keyring /usr/share/keyrings/k6-archive-keyring.gpg --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
       echo "deb [signed-by=/usr/share/keyrings/k6-archive-keyring.gpg] https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
       sudo apt-get update
       sudo apt-get install k6
       `
</details>
</br>
</br>

## 4. Menjalankan Aplikasi

Setelah semua prasyarat terpenuhi, langkah-langkahnya sama untuk semua sistem operasi.

1.  **Clone repositori ini:** [https://github.com/Angellieta/NestJS-Golang-Microservice-Architecture.git](https://github.com/Angellieta/NestJS-Golang-Microservice-Architecture.git)

    ```bash
    git clone https://github.com/Angellieta/NestJS-Golang-Microservice-Architecture.git
    cd fullstack-challenge
    ```

2.  **Jalankan dengan Docker Compose:**
    Buka terminal di direktori utama project dan jalankan perintah berikut untuk membangun dan memulai semua layanan.

    ```bash
    docker-compose up --build
    ```

3.  **Aplikasi Siap!**
    Semua permintaan ke aplikasi harus dilakukan melalui **BFF Service** - **BFF (Entrypoint)**: `http://localhost:5000` - **RabbitMQ Management UI**: `http://localhost:15672` (login: `guest` / `guest`)
    </br>
    </br>

## 5. Penggunaan API

Berikut adalah contoh **request** API, gunakan `localhost:5000` untuk semua _request_.

</br>


### Product Service

#### 1. Membuat Produk Baru

```powershell
Invoke-WebRequest -Uri http://localhost:5000/products -Method POST -ContentType "application/json" -Body '{"name": "Ice Cream", "price": 2000, "qty": 25}'
```

#### 2. Mengambil Produk Berdasarkan ID (Cached)

> Ganti <PRODUCT_ID> dengan ID valid yang muncul saat produk baru dibuat

```powershell
Invoke-WebRequest -Uri http://localhost:5000/products/<PRODUCT_ID>
```

</hr>

</br>

---

</br>


### Order Service

#### 3. Membuat Order Baru (dengan Validasi Produk)

> Ganti <PRODUCT_ID> dengan ID valid yang muncul saat produk baru dibuat

```powershell
Invoke-WebRequest -Uri http://localhost:5000/orders -Method POST -ContentType "application/json" -Body '{"productId": "<PRODUCT_ID>", "price": 2000, "qty": 3}'
```

#### 4. Mengambil Order Berdasarkan Product ID (Cached)

> Ganti <PRODUCT_ID> dengan ID valid yang muncul saat produk baru dibuat

```powershell
Invoke-WebRequest -Uri http://localhost:5000/orders/product/<PRODUCT_ID>
```

#### 5. Mengambil Ringkasan Produk dan Order

> Ganti <PRODUCT_ID> dengan ID valid yang muncul saat produk baru dibuat

```powershell
Invoke-WebRequest -Uri http://localhost:5000/order-summary/<PRODUCT_ID>
```

</br>

## 6. Pengujian

**Unit Tests**

- Product Service (NestJS):

```bash
cd product-service
npm test
```

- Order Service (Go):

```bash
cd order-service
go test ./...
```

**Performance Test**

- Pastikan semua layanan sedang berjalan.
- Pastikan sudah menginstal k6
- Dari direktori utama, jalankan:

```bash
k6 run performance-tests/create-order.js
```
