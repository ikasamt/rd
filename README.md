# rd - Redmine CLI Tool

A command-line interface tool for Redmine, designed to work seamlessly with Claude Code.

## Installation

```bash
go install github.com/ikasamt/rd@latest
```

## Configuration

Set the following environment variables:

```bash
export REDMINE_API_KEY="your-api-key"
export REDMINE_URL="https://your-redmine-instance.com"
```

## Usage

```bash
# List tickets
rd list

# Get ticket details
rd get <ticket-id>

# Create a new ticket
rd create --title "New feature" --description "Description"

# Update ticket status
rd update <ticket-id> --status "In Progress"
```

## Features

- Simple and intuitive command structure
- Full support for custom fields
- JSON output for integration with Claude Code
- Interactive mode for guided operations
- Configuration file support for customization

## License

MIT