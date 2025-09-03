# hyov-forpg

**Experimental PostgreSQL Terminal Viewer**

## Overview

hyov-forpg is an experimental terminal viewer for PostgreSQL query results. It provides a clean, readable display of PostgreSQL query output in the terminal using modern Go libraries.

## Features

- **PostgreSQL Query Execution**: Execute SQL queries directly from the command line
- **Enhanced Terminal Display**: Results formatted with `termhyo` for better readability
- **Interactive Paging**: Uses `ov` viewer by default for scrolling through large result sets
- **Flexible Output**: Option to output directly to stdout for scripting and piping
- **Modern Performance**: Built with `pgx` driver for optimized PostgreSQL connectivity
- **Graceful Cancellation**: Ctrl+C support for query interruption

## Requirements

- Go 1.23 or later
- PostgreSQL server (local or remote)
- Unix domain socket or network connection

## Installation

```bash
go install github.com/noborus/hyov-forpg@latest
```

## Usage

### Basic Query Execution

```bash
# Execute a query with interactive viewer (default)
hyov-forpg --query "SELECT version();"

# Output directly to stdout (no pager)
hyov-forpg --query "SELECT * FROM users;" --no-pager

# Short form
hyov-forpg -q "SELECT NOW();" -n
```

### Connection Configuration

```bash
# Specify connection string via command line
hyov-forpg -c "host=localhost user=postgres dbname=mydb" -q "SELECT 1;"

# Use configuration file (config.yaml)
hyov-forpg -q "SELECT * FROM table_name;"
```

### Configuration File

Create a `config.yaml` file in your current directory or `~/.config/hyov-forpg/`:

```yaml
db:
  connection: "host=/var/run/postgresql user=myuser dbname=mydb sslmode=disable"
```

## Command Line Options

| Option | Short | Description |
|--------|-------|-------------|
| `--query` | `-q` | SQL query to execute (required) |
| `--connection` | `-c` | Database connection string |
| `--no-pager` | `-n` | Disable pager, output to stdout |
| `--help` | `-h` | Show help message |

## Dependencies

- [pgx](https://github.com/jackc/pgx) - PostgreSQL driver
- [termhyo](https://github.com/noborus/termhyo) - Table formatting
- [ov](https://github.com/noborus/ov) - Terminal pager
- [cobra](https://github.com/spf13/cobra) - CLI framework
- [viper](https://github.com/spf13/viper) - Configuration management

## Example Connection Strings

```bash
# Unix domain socket (default)
"host=/var/run/postgresql user=postgres dbname=postgres sslmode=disable"

# TCP connection
"host=localhost port=5432 user=postgres dbname=mydb password=secret sslmode=require"

# Remote connection
"host=db.example.com port=5432 user=myuser dbname=production sslmode=require"
```

## Development Status

⚠️ **This is an experimental project**

- API and command-line interface may change
- Not recommended for production use
- Feedback and contributions welcome

## Contributing

This is an experimental project. Issues, suggestions, and pull requests are welcome as we explore the best approach for PostgreSQL terminal viewing.

## License

[MIT License](LICENSE) (if applicable)

## Similar Tools

- `psql` - PostgreSQL interactive terminal
- `pgcli` - PostgreSQL CLI with auto-completion
- Various database terminal UIs

---

*Note: This tool is in active development. Features and usage
