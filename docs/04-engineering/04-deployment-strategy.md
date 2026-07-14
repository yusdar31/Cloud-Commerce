# Deployment Strategy

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Draft

**Owner:** DevOps Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan strategi deployment untuk tiga environment:

- **Local Development** — Kind/Minikube untuk development sehari-hari
- **Cloud** — GCP (target) untuk production-like environment
- **On-Premise** — k3s di Proxmox/homelab untuk hybrid/portability demo

---

# 2. Architecture Principle

Satu source code dan Docker image yang sama untuk semua environment.

Perbedaan hanya pada:

- Infrastructure provisioning (Terraform)
- Kubernetes manifest values (Helm)
- Resource allocation
- External service endpoints

---

# 3. Environment Comparison

| Aspek | Local (Kind) | Cloud (GCP) | On-Prem (k3s) |
|-------|-------------|-------------|----------------|
| Cluster | Kind (1 node) | GKE (3 node) | k3s (3 node) |
| Ingress | MetalLB + NGINX | Cloud Load Balancer + NGINX | MetalLB + NGINX |
| TLS | mkcert (self-signed) | cert-manager + Let's Encrypt | cert-manager + Let's Encrypt |
| PostgreSQL | StatefulSet (lmits) | Cloud SQL | StatefulSet + PVC |
| Redis | Deployment (1 replica) | Memorystore | Deployment (1 replica) |
| Object Storage | MinIO (Deployment) | GCS | MinIO (Deployment) |
| NATS | Deployment (1 replica) | NATS JetStream (3 replicas) | Deployment (1 replica) |
| Monitoring | kube-prometheus-stack | Managed Prometheus + Grafana | kube-prometheus-stack |
| Logging | Loki (single binary) | Grafana Loki (managed) | Loki (single binary) |
| CI/CD | GitHub Actions → Kind | GitHub Actions → GKE | GitHub Actions → k3s |
| Secrets | Sealed Secrets | External Secrets + Secret Manager | Sealed Secrets |

---

# 4. Local Development (Kind)

## Target

Developer dapat menjalankan seluruh sistem di laptop tanpa cloud.

## Resource Requirement

| Resource | Minimal |
|----------|---------|
| CPU | 4 cores |
| RAM | 8 GB |
| Storage | 20 GB |
| Tools | Docker, Kind, kubectl, Helm |

## Architecture

```
Laptop
    │
Kind Cluster (1 node)
    │
├── NGINX Ingress
├── API Gateway
├── 9 Microservices
├── PostgreSQL (StatefulSet)
├── Redis (Deployment)
├── NATS (Deployment)
├── MinIO (Deployment)
└── kube-prometheus-stack
```

## Networking

```
Kind
    │
MetalLB (172.18.0.x)
    │
NGINX Ingress
    │
curl / browser → http://localhost:80
```

## Services

External akses via port-forward:

```
kubectl port-forward svc/api-gateway 8080:80
kubectl port-forward svc/grafana 3000:3000
kubectl port-forward svc/minio 9001:9001
```

## Provisioning

```
git clone cloudcommerce

make kind-create

make deploy-all

make kind-destroy
```

---

# 5. Cloud Deployment (GCP)

## Target

Environment production-like untuk demo portfolio.

## Resource Requirement

| Resource | Estimate |
|----------|----------|
| GKE Cluster | 3 x e2-medium (2 vCPU, 4 GB) |
| Cloud SQL | db-f1-micro (1 vCPU, 0.6 GB) |
| Total cost | ~$80-120/month (gunakan $300 credit) |

## Architecture

```
Internet
    │
Cloud Load Balancer
    │
GKE Cluster (3 node)
    │
├── cert-manager (Let's Encrypt)
├── NGINX Ingress Controller
├── API Gateway
├── 9 Microservices
├── Cloud SQL (PostgreSQL)
├── Memorystore (Redis)
├── GCS (Object Storage)
├── NATS JetStream (3 replicas)
└── Managed Grafana + Prometheus
```

## Networking

```
Internet
    │
Cloud DNS (cloudcommerce.my.id)
    │
Cloud Load Balancer (HTTPS)
    │
GKE Ingress
    │
NGINX Ingress Controller
    │
API Gateway (ClusterIP)
    │
Services
```

## Terraform Structure

```
infra/terraform/
│
├── cloud/
│   ├── main.tf           # Provider, backend
│   ├── network.tf        # VPC, subnet, firewall
│   ├── gke.tf            # GKE cluster
│   ├── databases.tf      # Cloud SQL, Memorystore
│   ├── storage.tf        # GCS bucket
│   ├── dns.tf            # Cloud DNS
│   └── outputs.tf
│
└── modules/
    ├── gke/
    ├── cloud-sql/
    └── memorystore/
```

## CI/CD Pipeline

```
Git Push
    │
GitHub Actions
    │
├── Terraform Plan
├── Terraform Apply
├── Docker Build & Push
└── Deploy Helm via ArgoCD
```

## Limitation

GCP free-tier credit $300 habis dalam ~3 bulan.
Disarankan untuk:

- Demo hanya saat diperlukan
- Gunakan shutdown scheduling (non-24/7)
- Dokumentasikan dengan screenshot/video

---

# 6. On-Premise Deployment (k3s)

## Target

Menunjukkan portability — sistem yang sama berjalan di infrastruktur sendiri.

## Resource Requirement

| Node | CPU | RAM | Storage | OS |
|------|-----|-----|---------|----|
| Master | 2 vCPU | 4 GB | 40 GB | Ubuntu 24.04 |
| Worker 1 | 2 vCPU | 4 GB | 40 GB | Ubuntu 24.04 |
| Worker 2 | 2 vCPU | 4 GB | 60 GB (data) | Ubuntu 24.04 |
| **Total** | **6 vCPU** | **12 GB** | **140 GB** | |

## Hardware Recommendation

| Setup | Spesifikasi | Notes |
|-------|-------------|-------|
| **Minimal** | Ryzen 5 / Core i5, 16 GB RAM, 256 GB NVMe | 1 physical machine → 3 VM via Proxmox |
| **Recommended** | Ryzen 7 / Core i7, 32 GB RAM, 512 GB NVMe | Lebih lega untuk monitoring + cadangan |
| **VM Storage** | 100-150 GB SSD | Persistent volumes untuk database |

## Architecture

```
Proxmox Host
    │
├── Master VM (2C/4G)
├── Worker-1 VM (2C/4G)
└── Worker-2 VM (2C/4G)
    │
k3s Cluster
    │
├── MetalLB (load balancer IP)
├── NGINX Ingress Controller
├── cert-manager (Let's Encrypt)
├── API Gateway
├── 9 Microservices
├── PostgreSQL (StatefulSet + PVC)
├── Redis (Deployment)
├── NATS (Deployment)
├── MinIO (Deployment + PVC)
├── kube-prometheus-stack
└── Loki (StatefulSet + PVC)
```

## Networking

```
Internet (port forwarding router)
    │
Proxmox Host (public IP)
    │
Port 443 → Master Node
    │
MetalLB (192.168.1.200-192.168.1.250)
    │
NGINX Ingress Controller
    │
API Gateway
    │
Services
```

## Terraform Structure

```
infra/terraform/
│
├── on-prem/
│   ├── main.tf           # Provider (Proxmox)
│   ├── vms.tf            # VM provisioning
│   ├── network.tf        # Bridge, firewall
│   ├── storage.tf        # Backup config
│   ├── k3s.tf            # k3s install (via ssh)
│   └── outputs.tf
│
└── modules/
    ├── proxmox-vm/
    └── k3s-install/
```

## Provisioning Flow

```
1. Provision VM via Terraform (Proxmox)
         │
2. Install k3s via Ansible script
         │
3. Deploy Helm charts (sama dengan cloud)
         │
4. Configure MetalLB IP range
         │
5. Install cert-manager + Ingress
```

## Ansible / Script

```
scripts/bootstrap-k3s.sh
scripts/join-worker.sh
scripts/deploy-all.sh
```

---

# 7. Helm Strategy

Satu Helm chart untuk semua environment, beda values file.

```
infra/helm/
│
├── charts/
│   ├── api-gateway/
│   ├── identity-service/
│   ├── catalog-service/
│   ├── order-service/
│   └── ... (per service)
│
├── environments/
│   ├── values-local.yaml
│   ├── values-cloud.yaml
│   └── values-onprem.yaml
│
└── global.yaml
```

Contoh perbedaan values:

```yaml
# values-local.yaml
postgresql:
  storageClass: standard

ingress:
  tls: false

replicas: 1
```

```yaml
# values-cloud.yaml
postgresql:
  cloudSQL: true
  connectionName: project:region:instance

ingress:
  tls: true
  certManager: true

replicas: 2
hpa: true
```

```yaml
# values-onprem.yaml
postgresql:
  storageClass: local-path
  size: 10Gi

ingress:
  tls: true
  certManager: true

replicas: 2
```

---

# 8. Makefile Targets

```makefile
# Local
kind-create       # Buat Kind cluster
kind-deploy       # Deploy semua service ke Kind
kind-destroy      # Hapus Kind cluster

# Cloud
gcp-init          # Terraform init untuk GCP
gcp-apply         # Provision GCP infrastructure
gcp-deploy        # Deploy Helm ke GKE
gcp-destroy       # Hancurkan semua resource GCP

# On-Prem
onprem-init       # Terraform init untuk Proxmox
onprem-apply      # Provision VM via Terraform
onprem-k3s        # Install k3s via Ansible
onprem-deploy     # Deploy Helm ke k3s
onprem-destroy    # Hapus semua VM
```

---

# 9. CI/CD Matrix

| Environment | Trigger | Action |
|-------------|---------|--------|
| Local | Manual (developer) | `make kind-deploy` |
| Dev/Staging | Push ke branch `develop` | Deploy ke Kind (CI) |
| Production (Cloud) | Merge ke `main` + manual approval | Deploy ke GKE via ArgoCD |
| Production (On-Prem) | Merge ke `main` + manual approval | Deploy ke k3s via ArgoCD |

---

# 10. Related Documents

- Technology Stack
- Monorepo Structure
- CI/CD Strategy
- Container Diagram
- Service Boundaries