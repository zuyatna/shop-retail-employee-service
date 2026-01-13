# shop-retail-employee-service
Manage employee shop retail

# Project Structures
```
shop-retail-employee/
├── cmd/
│   └── api/
│       └── main.go                 # entry point (tipis)
│
├── internal/
│   ├── app/                        # composition root
│   │   ├── app.go                  # init & lifecycle
│   │   ├── http.go                 # wiring handler + router
│   │   └── database.go             # DB connection
│   │
│   ├── domain/                     # enterprise business rules
│   │   └── employee.go             # entity murni
│   │
│   ├── usecase/                    # application business rules
│   │   ├── employee_usecase.go
│   │   ├── auth_usecase.go
│   │   └── employee_repository.go  # interface repo
│   │
│   ├── dto/                        # boundary object (API)
│   │   ├── auth/
│   │   │   ├── login_request.go
│   │   │   └── login_response.go
│   │   └── employee/
│   │       ├── create_request.go
│   │       ├── update_request.go
│   │       └── response.go
│   │
│   ├── adapter/                    # interface adapters
│   │   ├── http/
│   │   │   ├── auth_handler.go
│   │   │   ├── employee_handler.go
│   │   │   ├── middleware.go
│   │   │   └── routes.go
│   │   │
│   │   └── repo/
│   │       ├── record/             # persistence model / record
│   │       │   └── employee_record.go
│   │       └── postgres_employee_repo.go
│   │
│   ├── config/                     # configuration
│   │   └── config.go
│   │
│   └── util/                       # technical helper
│       ├── idgen/
│       │   └── uuidv7.go
│       └── jwtutil/
│           └── jwt.go
│
├── migrations/                     # DB schema & seed
│   ├── 001_init.sql
│   └── 002_seed_supervisor.sql
│
├── .env.example
├── go.mod
├── go.sum
└── README.md
```

### Dependency
```
domain
  ↑
usecase
  ↑
adapter (http, repo)
  ↑
app
  ↑
cmd
```

### Data Workflow
```
Handler → Usecase → Domain(Entity + Interface) → Adapter(DB)
```

### Simple Explanation
```
| Folder         | Tanggung Jawab                |
| -------------- | ----------------------------- |
| `cmd`          | Menjalankan aplikasi          |
| `app`          | Wiring & dependency injection |
| `domain`       | Aturan bisnis murni           |
| `usecase`      | Flow bisnis aplikasi          |
| `dto`          | Kontrak API                   |
| `adapter/http` | HTTP delivery                 |
| `adapter/repo` | DB implementation             |
| `repo/model`   | Representasi tabel            |
| `config`       | Konfigurasi                   |
| `util`         | Helper teknis                 |
| `migrations`   | Infrastruktur DB              |
```

### Docker Image
```
docker run --name shop-retail \
  -e POSTGRES_USER=your_user \
  -e POSTGRES_PASSWORD=your_password \
  -e POSTGRES_DB=db_name \
  -p 5432:5432 \
  -d postgres:15
```

### Check Running Docker Container
> sudo docker ps -a

### Enter The Container
> sudo docker exec -it shop-retail psql -U your_user

### Create Database
> CREATE DATABASE db_name;

### List Databases
> \l

### Connect to Database
> \c db_name

### Migration
Copy file migration into container
> sudo docker cp migrations/001_init.sql shop-retail:/001_init.sql

> sudo docker cp migrations/002_seed_supervisor.sql shop-retail:/002_seed_supervisor.sql

Note: You should update the supervisor account password after migration!

```
02_seed_supervisor
email: supervisor@shop.local
password: admin
```

Insert file into database
> sudo docker exec -it shop-retail psql -U your_user -d db_name -f 001_init.sql

> sudo docker exec -it shop-retail psql -U your_user -d db_name -f 002_seed_supervisor.sql

### Enter the Container PSQL
> sudo docker exec -it shop-retail psql -U your_user -d db_name
