# 🔐 Sampora Auth & Access Control (SSO)
**Versi:** 2.2 (Standardized Handshake)

Dokumentasi ini mengatur bagaimana aplikasi **POS** terhubung dengan ekosistem identitas **Sampora System**.

---

### 1. Alur Autentikasi (Redirect Handshake)
Sistem menggunakan mekanisme **Dynamic Return URL**. Seluruh kredensial dikelola terpusat oleh Auth Center.

#### Diagram Alur:
```mermaid
sequenceDiagram
    participant User
    participant POS as POS System (Client)
    participant Auth as Auth Center (Provider)

    User->>POS: Akses POS Dashboard
    POS->>POS: Cek Sesi Lokal (Kosong)
    POS->>User: Redirect ke auth.sampora.my.id?from=pos.sampora.my.id
    Auth->>User: Form Login (Input Email/Pass)
    User->>Auth: Submit Credentials
    Auth->>User: Sukses! Redirect ke pos.sampora.my.id?token=[JWT]
    User->>POS: Kembali ke POS dengan Token
    POS->>POS: Simpan Token & Bersihkan URL
    POS->>User: Dashboard POS Terbuka

### 2. Integrasi Backend POS
Tim backend cukup memvalidasi Token JWT yang dikirimkan oleh Frontend melalui Header.
Header: Authorization: Bearer <JWT>
Payload JWT: Berisi id (User ID), email, dan app_role (admin/cashier/supervisor).
Kebijakan: Gunakan pb_user_id sebagai Foreign Key pada tabel transaksi/stok di database POS.

### 3. Kebijakan Global Logout
Untuk keluar dari sesi secara bersih di seluruh subdomain:
Jalankan pb.authStore.clear().
Redirect user ke: https://auth.sampora.my.id?logout=true
