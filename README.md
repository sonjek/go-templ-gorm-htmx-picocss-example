# go-full-stack-example

Example of full-stack Go based Web app

[![GO](https://img.shields.io/badge/go-%233366CC.svg?logo=go&logoColor=white)](https://go.dev) [![templ](https://img.shields.io/badge/templ-%233366CC.svg?logo=htmx&logoColor=white)](https://github.com/a-h/templ) [![HTMX](https://img.shields.io/badge/htmx-%233366CC.svg?logo=htmx&logoColor=white)](https://github.com/bigskysoftware/htmx) [![SQLite](https://img.shields.io/badge/sqlite-%233366CC.svg?logo=sqlite&logoColor=white)](https://gitlab.com/cznic/sqlite) [![GORM](https://img.shields.io/badge/gorm-%233366CC.svg?logo=sqlite&logoColor=white)](https://github.com/go-gorm/gorm) [![tailwindcss](https://img.shields.io/badge/tailwindcss-%233366CC.svg?logo=tailwindcss&logoColor=white)](https://github.com/tailwindlabs/tailwindcss) [![daisyui](https://img.shields.io/badge/daisyui-%233366CC.svg?logo=daisyui&logoColor=white)](https://daisyui.com) [![ionicons](https://img.shields.io/badge/ionicons-%233366CC.svg?logo=ionic&logoColor=white)](https://github.com/ionic-team/ionicons)


Features:
- Comfortable and flexible component based templates via [templ](https://github.com/a-h/templ)
- CRUD functionality (Create, Read, Update, and Delete entries)
- Persistent storage via [SQLite](https://gitlab.com/cznic/sqlite) + ORM ([gorm](https://github.com/go-gorm/gorm))
- Error handling on server and user interface side
- Infinite Scrolling via lazy loading
- User friendly interface
- Interactive Modals for better UX
- Native light and dark mode support
- Preserve static files

## Quick start

```bash
# 1. Clone this repository
git clone https://github.com/sonjek/go-full-stack-example && cd go-full-stack-example

# 2. Run (with hot-reload)
make dev

# Or run (without hot-reload)
make start

# Or build a binary and run
make build && bin/app
```

The server starts on `:8089`.

The SQLite database is created automatically and migrations are applied on startup.

---
