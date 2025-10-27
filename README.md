# shop-retail-employee-service
A replication of employee shop retail using Go


### Docker Image
```
docker run --name shop-retail \
  -e POSTGRES_USER=your_user \
  -e POSTGRES_PASSWORD=your_password \
  -e POSTGRES_DB=employee_db \
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

Insert file into database
> sudo docker exec -it shop-retail psql -U your_user -d employee_db -f /001_init.sql

> sudo docker exec -it shop-retail psql -U your_user -d employee_db -f /002_seed_supervisor.sql
