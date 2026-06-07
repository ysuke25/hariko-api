# hariko

A lightweight API mock server written in Go. Define your API routes in YAML, get instant mock responses.

The name "hariko" comes from the Japanese word "張り子" (papier-mache) — it looks real on the outside but is hollow inside, just like a mock server.

## Features

- Define routes in a simple YAML file
- Static JSON responses with configurable status codes and headers
- Request logging
- Single binary, zero runtime dependencies

## Installation

### From source

```bash
go install github.com/ysuke25/hariko/cmd/hariko@latest
```

### From releases

Download the latest binary from [Releases](https://github.com/ysuke25/hariko/releases).

## Quick Start

1. Create a `routes.yaml`:

```yaml
server:
  port: 8080

routes:
  - method: GET
    path: /api/hello
    response:
      status: 200
      headers:
        Content-Type: application/json
      body: '{"message": "hello from hariko"}'
```

2. Run the server:

```bash
hariko -f routes.yaml
```

3. Test it:

```bash
curl http://localhost:8080/api/hello
```

## Usage

```
hariko -f <config-file> [-p <port>]
```

| Flag | Description | Default |
|------|-------------|---------|
| `-f` | Path to YAML config file | `routes.yaml` |
| `-p` | Override server port | (from config) |

## Configuration

See [examples/routes.yaml](examples/routes.yaml) for a full example.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

## License

MIT License. See [LICENSE](LICENSE).
