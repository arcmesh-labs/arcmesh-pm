# Contributing to arcmesh-pm

Thanks for your interest in contributing. `apm` is built by developers for developers — contributions of all kinds are welcome.

---

## Ways to contribute

- **Bug reports** — open an issue with steps to reproduce
- **Feature requests** — open an issue or start a Discussion under Ideas
- **Code** — fix a bug, add a feature, improve test coverage
- **Registry** — add a new MCP server to [arcmesh-registry](https://github.com/arcmesh-labs/arcmesh-registry)
- **Docs** — improve README, fix typos, add examples

---

## Getting started

```bash
git clone https://github.com/arcmesh-labs/arcmesh-pm
cd arcmesh-pm
make build       # builds bin/apm
make test        # runs tests
./bin/apm        # try it locally
```

Requirements: Go 1.21+

---

## Before you write code

1. Check if there's an existing issue or PR for what you want to do.
2. For anything non-trivial, open an issue first so we can align before you invest time.
3. Read `docs/DEVELOPMENT.md` — it covers architecture, conventions, and what not to do.

---

## Making a change

1. Fork the repo and create a branch: `git checkout -b your-branch`
2. Make your change — one thing per PR
3. Run `make test` and make sure everything passes
4. Open a PR with a clear description of what changed and why

There's no strict PR template, but a good description answers:
- What does this change?
- Why is it needed?
- Anything reviewers should pay attention to?

---

## Adding a new MCP server to the registry

That lives in a separate repo: [arcmesh-registry](https://github.com/arcmesh-labs/arcmesh-registry). See its `CONTRIBUTING.md` for the manifest format and validation steps.

---

## Questions

Open a [Discussion](https://github.com/arcmesh-labs/arcmesh-pm/discussions) — that's the right place for questions, ideas, and general conversation.