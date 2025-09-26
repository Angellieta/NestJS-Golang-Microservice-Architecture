// performance-tests/create-order.js

import http from "k6/http";
import { check, sleep } from "k6";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

// Konfigurasi Tes
export const options = {
  stages: [
    { duration: "30s", target: 1000 }, // Naikkan beban ke 1000 "user" selama 30 detik
    { duration: "30s", target: 1000 }, // Tahan beban di 1000 "user" selama 30 detik
    { duration: "10s", target: 0 }, // Turunkan beban ke 0
  ],
  thresholds: {
    // Syarat kelulusan tes:
    // 1. Error rate harus di bawah 1%
    http_req_failed: ["rate<0.01"],
    // 2. 95% dari request harus selesai di bawah 2000ms (2 detik)
    http_req_duration: ["p(95)<2000"],
  },
};

export default function () {
  const url = "http://localhost:8080/orders";

  // Membuat payload dengan productId yang unik untuk setiap request
  // untuk mensimulasikan order produk yang berbeda-beda.
  const payload = JSON.stringify({
    productId: uuidv4(),
    price: 1500,
    qty: 1,
  });

  const params = {
    headers: {
      "Content-Type": "application/json",
    },
  };

  // Mengirim request POST
  const res = http.post(url, payload, params);

  // Mengecek apakah responsnya adalah 201 Created
  check(res, {
    "is status 201": (r) => r.status === 201,
  });

  sleep(1); // Menunggu 1 detik sebelum iterasi berikutnya
}
