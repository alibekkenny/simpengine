# simpengine 💘
**simpengine** is a fun side-project that helps simps shoot their shot at their simp targets.  
Built with **Go 1.24**, **Postgres 17**, and **MinIO** (soon) for object storage. Runs fully in Docker.
Live demo: simpengine.xyz

---

## 🚀 Features
- Backend written in **Golang**
- Database: **Postgres**
- Object storage: **MinIO** (soon)
- JWT-based auth
- Simple migration system with [golang-migrate](https://github.com/golang-migrate/migrate)

---

## 📦 Setup

### Requirements
- Docker & Docker Compose installed
- Go (for local development)

### Run with Docker Compose
```bash
docker-compose up --build
```
The following services will start:
- Backend: http://localhost:8080
- Postgres: localhost:5433
- MinIO API: http://localhost:9000
- MinIO Console: http://localhost:9001

Apply migrations with golang-migrate:
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/simpengine?sslmode=disable" up
```
