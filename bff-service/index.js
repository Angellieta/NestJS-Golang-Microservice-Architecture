// bff-service/index.js
import express from "express";
import axios from "axios";
import { v4 as uuidv4 } from "uuid";

const app = express();
app.use(express.json());

const PORT = 5000;
const PRODUCT_SERVICE_URL = process.env.PRODUCT_SERVICE_URL;
const ORDER_SERVICE_URL = process.env.ORDER_SERVICE_URL;
const CORRELATION_ID_HEADER = "x-correlation-id";

// Middleware untuk menangani Correlation ID
app.use((req, res, next) => {
  let correlationId = req.headers[CORRELATION_ID_HEADER];
  if (!correlationId) {
    correlationId = uuidv4();
    req.headers[CORRELATION_ID_HEADER] = correlationId;
  }
  // Atur ID di header respons juga, agar klien bisa melihatnya
  res.set(CORRELATION_ID_HEADER, correlationId);
  next();
});

// Helper untuk meneruskan header
const getForwardingHeaders = (req) => {
  const headers = {};
  const correlationId = req.headers[CORRELATION_ID_HEADER];
  if (correlationId) {
    headers[CORRELATION_ID_HEADER] = correlationId;
  }
  return headers;
};

// Endpoint proxy sederhana ke product-service
app.get("/products/:id", async (req, res) => {
  try {
    const { id } = req.params;
    const response = await axios.get(`${PRODUCT_SERVICE_URL}/products/${id}`, {
      headers: getForwardingHeaders(req),
    });
    res.status(response.status).json(response.data);
  } catch (error) {
    res
      .status(error.response?.status || 500)
      .json(error.response?.data || { message: "An error occurred" });
  }
});

app.post("/products", async (req, res) => {
  try {
    const response = await axios.post(
      `${PRODUCT_SERVICE_URL}/products`,
      req.body,
      {
        headers: getForwardingHeaders(req),
      }
    );
    res.status(response.status).json(response.data);
  } catch (error) {
    res
      .status(error.response?.status || 500)
      .json(error.response?.data || { message: "An error occurred" });
  }
});

// Endpoint proxy sederhana ke order-service
app.post("/orders", async (req, res) => {
  try {
    const response = await axios.post(`${ORDER_SERVICE_URL}/orders`, req.body, {
      headers: getForwardingHeaders(req),
    });
    res.status(response.status).json(response.data);
  } catch (error) {
    res
      .status(error.response?.status || 500)
      .json(error.response?.data || { message: "An error occurred" });
  }
});

// Endpoint gabungan untuk mendapatkan detail produk beserta daftar ordernya
app.get("/order-summary/:productId", async (req, res) => {
  try {
    const { productId } = req.params;
    const headers = getForwardingHeaders(req);

    const productResponse = await axios.get(
      `${PRODUCT_SERVICE_URL}/products/${productId}`,
      { headers }
    );
    const orderResponse = await axios.get(
      `${ORDER_SERVICE_URL}/orders/product/${productId}`,
      { headers }
    );

    const productData = productResponse.data;
    const ordersData = orderResponse.data;

    // Menggabungkan hasilnya menjadi satu objek respons
    const summary = {
      product: {
        id: productData.id,
        name: productData.name,
        price: productData.price,
        currentStock: productData.qty,
      },
      orders: ordersData,
      totalOrders: ordersData.length,
    };

    res.status(200).json(summary);
  } catch (error) {
    // Penanganan error jika salah satu service gagal merespons
    res.status(error.response?.status || 500).json(
      error.response?.data || {
        message: "An error occurred fetching order summary",
      }
    );
  }
});

app.listen(PORT, () => {
  console.log(`BFF Service is running on port ${PORT}`);
});
