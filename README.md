# 🚀 Web POS Backend with Online Ordering & Midtrans Integration

Backend untuk **Point of Sales (POS) berbasis Web** yang sudah mendukung **Online Ordering** dan integrasi dengan **Midtrans Payment Gateway**.  
Dibangun dengan arsitektur terpisah antara frontend & backend agar lebih fleksibel dan scalable. 💡

---

## ✨ Fitur Utama
- 🔐 **Authentication & Authorization** (JWT, Role-based Access)
- 📦 **Manajemen Produk & Inventori** (CRUD Produk, Kategori, Stok)
- 🛒 **Penjualan (POS Screen)** → Cart, Checkout, Struk
- 👥 **Manajemen Pelanggan** (opsional, histori pembelian)
- 📊 **Laporan & Dashboard** (penjualan, produk terlaris, stok menipis)
- 👤 **User Management** (Admin, Kasir, Supervisor)
- 🌐 **Online Ordering** → customer bisa order via web
- 💳 **Integrasi Midtrans (Snap API + Webhook)**

---

## 🗂️ Desain Database (Ringkasan Tabel)
- **users** → data akun user
- **products, categories** → produk & kategori
- **transactions, transaction_items** → penjualan offline
- **customers** → data pelanggan
- **stock_movements** → histori pergerakan stok
- **orders, order_items** → order online
- **payments** → integrasi Midtrans
- **webhook_logs** → log notifikasi Midtrans

---

## 🏗️ Arsitektur Backend
- 🔑 **Auth Service** → login, JWT, role-based middleware  
- 📦 **Product Service** → CRUD produk, stok  
- 💰 **Transaction Service** → penjualan offline/online, pembayaran  
- 📊 **Report Service** → laporan penjualan  
- 👥 **User Service** → manajemen akun  
- 🌐 **Order Service** → online ordering  
- 💳 **Payment Gateway Integration** → Midtrans Snap API + Webhook  

---

## 🔌 API Endpoint (REST)
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/api/auth/login` | Login user |
| `POST` | `/api/orders` | Buat order baru |
| `POST` | `/api/payments/:order_id/create` | Buat pembayaran via Midtrans |
| `POST` | `/api/webhooks/midtrans` | Endpoint notifikasi dari Midtrans |
| `GET`  | `/api/orders/:id/status` | Cek status order |

---

## 💳 Alur Integrasi Midtrans
1. 🛒 Customer buat order → disimpan ke tabel `orders` (status: `pending`)  
2. ⚡ Backend memanggil **Midtrans Snap API** → generate Snap Token  
3. 💾 Simpan record ke tabel `payments` (status: `pending`)  
4. 🎨 Frontend tampilkan **Snap Redirect/Popup** ke customer  
5. 💳 Customer menyelesaikan pembayaran  
6. 📩 Midtrans kirim **Webhook Notification** → backend menerima  
7. 🔐 Backend verifikasi `signature_key` dari Midtrans  
8. ✅ Update tabel `payments` & `orders` → status `paid` / `failed`  

---

## 🔒 Best Practices & Keamanan
- 🧪 Gunakan **sandbox/test key** saat development  
- 🔐 Selalu verifikasi `signature_key` dari Midtrans pada webhook  
- 🌐 Gunakan **HTTPS** untuk endpoint webhook  
- 📝 Simpan **raw webhook** di tabel `webhook_logs` untuk audit  
- ♻️ Pastikan update status pembayaran **idempotent** (hindari double update)  
- 🔑 Simpan **API keys di environment variables**, jangan hardcode  

---

## 📌 Tech Stack

- 🖥️ **Backend**: Golang (Gin)
- 🗄️ **Database**: PostgreSQL / MySQL
- 💳 **Payment**: Midtrans (Snap API + Webhook)
- 🔑 **Auth**: JWT, Role-based Access

---

## ⚡ Cara Menjalankan
```bash
# 1. Clone repository
git clone https://github.com/username/pos-backend.git
cd pos-backend

# 2. Setup environment
cp .env.example .env
# lalu isi dengan konfigurasi database & Midtrans API keys

# 3. Jalankan server
go run main.go
```

---

## 📜 Lisensi

- MIT License © 2025 [Naza Zulfiqi](https://www.nazazulfiqi.me/)