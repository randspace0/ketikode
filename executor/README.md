# Ketikode Executor

HTTP service that compiles and runs untrusted code inside [isolate](https://github.com/ioi/isolate).

## API

`POST /api/exec` (form fields: `code`, `language`, `stdin`, `timeout` seconds, `memory` MB)

Returns `phase` (`compile`|`run`), `stdout`, `stderr`, `success`, `status` (`OK`|`RE`|`SG`|`TO`), `exit_code`, `time`, `memory`.
Compile errors return `phase: "compile"` with the compiler output in `stderr`.

`GET /healthz` returns 200.

Languages: `c`, `cpp`, `python3`, `ruby`, `rust`, `java` (see `languages.go`).

## Configuration

| Env | Default | |
|---|---|---|
| `ADDR` | `:8001` | listen address |
| `ISOLATE_BIN` | `isolate` | isolate binary |
| `ISOLATE_BOXES` | `8` | max concurrent runs |
| `ISOLATE_CGROUPS` | unset | `1` enables `--cg` (exact memory limits, needs cgroup v2 set up for isolate) |

Without cgroups, memory is limited by address space (`--mem`), and by `-Xmx` for Java.

## Tests

```
go test ./...
EXECUTOR_INTEGRATION=1 go test ./...   # needs isolate and root
```

## License

[Mozilla Public License Version 2.0](LICENSE)

## Author

- **Andra Antariksa**
