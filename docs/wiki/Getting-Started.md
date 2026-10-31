## 1. Install Go 1.24+

**macOS (Apple Silicon)**

```bash
brew install go
```

**Linux (amd64 or arm64)**

```bash
GO_VERSION=1.24.7
ARCH=$(uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')
curl -fsSLO "https://go.dev/dl/go${GO_VERSION}.linux-${ARCH}.tar.gz"
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-${ARCH}.tar.gz"
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile && source ~/.profile
```

Check with `go version`.

## 2. Install the Temporal CLI

**macOS, or Linux with Homebrew**

```bash
brew install temporal
temporal --version
```

**Linux without Homebrew:** download the archive for your architecture from Temporal's [CLI setup page](https://docs.temporal.io/cli/setup-cli), extract it, and put `temporal` on your `PATH`.

## 3. Clone the repos as siblings

shuttle builds against a local checkout of weft until weft is published, so `go.mod` replaces it with `../weft`.

```bash
mkdir -p ~/src/heddle && cd ~/src/heddle
git clone https://github.com/heedle-ai/weft.git
git clone https://github.com/heedle-ai/shuttle.git
cd shuttle
go test -race ./...
```

The tests use Temporal's in-process test suite, so they don't need a server.

## 4. Run the example against a real server

Terminal 1:

```bash
temporal server start-dev
```

This serves gRPC on `localhost:7233` and the Web UI on http://localhost:8233. It keeps data in memory; add `--db-filename temporal.db` to keep it between restarts.

Terminal 2:

```bash
cd ~/src/heddle/shuttle
go run ./examples/echo
```

Expected output (the workflow ID and run ID on stderr will differ):

```text
workflow 3f1c… run 9a7e…
hello  hello
join   joined: | hello world
world  world
```

`hello` and `world` have no dependencies, so they run in parallel. `join` runs after both and receives their outputs. `examples/echo/main_test.go` checks these exact values.

Open http://localhost:8233 and choose the `PlanWorkflow` execution to see each activity in its event history.

Use `-address host:port` to point the example at a different Temporal frontend.

## Troubleshooting

| Symptom | Cause |
|---|---|
| `connect to Temporal at localhost:7233` | The dev server isn't running, or it's on another port |
| `go: github.com/heedle-ai/weft … no such file` | weft isn't cloned next to shuttle |
| The workflow stays in *Running* | No worker is polling the task queue; the example starts its own, so check it's still running |

Next: [Plans](Plans).
