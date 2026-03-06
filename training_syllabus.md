# NDS Go Starter - 5 Days Intensive Training Syllabus

Silabus ini dirancang khusus untuk mempercepat adaptasi developer junior (terutama yang memiliki latar belakang Spring Boot atau .NET) ke dalam standar pengembangan Go di NDS.

---

## **Day 1: Project DNA & Environment**
*Fokus: Memahami "Kenapa" dan "Bagaimana" proyek ini disusun.*

- **1.1 Pengenalan Arsitektur:**
    - Mengapa Feature-Driven Development?
    - Perbandingan: Go Handler vs Spring @RestController.
    - Pendekatan Manual Dependency Injection (tanpa Magic).
- **1.2 Environment Setup:**
    - Konfigurasi `.env` dan `internal/env`.
    - Dependency management dengan `go mod tidy`.
- **1.3 Struktur Proyek:**
    - Membedah isi `cmd/`, `internal/core/`, dan `internal/features/`.

## **Day 2: Data Persistence & CRUD Flow**
*Fokus: Mengalirkan data dari Database ke API.*

- **2.1 Database Migration:**
    - Version control untuk database menggunakan Goose.
    - Praktik membuat file `.sql` migrasi yang aman.
- **2.2 SQLC (Type-Safe SQL):**
    - Menulis raw query SQL dan men-generate kode Go.
    - Keuntungan SQLC dibanding ORM tradisional.
- **2.3 Repository & Service Pattern:**
    - Implementasi Repository interface.
    - Service layer sebagai pengatur logika bisnis.

## **Day 3: Interaction Layer & Security**
*Fokus: Komunikasi HTTP dan Perlindungan Data.*

- **3.1 HTTP Handlers dengan Chi Router:**
    - Request decoding & Schema validation menggunakan `go-playground/validator`.
    - Standardized JSON Response.
- **3.2 Middleware Power:**
    - Logging, Recoverer, dan Timeout middleware.
    - Implementasi Custom Middleware.
- **3.3 JWT Authentication:**
    - Bagaimana `AuthMiddleware` bekerja.
    - Mengamankan route dengan `r.Group`.

## **Day 4: Quality Assurance & Documentation**
*Fokus: Membuat kode yang reliabel dan mudah dibaca.*

- **4.1 Unit Testing with Testify:**
    - Menulis test case dasar.
    - Assertions (NoError, Equal, Len).
- **4.2 Mocking for Isolation:**
    - Menggunakan `testify/mock` untuk meniru Repository di level Service.
    - Mengapa kita tidak mengetes DB langsung saat unit test?
- **4.3 Auto-Documentation (Swagger):**
    - Menulis Swagger annotations di Handler.
    - Generate JSON/YAML menggunakan `swag init`.

## **Day 5: Advanced Topics & Capstone Project**
*Fokus: Kasus nyata dan tantangan mandiri.*

- **5.1 Database Transactions:**
    - Mengelola atomicity (All or Nothing) menggunakan `WithTx`.
- **5.2 Background Workers:**
    - Membuat service yang berjalan di latar belakang (misal: cleanup session).
- **5.3 Capstone Project (Hands-on):**
    - **Tugas:** Implementasikan fitur "Category" secara mandiri:
        1. Buat tabel kategori.
        2. Buat CRUD (Create & List) menggunakan SQLC.
        3. Tambahkan Unit Test & Mocking.
        4. Munculkan di Swagger.

---

### **Reference Links**
- **Developer Guide:** `developer_guide.md`
- **Go Documentation:** [golang.org/doc](https://golang.org/doc)
- **SQLC Documentation:** [sqlc.dev](https://sqlc.dev)
