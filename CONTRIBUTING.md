# Contributing to hariko

Thank you for considering contributing to hariko!

## Development Setup

```bash
git clone https://github.com/ysuke25/hariko.git
cd hariko
make build
```

## Running Tests

```bash
make test
```

## Making Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/my-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linter (`make lint`)
6. Commit with conventional commits (`feat:`, `fix:`, `docs:`, etc.)
7. Push and open a Pull Request

## Commit Message Format

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat: add new feature`
- `fix: fix a bug`
- `docs: update documentation`
- `test: add or update tests`
- `refactor: code refactoring`
- `chore: maintenance tasks`

## Code Style

- Run `gofmt` before committing
- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Keep functions small and focused
