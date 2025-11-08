# shop-retail-employee-service
A replication of employee shop retail using Go

# Project Structures
shop-retail-employee/
├── cmd/
|   └── api/
|       └── main.go
├── internal/
|   ├── adapter/
|   |   ├── http/
|   |   |    ├── auth_handler.go
|   |   |    ├── employee_handler.go
|   |   |    └── middleware.go 
|   |   └── repo/
|   |        └── postgres_employee_repo.go
|   ├── config/
|   |   └── config.go
|   ├── model/
|   |   └── employee.go
|   ├── router/
|   |   └── employee_routes.go
|   ├── usecase/
|   |   ├── auth_usecase.go
|   |   └── employee_usecase.go
|   ├── util/
|   |   ├── idgen/
|   |   |   └── uuidv7.go
|   |   └── jwtutil/
|   |       └── jwt.go
├── migrations/
|   ├── 001_init.sql
|   └── 002_seed_supervisor.sql
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md

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
