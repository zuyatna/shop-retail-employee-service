# Scope Teknis & Task List (Portfolio Strategis)

Dokumen ini berisi scope teknis dan task list untuk fitur portofolio strategis yang direkomendasikan.

## Fitur 1: RBAC + Audit Log (Prioritas 1)

### Scope Teknis
- **RBAC** berbasis role: `supervisor`, `staff`.
- **Policy** akses untuk endpoint employee (create/update/delete).
- **Audit log** untuk aksi kritis: create/update/delete employee + login/logout.
- **Standard error response** saat akses ditolak (403) atau token invalid (401).

### Task List
1. **Domain**
   - Tambah enum/const role pada domain user/employee (jika belum ada).
   - Definisikan struktur event audit (mis. `AuditEvent`).
2. **Usecase**
   - Tambahkan validasi role untuk usecase employee.
   - Buat usecase `RecordAudit` yang menerima metadata aksi.
3. **Adapter Repo**
   - Buat tabel `audit_logs` + migration.
   - Implementasi repo `AuditLogRepository`.
4. **Adapter HTTP**
   - Middleware untuk role check (RBAC).
   - Inject `actor_id`, `actor_role`, `action`, `resource` ke audit log.
5. **DTO**
   - Pastikan response error konsisten.
6. **Testing**
   - Unit test usecase RBAC (allowed/denied).
   - Integration test repo audit log.

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
