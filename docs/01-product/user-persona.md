User Personas: CloudCommerce MVP

Dokumen ini menjabarkan profil pengguna utama platform CloudCommerce. Memahami persona ini sangat penting karena setiap profil memiliki ekspektasi yang berbeda terhadap performa, keamanan, dan ketersediaan sistem, yang secara langsung memengaruhi rancangan infrastruktur dan Non-Functional Requirements (NFRs).

1. The Platform Operator (Internal)

Nama: Yusdar (22 Tahun)

Peran: DevOps Engineer / Cloud Engineer

Tingkat Keahlian Teknis:  Tinggi (SME di Kubernetes, CI/CD, Cloud)

Profil Singkat

Yusdar adalah pengelola sistem CloudCommerce. Tugas utamanya adalah memastikan platform berjalan mulus, aman, dan efisien secara biaya. Ia mengelola ratusan tenant (toko) yang berjalan di atas satu infrastruktur cluster yang sama.

Goals & Motivations

Mengotomatisasi seluruh proses provisioning infrastruktur dan deployment aplikasi.

Mempertahankan SLA (Service Level Agreement) 99.9% uptime.

Menjaga biaya infrastruktur cloud serendah mungkin (FinOps) tanpa mengorbankan stabilitas sistem.

Memiliki visibilitas penuh terhadap kesehatan sistem (metrics, logs, traces).

Pain Points (Masalah yang Dihadapi)

Toil (pekerjaan manual yang berulang) untuk men-setup tenant baru.

Kurangnya peringatan dini (alerting) jika ada layanan yang bermasalah, sehingga baru tahu saat pelanggan komplain.

Noisy Neighbor problem (satu merchant yang sedang viral menghabiskan seluruh resource CPU/RAM, membuat toko lain down).

DevOps & Infrastructure Implications

GitOps & IaC: Membutuhkan Terraform untuk provisioning GCP dan ArgoCD untuk deployment continuous.

Observability Stack: Wajib mengimplementasikan Prometheus, Grafana, dan Alertmanager untuk monitoring proaktif.

Resource Quotas & Rate Limiting: Konfigurasi Traefik API Gateway untuk membatasi request per tenant, serta requests/limits di level Pod K8s.

Cost Optimization: Menjalankan cluster K3s di atas Preemptible/Spot VMs dengan kapabilitas auto-healing.

2. The Tenant (Merchant / Seller)

Nama: Budi (34 Tahun)

Peran: Pemilik UMKM / Kreator Digital

Tingkat Keahlian Teknis: Rendah - Menengah

Profil Singkat

Budi memiliki brand pakaian lokal. Ia lelah berjualan di marketplace karena perang harga, biaya admin/komisi yang terus naik, dan tidak bisa mendapatkan data analitik pelanggannya secara penuh. Ia ingin toko online mandiri yang terlihat profesional namun mudah dikelola.

Goals & Motivations

Memiliki toko online dengan link sendiri dalam hitungan menit tanpa perlu mengerti coding atau hosting.

Proses manajemen katalog produk dan inventaris yang simpel.

Bisa langsung menerima pembayaran dari berbagai metode populer di Indonesia (QRIS, GoPay, Virtual Account).

Data pelanggan (nama, email, riwayat beli) tersimpan eksklusif untuk tokonya sendiri (tidak dibagi dengan toko lain).

Pain Points (Masalah yang Dihadapi)

Biaya berlangganan platform e-commerce global (seperti Shopify) terlalu mahal karena menggunakan USD.

Kesulitan jika harus mengintegrasikan payment gateway lokal secara manual.

Khawatir datanya bocor atau tercampur dengan data kompetitor di platform yang sama.

DevOps & Infrastructure Implications

Frictionless Onboarding: Proses pembuatan schema database atau alokasi tenant_id harus otomatis (event-driven) saat registrasi berhasil.

Strict Data Isolation: Implementasi Row-Level Security (RLS) di PostgreSQL atau pemisahan schema logikal untuk menjamin Merchant A tidak bisa melihat data Merchant B.

Zero-Downtime Deployment: Pembaruan aplikasi oleh tim Platform (Alex) tidak boleh mengganggu operasional toko Budi.

3. The End-Customer (Buyer)

Nama: Citra (24 Tahun)

Peran: Pembeli Online Aktif

Tingkat Keahlian Teknis: Menengah (Terbiasa menggunakan aplikasi smartphone dan transaksi digital)

Profil Singkat

Citra adalah pengikut setia brand Budi di Instagram. Ia mengklik link toko Budi di bio Instagram untuk membeli produk rilisan terbaru yang sedang diskon atau flash sale.

Goals & Motivations

Pengalaman berbelanja (browsing katalog, memasukkan ke keranjang) yang mulus dan instan tanpa jeda muat yang lama.

Proses checkout dan pembayaran yang jelas, aman, dan memvalidasi stok secara real-time.

Mendapatkan notifikasi instan via email atau WhatsApp ketika pembayaran berhasil dan barang dikirim.

Pain Points (Masalah yang Dihadapi)

Halaman web lambat dimuat, terutama saat koneksi internet kurang stabil.

Produk yang sudah di-checkout ternyata stoknya kosong (karena race condition dengan pembeli lain).

Pesan error yang tidak jelas saat proses pembayaran gagal.

DevOps & Infrastructure Implications

Low Latency Catalog: API pembacaan katalog produk harus sangat cepat (P95 < 100ms), membutuhkan lapisan caching Redis yang efisien.

High Availability & Event-Driven: Proses checkout tidak boleh gagal meskipun layanan notifikasi sedang mati. Hal ini dicapai menggunakan NATS message broker untuk memisahkan proses (decoupling).

Concurrency Control: Desain sistem pada Order Service harus mampu menangani ribuan request checkout bersamaan (flash sale) tanpa membuat database PostgreSQL terkunci (deadlock).