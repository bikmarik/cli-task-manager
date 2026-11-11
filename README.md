# CLI Task Manager

A simple command-line task manager written in Go. Create, list, complete, and delete tasks with JSON persistence.

## Features

- Add and delete tasks
- Mark tasks as completed
- List tasks in a formatted table
- Save tasks locally in JSON
- No external dependencies

## Installation

Requires Go installed.

```sh
git clone https://github.com/bikmarik/cli-task-manager.git
cd cli-task-manager
go build -o task .
```

## Usage

```sh
./task add "Learn Go"
./task list
./task complete 1
./task delete 1
./task help
```

Tasks are automatically saved in your user configuration directory.

## Built With

- Go
- Go standard library (`encoding/json`, `os`, `text/tabwriter`)

## License

No license specified.

