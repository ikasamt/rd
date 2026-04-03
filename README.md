# rd - Redmine CLI Tool

A command-line interface tool for Redmine, designed to work seamlessly with Claude Code.

## Installation

```bash
go install github.com/ikasamt/rd@latest
```

## Configuration

You can configure rd via flags, environment variables, or a `.rd` file.

Priority: flags > env vars > `.rd` in CWD > `.rd` at Git root > `~/.rd`.

Environment variables:

```bash
export REDMINE_API_KEY="your-api-key"
export REDMINE_URL="https://your-redmine-instance.com"
```

Optional `.rd` file (key=value):

```
# .rd
REDMINE_URL=https://your-redmine-instance.com
REDMINE_API_KEY=your-api-key
```

You can also use short keys in `.rd`:

```
url=https://your-redmine-instance.com
key=your-api-key
```

## Usage

### List issues

```bash
rd list
rd list --project myproject --status open --assignee me
rd list --oneline
rd list --csv
rd list --json
```

### Get issue details

```bash
rd get 123
rd get 123 --no-comments
rd get 123 --json
```

### Create issue

```bash
rd create --project myproject --title "Bug report" --description "Details here"
rd create --project myproject --title "Task" --tracker 2 --priority 3 --assignee 5
rd create --project myproject --title "With custom field" --field "Field Name=value"
rd create --interactive
```

### Update issue

```bash
rd update 123 --status 2
rd update 123 --assign me
rd update 123 --description "Updated description"
rd update 123 --priority 3 --done-ratio 50
rd update 123 --version "v1.0" --due-date 2025-12-31
rd update 123 --field "Custom Field=value"
rd update 123 --note "Progress update"
```

### Add comment

```bash
rd comment 123 "This is a comment"
```

### Search

```bash
rd search "keyword"
rd search "keyword" --all-types
rd search "keyword" --wiki --titles-only
rd search "keyword" --all --oneline
rd search "keyword" --json
```

### Global flags

```bash
rd --url https://redmine.example.com --key YOUR_API_KEY list
rd --json list
rd --debug get 123
```

## Features

- Simple and intuitive command structure
- Full support for custom fields (name-based resolution)
- JSON output for integration with Claude Code
- Interactive mode for issue creation
- Flexible configuration (flags, env vars, `.rd` file)
- Version name resolution for `--version` flag
- `--assign me` resolves current user automatically
- Search across issues, wiki, news, documents, and more

## License

MIT
