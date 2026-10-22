# Ketikode

Online code editor and runner. Monorepo:

| Dir | |
|---|---|
| `ketide/` | Web app (Go, Echo, Postgres) |
| `executor/` | Sandboxed compile/run service (Go, [isolate](https://github.com/ioi/isolate)) |
| `ketikode/` | Early UI prototype |

## Run

```
cp .env.example .env   # set SESSION_SECRET
docker compose up --build
```

Open http://localhost:1324. The executor container runs `--privileged` because isolate needs namespaces and mounts; it is not published on the host.

## Test

```
(cd ketide && go test ./...)
(cd executor && go test ./...)
docker build --target integration -t executor-integration executor && docker run --rm --privileged executor-integration
```
