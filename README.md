# logslice

Fast log file slicer that extracts time-range segments from large structured log files.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Extract log entries between two timestamps:

```bash
logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" --file app.log
```

Read from stdin and write to a file:

```bash
cat app.log | logslice --from "2024-01-15T08:00:00Z" --to "2024-01-15T09:00:00Z" > slice.log
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--file` | Path to the log file | stdin |
| `--from` | Start of the time range (RFC3339) | required |
| `--to` | End of the time range (RFC3339) | required |
| `--format` | Timestamp format in log entries | `RFC3339` |
| `--field` | Timestamp field name (JSON logs) | `time` |

### Supported Log Formats

- JSON structured logs (e.g., `{"time":"...","level":"info","msg":"..."}`)
- Common plaintext log formats with leading timestamps

## Building from Source

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build ./...
```

## License

MIT © 2024 yourusername