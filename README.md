# Programming Blog (Go • Clean Architecture)

Мини-блог на Go с чистой архитектурой: порты, адаптеры, юзкейс, JWT-аутентификация и серверные шаблоны. Лёгкий, быстрый, модульный.

## Стек
Go, Gin, GORM + PostgreSQL, JWT + bcrypt, html/template, SMTP (контакт-форма), Docker-ready.

## Возможности
- CRUD постов и категорий, список/деталь, фильтр по категориям
- Регистрация и вход по JWT
- Контакт-форма (SMTP)
- Веб-UI на Go templates
- Чистые слои: `domain / usecase / adapter / infrastructure`

## Быстрый старт
```bash
git clone https://github.com/100bench/programming_blog
cd programming_blog
go mod tidy

# .env (минимум)
# DB_HOST=localhost
# DB_USER=devuser
# DB_PASSWORD=devpass
# DB_NAME=devsearch_go
# DB_PORT=5432
# JWT_SECRET=supersecret
# SMTP_HOST=smtp.example.com
# SMTP_PORT=587
# SMTP_USERNAME=user
# SMTP_PASSWORD=pass
# SMTP_FROM=noreply@example.com
# PORT=8080

# миграции
psql -h $DB_HOST -U $DB_USER -d $DB_NAME \
  -f internal/adapter/persistence/postgres/migrations/0001_initial.up.sql

# запуск
go run cmd/main.go
# → http://localhost:8080
