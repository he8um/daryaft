#!/bin/sh
set -eu

PORT="${1:-8091}"
SERVER_DIR="/tmp/daryaft-qa-server"

if ! command -v python3 >/dev/null 2>&1; then
	echo "python3 is required to run the manual QA server" >&2
	exit 1
fi

mkdir -p "$SERVER_DIR"

SERVER_DIR="$SERVER_DIR" PORT="$PORT" python3 - <<'PY'
import os
from pathlib import Path

server_dir = Path(os.environ["SERVER_DIR"])
port = os.environ["PORT"]

(server_dir / "file.txt").write_text("Daryaft manual QA test file\n", encoding="utf-8")

chunk = bytes(range(256)) * 4096
with (server_dir / "big.bin").open("wb") as handle:
    for _ in range(24):
        handle.write(chunk)

(server_dir / "urls.txt").write_text(
    f"http://localhost:{port}/file.txt\n"
    f"http://localhost:{port}/big.bin\n",
    encoding="utf-8",
)
PY

cat <<EOF
Manual QA server files are ready in $SERVER_DIR

Use another terminal for Daryaft commands, for example:
  go run . http://localhost:$PORT/file.txt -o /tmp/daryaft-qa-out
  go run . -f $SERVER_DIR/urls.txt -o /tmp/daryaft-qa-out
  go run . inspect http://localhost:$PORT/file.txt

Starting local HTTP server on http://127.0.0.1:$PORT
Press Ctrl+C to stop it.
EOF

cd "$SERVER_DIR"
exec python3 -m http.server "$PORT" --bind 127.0.0.1
