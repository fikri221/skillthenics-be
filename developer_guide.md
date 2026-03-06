# NDS Go Starter - Developer Guide

Selamat datang di proyek **NDS Go Starter**. Panduan ini dirancang untuk membantu developer (terutama junior) dalam memahami struktur, konfigurasi, dan cara berkontribusi pada codebase ini.

---

## 1. Konfigurasi & Variabel Lingkungan

Proyek ini menggunakan variabel lingkungan (environment variables) untuk konfigurasi. File `.env` di root direktori digunakan untuk pengembangan lokal.

### Daftar Konfigurasi Utama

| Kunci            | Kegunaan                                   | Contoh Nilai                           |
| :--------------- | :----------------------------------------- | :------------------------------------- |
| `GOOSE_DRIVER`   | Driver database untuk migrasi              | `mysql`                                |
| `GOOSE_DBSTRING` | Koneksi string database (DSN)              | `root:pass@tcp(localhost:3306)/dbname` |
| `AUTH_ENABLED`   | Mengaktifkan/menonaktifkan validasi JWT    | `true` / `false`                       |
| `JWT_SECRET`     | Kunci rahasia untuk tanda tangan token JWT | `your-secret-key`                      |
| `LOG_FILE_PATH`  | Jalur file untuk output log internal       | `logs/app.log`                         |
| `CORS_ENABLED`   | Mengaktifkan kebijakan CORS                | `true`                                 |

---

## 2. Arsitektur Proyek (Feature-Driven)

Codebase ini mengikuti pola **Feature-Driven Development**. Semua logika terkait satu fitur dikelompokkan dalam satu paket di `internal/features`.

### Struktur Folder Fitur (Contoh: `products`)

- `domain.go`: Definisi struktur data (struct) yang digunakan di level aplikasi.
- `repository.go`: Definisi interface database dan implementasinya (biasanya membungkus SQLC).
- `service.go`: Tempat logika bisnis utama (business logic layer).
- `handler.go`: Menangani HTTP request, validasi input, dan Swagger documentation.
- `register.go`: Tempat mendaftarkan route ke router Chi dan melakukan Dependency Injection.

---

## 3. Membuat API CRUD: Step-by-Step

Ikuti langkah-langkah berikut untuk membuat fitur baru:

### Langkah 1: Database Migration

1. Buat file migrasi baru di `internal/adapters/mysql/migrations`.
2. Gunakan format `YYYYMMDDHHMMSS_name.sql`.
3. Atau bisa menggunakan command dari goose
   ```bash
   goose create <nama_migrasi> sql
   ```
4. Jalankan migrasi menggunakan Goose:
   ```bash
   goose -dir internal/adapters/mysql/migrations mysql "user:pass@tcp(host:port)/dbname" up
   ```

### Langkah 2: SQLC Queries

1. Tambahkan query SQL di folder `internal/core/repository/queries` (misal: `categories.sql`).
2. Jalankan perintah generate:
   ```bash
   sqlc generate
   ```

### Langkah 3: Domain & Repository

1. Definisikan entity di `domain.go`.
2. Buat interface `Repository` di `repository.go` dan implementasikan metodenya menggunakan SQLC querier.

### Langkah 4: Service Layer

1. Buat interface `Service` di `service.go`.
2. Implementasikan logika bisnis di service tersebut. Gunakan Repository untuk akses data.

### Langkah 5: Handler & Swagger

1. Buat handler di `handler.go`.
2. Tambahkan komentar Swagger di atas fungsi handler (Summary, Tags, Param, Success, Router).
3. Gunakan `json.DecodeAndValidate` untuk validasi input.
4. Perbarui dokumentasi Swagger:
   ```bash
   swag init -g cmd/api.go -o docs
   ```

### Langkah 6: Register Route

1. Daftarkan handler ke router di `register.go`.
2. Panggil fungsi `Register` fitur tersebut di `cmd/api.go`.

---

## 4. Training Course: Fitur Categories

Berikut adalah implementasi lengkap untuk fitur **Category**.

### `domain.go`

```go
type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
```

### `repository.go`

```go
type Repository interface {
	CreateCategory(ctx context.Context, id, name string) error
	ListCategories(ctx context.Context) ([]Category, error)
}

type repo struct {
    db *repository.Queries
}

func (r *repo) CreateCategory(ctx context.Context, id, name string) error {
    return r.db.CreateCategory(ctx, repository.CreateCategoryParams{
        ID:   id,
        Name: name,
    })
}
```

### `service.go`

```go
type Service interface {
    CreateCategory(ctx context.Context, name string) error
}

type svc struct {
    repo Repository
}

func (s *svc) CreateCategory(ctx context.Context, name string) error {
    id := ksuid.New().String()
    return s.repo.CreateCategory(ctx, id, name)
}
```

### `handler.go` (Swagger & Validation)

```go
type createReq struct {
    Name string `json:"name" validate:"required,min=3"`
}

// CreateCategory godoc
// @Summary      Create a Category
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        request  body      createReq  true  "Category info"
// @Success      201      {object}  json.Response
// @Router       /categories [post]
func (h *handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
    var req createReq
    if err := json.DecodeAndValidate(r, &req); err != nil {
        json.WriteError(w, r, err)
        return
    }
    // ... call service ...
    json.Write(w, r, http.StatusCreated, "Created")
}
```

### `service_test.go` (Mocking & Unit Testing)

```go
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateCategory(ctx context.Context, id, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func TestCreateCategory(t *testing.T) {
	mockRepo := new(MockRepository)
	service := NewService(mockRepo)

	mockRepo.On("CreateCategory", mock.Anything, mock.Anything, "Electronics").Return(nil)

	err := service.CreateCategory(context.Background(), "Electronics")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
```

### Menjalankan Unit Test

Gunakan perintah berikut untuk menjalankan test di satu package:

```bash
go test ./internal/features/categories/...
```

---

## 5. Membuat Worker (Background Service)

Worker digunakan untuk tugas yang berjalan di background secara periodik.

### Langkah 1: Implementasi Interface

Buat struct yang mengimplementasikan interface `worker.Worker` yang ada di `internal/core/worker`.

```go
type myWorker struct {
    interval time.Duration
}

func (w *myWorker) Start(ctx context.Context) {
    ticker := time.NewTicker(w.interval)
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            // Lakukan tugas background di sini
            fmt.Println("Worker running...")
        }
    }
}
```

### Langkah 2: Registrasi di `cmd/api.go`

Tambahkan instance worker Anda ke aplikasi:

```go
w := features.NewMyWorker(10 * time.Minute)
app.addWorker(w)
```

## 6. Integrasi Kafka

Proyek ini menggunakan `segmentio/kafka-go` untuk messaging.

### Konsep Producer (Kirim Pesan)

Producer diletakkan di dalam Service. Contoh pada `internal/features/orders/service.go`:

1. Tambahkan `*kafka.Writer` ke struct service.
2. Gunakan `s.kafkaWriter.WriteMessages` untuk mengirim event.

### Konsep Consumer (Terima Pesan)

Consumer diimplementasikan sebagai `worker.Worker`. Contoh pada `internal/features/notifications/worker.go`:

1. Gunakan `*kafka.Reader` untuk membaca pesan.
2. Jalankan loop `ReadMessage` di dalam fungsi `Start`.

### Registrasi di `cmd/api.go`

```go
kafkaAddr := "localhost:9092"
writer := kafkaAdapter.NewWriter(kafkaAddr, "topic-name")
reader := kafkaAdapter.NewReader(kafkaAddr, "topic-name", "group-id")

// Tambahkan worker ke app
app.addWorker(notifications.NewNotificationWorker(reader))
```

---

# Pemahaman Konsep dengan framework lain

| Konsep             | Spring Boot       | .NET             | Go                     |
| ------------------ | ----------------- | ---------------- | ---------------------- |
| Controller         | @RestController   | ControllerBase   | Handler struct         |
| Routing            | Annotation        | Attribute        | Router config          |
| DI                 | @Autowired        | Constructor DI   | Manual inject          |
| Middleware         | Filter            | Middleware       | Middleware             |
| Validation         | @Valid            | ModelState       | Manual / validator lib |
| Exception handling | @ControllerAdvice | Exception filter | Manual wrapper         |

---

## 7. Pemahaman Metode Database (SQLC & database/sql)

Berikut adalah ringkasan kegunaan metode-metode database yang sering digunakan dalam repositori (terutama pada kode hasil generate SQLC):

| Method              | Kegunaan Utama                                                              | Tipikal Operasi SQL          | Hasil Kembalian (Return) |
| :------------------ | :-------------------------------------------------------------------------- | :--------------------------- | :----------------------- |
| **ExecContext**     | Menjalankan perintah yang tidak mengembalikan baris data.                   | `INSERT`, `UPDATE`, `DELETE` | `(sql.Result, error)`    |
| **QueryContext**    | Mengambil banyak baris data dari database.                                  | `SELECT` (banyak baris)      | `(*sql.Rows, error)`     |
| **QueryRowContext** | Mengambil tepat satu baris data.                                            | `SELECT` (single row/ID)     | `*sql.Row`               |
| **PrepareContext**  | Menyiapkan template query (Prepared Statement) untuk keamanan dan performa. | Semua jenis query            | `(*sql.Stmt, error)`     |

> [!TIP]
> Selalu gunakan versi `...Context` (misal: `ExecContext` bukan `Exec`) agar aplikasi dapat menangani **timeout** dan **cancellation** dengan baik melalui object `context.Context`.
