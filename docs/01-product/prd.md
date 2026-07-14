Product Requirements Document (PRD): CloudCommerce MVP

Document Version: 1.0
Status: Approved for MVP Development
Product: CloudCommerce (Multi-Tenant SaaS E-commerce)

1. Introduction & Purpose

Dokumen ini menjelaskan spesifikasi fungsional dan non-fungsional untuk fase Minimum Viable Product (MVP) dari CloudCommerce. Tujuan MVP ini adalah memvalidasi arsitektur dasar multi-tenant, kelancaran pipeline CI/CD infrastruktur, dan keberhasilan pemrosesan transaksi e-commerce end-to-end dengan logical data isolation yang ketat antar merchant.

2. User Roles & Personas

Sistem harus mendukung tiga jenis entitas pengguna dengan hak akses yang terpisah:

Platform Admin (Internal): Operator sistem/DevOps yang mengelola infrastruktur K8s, memantau resource usage, dan memiliki akses ke level database/log (melalui tools observability).

Merchant / Seller (Tenant): Pemilik toko. Mereka mendaftar ke platform, mendapatkan Workspace/Tenant ID, dan memiliki akses ke Seller Dashboard untuk mengelola produk dan pesanan milik mereka sendiri.

Buyer (End-Customer): Pembeli anonim atau terdaftar yang mengunjungi Storefront spesifik milik Merchant, menambahkan produk ke keranjang, dan melakukan checkout.

3. System Components Scope

MVP akan diimplementasikan melalui 9 microservices dan 1 API Gateway:

NGINX Ingress Controller: TLS termination dan external traffic routing ke API Gateway.

API Gateway (Go + Gin): Routing internal, autentikasi JWT, rate limiting per tenant, request logging, dan API versioning. Gateway tidak memiliki business logic.

Identity Service: Autentikasi, otorisasi, dan manajemen JWT.

Tenant Service: Manajemen merchant, konfigurasi toko, dan branding.

Catalog Service: Manajemen produk dan kategori dengan caching Redis.

Inventory Service: Manajemen stok dan reservasi produk.

Order Service: Orkestrasi checkout dan manajemen state pesanan.

Payment Service: Integrasi Payment Gateway (Midtrans/Xendit) dan penanganan webhook.

Notification Service: Event-driven notification dispatcher (email via Mailhog).

Analytics Service: Dashboard metrics dan laporan penjualan (event-driven, data dari NATS).

Subscription Service: Manajemen paket langganan dan siklus trial tenant.

Frontend (Next.js): Menyajikan UI untuk Seller Dashboard dan Buyer Storefront.

4. Functional Requirements (Epics & User Stories)

Epic 1: Tenant Onboarding & Identity

FR 1.1 - Merchant Registration: Pengguna dapat mendaftar dengan email dan password. Sistem otomatis membuatkan tenant_id unik dan status toko "Active".

FR 1.2 - JWT Generation: Saat login berhasil, sistem mengembalikan JWT yang memuat informasi user_id, role, dan tenant_id.

FR 1.3 - Context Enforcement: API Gateway dan backend services WAJIB memvalidasi tenant_id dari JWT untuk setiap request yang masuk agar Merchant A tidak bisa mengakses data Merchant B.

Epic 2: Catalog & Inventory Management

FR 2.1 - Product CRUD: Merchant dapat membuat, membaca, mengubah, dan menghapus produk (Nama, Deskripsi, Harga, Stok).

FR 2.2 - Tenant Isolation: Query produk dari Storefront harus difilter secara ketat berdasarkan URL / identitas toko pembeli (misal: mengambil katalog spesifik untuk tenant_id = X).

FR 2.3 - Redis Caching: Daftar produk per tenant harus di-cache di Redis untuk mempercepat respons pembacaan dari Storefront. Cache invalidation terjadi setiap kali ada perubahan data produk.

Epic 3: Order & Checkout Orchestration

FR 3.1 - Shopping Cart: Buyer dapat menambahkan produk ke keranjang belanja (disimpan sementara via Redis atau Local Storage frontend).

FR 3.2 - Checkout Process: Saat checkout, Order Service harus membuat record pesanan dengan status PENDING, mengkalkulasi total harga, dan mengamankan stok produk sementara (soft-reserve).

FR 3.3 - Event Publishing: Setelah pesanan dibuat, Order Service mem-publikasikan event OrderCreated ke Message Broker (NATS).

Epic 4: Payment Processing & Webhooks

FR 4.1 - Payment PG Integration: Payment Service yang mendengarkan event OrderCreated akan memanggil API Midtrans/Xendit (Mode Sandbox) untuk mendapatkan Payment URL/Token, dan menyimpannya ke database.

FR 4.2 - Webhook Listener: Sistem menyediakan endpoint public yang aman (memvalidasi signature key) untuk menerima notifikasi pembayaran sukses dari Midtrans/Xendit.

FR 4.3 - Order State Update: Setelah webhook valid diterima, Payment Service mem-publikasikan event PaymentSuccess. Order Service mengonsumsi event ini dan mengubah status pesanan dari PENDING menjadi PAID.

Epic 5: Asynchronous Notifications

FR 5.1 - Email Dispatch: Notification Service mengonsumsi event PaymentSuccess dan mengirimkan email konfirmasi ke Buyer (menggunakan SMTP mock seperti Mailhog untuk MVP).

5. Non-Functional Requirements (NFRs)

Persyaratan non-fungsional ini adalah fondasi bagi desain infrastruktur K3s:

5.1. Performance & Scalability

Rate Limiting: API Gateway harus membatasi request per tenant_id (misal: max 100 req/menit) untuk mencegah noisy neighbor attack (satu toko membebani seluruh server).

Latency: Respons API untuk pembacaan katalog produk (cache hit) harus < 100ms. Pembuatan order harus < 300ms.

Statelessness: Semua microservices tidak boleh menyimpan state lokal (sesi/file) agar K3s dapat mematikan/menghidupkan pod (auto-scaling atau spot VM termination) kapan saja tanpa korupsi data.

5.2. Reliability & Resilience

Idempotency: Payment Webhook dan Order Service harus idempotent. Jika Midtrans mengirimkan webhook yang sama dua kali, sistem tidak boleh memproses pembayaran ganda.

Graceful Shutdown: Aplikasi Node.js harus memproses SIGTERM dari Kubernetes secara graceful (menyelesaikan koneksi database dan request HTTP yang berjalan sebelum mati).

5.3. Observability

Structured Logging: Semua service memancarkan log dalam format JSON.

Correlation ID: Setiap request yang masuk dari API Gateway harus disematkan X-Request-ID unik yang diteruskan ke seluruh komunikasi internal (NATS & HTTP) untuk distributed tracing.

Metrics Endpoint: Semua layanan harus mengekspos endpoint /metrics untuk di-scrape oleh Prometheus (CPU, RAM, API response time, active connections).

5.4. Security

Secret Management: Kredensial database, API Key Midtrans, dan JWT Secret tidak boleh ada dalam repositori kode. Harus di-inject menggunakan Kubernetes Secrets (atau Sealed Secrets).

Network Policies: Layanan internal (seperti Database dan NATS) tidak boleh terekspos ke internet publik. Hanya API Gateway yang memiliki External IP.

6. Integrations & Dependencies

PostgreSQL: Penyimpanan data persisten utama.

Redis: Caching dan manajemen state sementra.

NATS: Message broker ringan yang berkinerja tinggi untuk memfasilitasi komunikasi asinkron (Pub/Sub) antar microservices. NATS dipilih secara strategis dibandingkan solusi yang lebih berat (seperti Kafka) untuk menjaga efisiensi resource (memory footprint rendah) di dalam arsitektur K3s yang cost-optimized. Komponen ini sangat vital untuk mendekopel layanan; misalnya, ketika Order Service mempublikasikan event OrderCreated, layanan lain seperti Payment dan Notification dapat mengonsumsi event tersebut secara mandiri tanpa memblokir alur transaksi utama atau menciptakan single point of failure.

Midtrans / Xendit API: Eksternal API untuk pemrosesan pembayaran (Sandbox).

Mailhog: Sistem SMTP lokal/virtual untuk testing notifikasi email.

7. Out of Scope (For Future Phases)

Subscription billing (pemrosesan pembayaran langganan dari merchant). Subscription Service untuk manajemen paket/trial tetap ada di arsitektur, tetapi integrasi payment gateway untuk subscription fee tidak termasuk MVP.

Custom Domain Mapping (merchant membawa domain .com mereka sendiri).

Pengiriman / Integrasi logistik pihak ketiga (JNE, SiCepat).