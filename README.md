# 🚀 Finance Tracker - Where you follow your money

***Finance Tracker** is a modern, full-stack web application for personal finance tracking. Built with Go, Gin, HTMX, and PostgreSQL, it offers robust authentication (JWT & PASETO) and a reactive user experience. The project leverages sqlc for type-safe database access and Air for rapid development.*


## ✨ Features

- 🔐 **User Authentication**
  Secure login & registration via JWT and PASETO tokens.

- 💶 **Finance Tracking**
  Track expenses, incomes, and categorize transactions.

- 📦 **Modern UI**
  Interactive frontend powered by HTMX for seamless user experience.

- 🗄️ **Type-safe Database Access**
  Uses sqlc to generate Go code from SQL queries.

- 💨 **Rapid Development**
  Hot reloading with Air for efficient development workflow.


## ⚒️ Tech Stack

- **Backend:** Go, Gin
- **Frontend:** HTMX, HTML/CSS
- **Database:** PostgreSQL
- **ORM/DB Access:** sqlc
- **Authentication:** JWT, PASETO
- **Dev Tools:** Air (live reload)

## ⚙️ Getting Started

### Prerequisites

- Go 1.20+
- PostgreSQL 13+
- [sqlc](https://sqlc.dev/)
- [Air](https://github.com/cosmtrek/air)

### 🛠 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/finance_tracker.git
   cd finance_tracker
   ```

2. **Setup the database**
   ```bash
   make postgres && make createdb && make migrationup
   ```

3. **Run the application (with Air for hot reload)**
   ```bash
   make air
   ```

4. **Run the application directly**
   ```bash
   make sqlc && make templ && make server
   ```


## 🎮 Usage

- Visit `http://localhost:8080` in your browser.
- Register a new account or log in.
- Start tracking your expenses and incomes.


## 📚 Acknowledgements
- [Gin](https://github.com/gin-gonic/gin)
- [HTMX](https://htmx.org/)
- [sqlc](https://sqlc.dev/)
- [Air](https://github.com/cosmtrek/air)
- [PASETO](https://paseto.io/)

## 🤝 Contributing
Pull requests are welcome! For major changes, please open an issue first to discuss what you would like to change.

## 📄 License
MIT License - See [LICENSE](LICENSE)

Happy tracking! 🚀

---

**Developed with ❤️ by [Moth13] | 2025**
[Full documentation](docs/) | [Advanced examples](examples/)
