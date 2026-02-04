# Scope Teknis & Task List (Portfolio Strategis)

Dokumen ini berisi scope teknis dan task list untuk fitur portofolio strategis yang direkomendasikan.

## Fitur 1: RBAC + Audit Log (Prioritas 1)

### Scope Teknis
- **RBAC** berbasis role: `supervisor`, `staff`.
- **Policy** akses untuk endpoint employee (create/update/delete) + membaca detail karyawan.
- **Audit log** untuk aksi kritis: create/update/delete employee + login/logout.
- **Konteks audit** mencakup `actor_id`, `actor_role`, `action`, `resource`, `resource_id`, `ip`, `user_agent`, `timestamp`.
- **Standard error response** saat akses ditolak (403) atau token invalid (401).

### Definisi Role & Policy (contoh)
- `supervisor`: full access (create/update/delete/list/detail).
- `staff`: read-only (list/detail).

### Skema Audit Log (contoh)
Kolom minimum yang direkomendasikan:
- `id` (uuid)
- `actor_id` (uuid)
- `actor_role` (text)
- `action` (text) — `employee.create`, `employee.update`, `employee.delete`, `auth.login`, `auth.logout`
- `resource` (text) — `employee`, `auth`
- `resource_id` (uuid, nullable)
- `ip` (text, nullable)
- `user_agent` (text, nullable)
- `created_at` (timestamp)

### Alur Request (ringkas)
1. Middleware auth mengekstrak `actor_id` + `actor_role`.
2. Middleware RBAC memvalidasi policy per endpoint.
3. Handler memanggil usecase, lalu usecase mencatat audit melalui `RecordAudit`.

### Task List
1. **Domain**
   - Tambah enum/const role pada domain user/employee (jika belum ada).
   - Definisikan struktur event audit (mis. `AuditEvent`).
2. **Usecase**
   - Tambahkan validasi role untuk usecase employee.
   - Buat usecase `RecordAudit` yang menerima metadata aksi + context.
   - Pastikan usecase memanggil audit log untuk create/update/delete.
3. **Adapter Repo**
   - Buat tabel `audit_logs` + migration.
   - Implementasi repo `AuditLogRepository`.
4. **Adapter HTTP**
   - Middleware untuk role check (RBAC).
   - Inject audit context: `actor_id`, `actor_role`, `action`, `resource`, `resource_id`, `ip`, `user_agent`.
5. **DTO**
   - Pastikan response error konsisten.
6. **Testing**
   - Unit test usecase RBAC (allowed/denied).
   - Integration test repo audit log.
   - Handler test: akses staff ke endpoint write harus 403.

---

## Fitur 2: Observability Pack (Prioritas 2)

### Scope Teknis
- **Structured logging** (JSON) dengan `request_id`.
- **Metrics** dasar: request count & latency per endpoint.
- **Health & readiness** endpoint.

### Task List
1. **App/Config**
   - Tambahkan konfigurasi logger (level, format).
   - Generate `request_id` (uuid).
2. **Adapter HTTP**
   - Middleware logging (request/response).
   - Middleware metrics (counter + histogram).
   - Endpoint `/healthz` dan `/readyz`.
3. **Testing**
   - Unit test middleware logging & request_id presence.
   - Smoke test endpoint health/readiness.

---

## Fitur 3: Upload Avatar ke MinIO (Prioritas 3)

### Scope Teknis
- Endpoint upload avatar employee.
- Simpan metadata file (bucket, path, content-type, size).
- Validasi file type & size limit.

### Task List
1. **Usecase**
   - Validasi file + generate object key.
   - Simpan metadata ke DB.
2. **Adapter Repo**
   - Tambah tabel `employee_avatars` + migration.
3. **Adapter HTTP**
   - Endpoint upload (multipart/form-data).
   - Upload file ke MinIO.
4. **Testing**
   - Unit test validation file type/size.
   - Integration test upload flow (mock MinIO).
