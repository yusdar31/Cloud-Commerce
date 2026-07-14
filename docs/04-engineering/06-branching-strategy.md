# Branching Strategy

**Project:** CloudCommerce

**Version:** 1.0.0

**Status:** Approved

**Owner:** Engineering Team

**Last Updated:** July 2026

---

# 1. Purpose

Dokumen ini mendefinisikan branching strategy dan Git workflow untuk CloudCommerce.

Tujuan:

- Commit history yang rapi dan profesional
- Memudahkan code review
- Mendukung CI/CD pipeline
- Memudahkan tracking fitur dan bug

---

# 2. Branch Model

CloudCommerce menggunakan **Trunk-based Development** yang disederhanakan.

```
main ────────●────────────●────────────●────────
                \        /            /
feature/xxx      ●──●──●             /
                                    /
fix/xxx          ●──●──●───────────●
```

## Branch Types

| Branch | Naming | Source | Merge To |
|--------|--------|--------|----------|
| Main | `main` | - | - |
| Feature | `feat/{nama-fitur}` | `main` | `main` |
| Fix | `fix/{nama-bug}` | `main` | `main` |
| Chore | `chore/{task}` | `main` | `main` |
| Docs | `docs/{topik}` | `main` | `main` |

## Branch Naming Convention

```
feat/identity-service
feat/catalog-crud
feat/storefront-landing

fix/cors-error
fix/stock-race-condition

chore/docker-compose
chore/helm-setup

docs/api-guidelines
docs/readme
```

---

# 3. Commit Convention (Conventional Commits)

```
<type>(<scope>): <description>

feat(catalog): add product CRUD endpoints
fix(order): fix stock reservation race condition
chore(docker): add Dockerfile for identity service
docs(readme): add architecture diagram
```

## Types

| Type | Usage |
|------|-------|
| `feat` | New feature |
| `fix` | Bug fix |
| `chore` | Tooling, config, dependencies |
| `docs` | Documentation only |
| `refactor` | Code change without feature/fix |
| `test` | Adding or fixing tests |
| `style` | Formatting, linting |
| `perf` | Performance improvement |

## Scope (optional)

| Scope | Area |
|-------|------|
| `(identity)` | Identity service |
| `(catalog)` | Catalog service |
| `(gateway)` | API Gateway |
| `(storefront)` | Buyer frontend |
| `(dashboard)` | Seller dashboard |
| `(docker)` | Docker config |
| `(k8s)` | Kubernetes manifest |
| `(ci)` | GitHub Actions |
| `(helm)` | Helm chart |

Examples:

```
feat(catalog): implement product search by name
fix(order): handle payment webhook timeout
chore(ci): add Go lint step to pipeline
docs(readme): add deployment prerequisites
refactor(identity): extract JWT validation to middleware
test(payment): add webhook idempotency test
```

---

# 4. Pull Request Process

## PR Title

Sama dengan commit convention:

```
feat(catalog): add product CRUD endpoints
```

## PR Template

```markdown
## Description

Brief description of changes.

## Type

- [ ] Feature
- [ ] Bug fix
- [ ] Chore
- [ ] Docs

## Related Issues

Closes #123

## Checklist

- [ ] Code compiles
- [ ] Tests pass
- [ ] Lint passes
- [ ] Self-review done

## Screenshots (if applicable)

## Notes for Reviewer

Any additional context?
```

## Rules

- PR minimal 1 reviewer approval
- PR harus squash merge ke `main`
- PR description harus jelas

---

# 5. Commit History

Setelah squash merge, `main` akan memiliki commit yang rapi:

```
feat(catalog): add product CRUD endpoints
feat(identity): implement JWT authentication
chore(docker): add Docker Compose for local dev
fix(order): resolve stock race condition
chore(ci): setup GitHub Actions pipeline
docs(readme): update architecture diagram
```

No merge commits (`Merge branch 'feature/xxx' into main`) on `main`.

---

# 6. Workflow Example

## Start New Feature

```bash
git checkout main
git pull
git checkout -b feat/catalog-search
```

## During Development

```bash
git add .
git commit -m "feat(catalog): add search query parameter"
git push origin feat/catalog-search
```

## Create PR

- Push branch
- Open PR on GitHub: `feat/catalog-search` → `main`
- Wait for CI (lint + test + build)
- Request review

## Merge

After approval, squash merge to `main`.

---

# 7. CI Integration

Setiap push akan trigger:

| Event | Action |
|-------|--------|
| Push to `feat/*` | Lint + Test + Build |
| Push to `fix/*` | Lint + Test + Build |
| Push to `main` | Lint + Test + Build + Docker |
| PR to `main` | Lint + Test + Build |

CI harus hijau sebelum PR bisa di-merge.

---

# 8. Release Process

Tidak ada release branch terpisah.

Cukup tag di `main`:

```bash
git checkout main
git tag v0.1.0
git push origin v0.1.0
```

Tag akan trigger deployment pipeline.

## Versioning (Semantic)

```
v0.1.0    # MVP release
v0.2.0    # New feature
v0.2.1    # Bug fix
v1.0.0    # Production ready
```

---

# 9. Repository Rules

- Jangan commit langsung ke `main`
- Jangan commit `.env`, secrets, atau credentials
- Jangan commit `node_modules/`, `vendor/`, binary files
- Gunakan `.gitignore` yang tepat

---

# 10. Related Documents

- Coding Standards
- CI/CD Strategy
- Release Strategy
- Monorepo Structure
