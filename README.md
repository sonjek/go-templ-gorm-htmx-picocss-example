# go-templ-htmx-picocss-example

Example of a Web CRUD app based on [Go](https://github.com/golang/go) + ([sqlite](https://gitlab.com/cznic/sqlite) + [gorm](https://github.com/go-gorm/gorm)) + [templ](https://github.com/a-h/templ) + [</>htmx](https://github.com/bigskysoftware/htmx) + [PicoCSS](https://github.com/picocss/pico) + [Ionicons](https://github.com/ionic-team/ionicons)


![demo](demo.gif)

Features:
- Comfortable and flexible component based templates via [templ](https://github.com/a-h/templ)
- CRUD functionality
- Persistent storage via [SQLite](https://gitlab.com/cznic/sqlite) + ORM ([gorm](https://github.com/go-gorm/gorm))
- Modal windows
- Error handling on user interface side
- Infinite scroll (Lazy Loading)
- User frendly interface
- Preserve static files

## Develop

Available makefile actions:
```sh
% make
Usage:  make COMMAND

Commands:
  build          Compile templ files and build application
  start          Build and start application
  build-docker   Build Docker container image with this app
  run-docker     Run Docker container image with this app
  test           Run unit tests
  tidy           Removes unused dependencies and adds missing ones
  update-deps    Update go dependencies
  get-deps       Download go dependencies
  generate-web   Compile templ files via github.com/a-h/templ/cmd/templ
  air            Build and start application in live reload mode via air
  format         Fix code format issues
  audit          Quality checks
  check-go       Check that Go is installed
  help           Display help
```

## Local Development Setup

To get started, follow these steps:

1) Run `make start` to download go dependencies, compile templ files, build application and finally start application.
```sh
% make start
go mod download
go tool templ generate
(✓) Complete [ updates=7 duration=13.893847ms ]
go run ./cmd/app
Starting web interface on port: 8089
```

You should now be able to access application in your web browser at http://localhost:8089
