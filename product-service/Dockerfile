# product-service/Dockerfile

# 1. Base Image: Menggunakan Node.js versi 18
FROM node:18-alpine

# 2. Set working directory di dalam container
WORKDIR /usr/src/app

# 3. Copy package.json dan package-lock.json
COPY package*.json ./

# 4. Install dependencies
RUN npm install

# 5. Copy seluruh kode aplikasi
COPY . .

# 6. Build aplikasi NestJS
RUN npm run build

# 7. Expose port yang digunakan aplikasi
EXPOSE 3000

# 8. Perintah untuk menjalankan aplikasi saat container dimulai
CMD ["node", "dist/main"]