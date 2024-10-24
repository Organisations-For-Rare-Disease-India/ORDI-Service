# Project ORDI

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### prerequisites

#### Templ binary

Templ is used for HTML templating in go -  [link](https://templ.guide/).

Use the command to setup `go install github.com/a-h/templ/cmd/templ@latest`.

#### Tailwind css

Follow the installation setup provided [here](https://tailwindcss.com/docs/installation).

#### .env file setup

Refer the `.env.example` file and setup .env file.

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```
