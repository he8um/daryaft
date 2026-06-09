# HTTP Customization QA Checklist

Manual QA checklist for the v1.3.0 HTTP request customization feature.

Build a local binary first:

```bash
go build -o /tmp/daryaft-qa .
alias dq=/tmp/daryaft-qa
```

---

## 1. Flags visible in help

```bash
dq download --help
dq inspect --help
```

Expected: `--proxy`, `--header`, `--user-agent`, `--username`, `--password` appear in flags list.

---

## 2. Custom header — download

```bash
dq download https://httpbin.org/get --header "X-Test: hello" --dry-run
```

Expected: dry-run succeeds, output shows `X-Test: hello` in HTTP section.

```bash
dq download https://httpbin.org/get --header "X-Test: hello" --output /tmp/qa-header/
cat /tmp/qa-header/get
```

Expected: downloaded JSON contains `"X-Test": "hello"`.

---

## 3. Multiple headers — download

```bash
dq download https://httpbin.org/get \
  --header "X-A: one" \
  --header "X-B: two" \
  --dry-run
```

Expected: both headers shown in dry-run.

---

## 4. User-Agent — download

```bash
dq download https://httpbin.org/user-agent --user-agent "DaryaftTest/1.0" --output /tmp/qa-ua/
cat /tmp/qa-ua/user-agent
```

Expected: `"user-agent": "DaryaftTest/1.0"`.

---

## 5. Proxy — download

```bash
# Requires a local proxy. If mitmproxy or similar is available:
dq download https://httpbin.org/get --proxy http://127.0.0.1:8080 --output /tmp/qa-proxy/ --dry-run
```

Expected: dry-run shows `Proxy: http://127.0.0.1:8080`.

Invalid proxy:

```bash
dq download https://example.com/file.zip --proxy socks5://127.0.0.1:1080 --dry-run
```

Expected: error containing `unsupported scheme "socks5"`.

---

## 6. Basic Auth — download

```bash
dq download https://httpbin.org/basic-auth/alice/pass \
  --username alice \
  --password pass \
  --output /tmp/qa-auth/
```

Expected: download succeeds (200 OK, JSON body).

Password without username:

```bash
dq download https://example.com/file.zip --password secret --dry-run
```

Expected: error `--password requires --username`.

---

## 7. Dry-run redaction

```bash
dq download https://example.com/file.zip \
  --username alice \
  --password topsecret \
  --header "Authorization: Bearer mytoken" \
  --dry-run
```

Expected:
- Output does NOT contain `topsecret`.
- Output does NOT contain `mytoken`.
- Output CONTAINS `[REDACTED]`.

---

## 8. Authorization + Basic Auth conflict

```bash
dq download https://example.com/file.zip \
  --username alice \
  --password pass \
  --header "Authorization: Bearer tok" \
  --dry-run
```

Expected: error about conflicting auth methods.

---

## 9. Invalid header format

```bash
dq download https://example.com/file.zip --header "NoColon" --dry-run
dq inspect https://example.com/file.zip --header "NoColon"
```

Expected: clear error about missing colon.

---

## 10. Inspect with custom header

```bash
dq inspect https://httpbin.org/get --header "X-Inspect: yes"
```

Expected: succeeds and shows metadata. Header is applied on HEAD and GET fallback.

---

## 11. Inspect with Basic Auth

```bash
dq inspect https://httpbin.org/basic-auth/alice/pass \
  --username alice \
  --password pass
```

Expected: metadata for 200 response.

---

## 12. Verbose output — no password leakage

```bash
dq download https://httpbin.org/basic-auth/alice/pass \
  --username alice \
  --password topsecret \
  --verbose \
  --output /tmp/qa-verbose/ 2>&1
```

Expected: output contains `[REDACTED]` for auth, not raw password.

---

## 13. Batch download with custom header

```bash
printf 'https://httpbin.org/get\nhttps://httpbin.org/headers\n' > /tmp/qa-batch-urls.txt
dq download -f /tmp/qa-batch-urls.txt \
  --header "X-Batch: yes" \
  --output /tmp/qa-batch/
```

Expected: both downloads succeed; each downloaded JSON shows `"X-Batch": "yes"`.

---

## 14. TUI not affected

```bash
dq
```

Expected: TUI opens normally. HTTP customization flags are CLI-only; TUI launches without any HTTP options passed.

---

## 15. Environment credential fallback — download

```bash
export DARYAFT_USERNAME=alice
export DARYAFT_PASSWORD=pass
dq download https://httpbin.org/basic-auth/alice/pass --output /tmp/qa-env-auth/
unset DARYAFT_USERNAME DARYAFT_PASSWORD
```

Expected: download succeeds. Credentials come from env, not flags.

CLI flag overrides env:

```bash
export DARYAFT_USERNAME=wrong
export DARYAFT_PASSWORD=wrong
dq download https://httpbin.org/basic-auth/alice/pass \
  --username alice --password pass \
  --output /tmp/qa-env-override/
unset DARYAFT_USERNAME DARYAFT_PASSWORD
```

Expected: download succeeds using the flag values, not env values.

Password env without username:

```bash
export DARYAFT_PASSWORD=secret
dq download https://example.com/file.zip --dry-run
unset DARYAFT_PASSWORD
```

Expected: error `--password requires --username`.

Dry-run redaction with env credentials:

```bash
export DARYAFT_USERNAME=alice
export DARYAFT_PASSWORD=topsecret
dq download https://example.com/file.zip --dry-run
unset DARYAFT_USERNAME DARYAFT_PASSWORD
```

Expected: output contains `[REDACTED]`, not `topsecret`.

---

## 16. Environment credential fallback — inspect

```bash
export DARYAFT_USERNAME=alice
export DARYAFT_PASSWORD=pass
dq inspect https://httpbin.org/basic-auth/alice/pass
unset DARYAFT_USERNAME DARYAFT_PASSWORD
```

Expected: inspect succeeds.

---

## Pass criteria

All numbered cases produce the expected output or error. No raw secrets appear in any output.
