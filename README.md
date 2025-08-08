# 📅 MurimHelper — Schedule & Discipline Assistant

MurimHelper is a personal schedule management backend written in **Go** with a clean architecture (Handler → Usecase → Repository).  
It supports PostgreSQL migrations, recurring events, and a weekly schedule view, and is built for future integration with **natural language input** and **notifications**.

---

## 🚀 Features
- Create, read, update, and delete schedules
- Weekly schedule view
- Mark schedules as completed
- PostgreSQL database with migrations (`golang-migrate`)
- Context-based request handling
- Ready for Natural Language Processing integration (coming soon)

---

## 🛠️ Tech Stack
- **Language:** Go
- **Router:** Gorilla Mux
- **Database:** PostgreSQL
- **Migrations:** golang-migrate
- **Architecture:** Clean architecture (Domain, Repository, Usecase, Handler layers)
- **Container:** Docker (optional, for running migrations)

---
