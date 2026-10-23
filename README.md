# Ketikode

Online code editor and runner. Monorepo:

| Dir | |
|---|---|
| `backend/` | Web app (Go, Echo, Postgres) |
| `executor/` | Sandboxed compile/run service (Go, [isolate](https://github.com/ioi/isolate)) |
| `frontend/` | Early UI prototype |

## Run

```
cp .env.example .env   # set SESSION_SECRET (32+ chars: openssl rand -hex 32)
docker compose up --build
```

Open http://localhost:1324. Set `ADMIN_USER` and `ADMIN_PASSWORD` to enable `/admin` (submissions list, basic auth); serve over HTTPS and set `COOKIE_SECURE=1` in production. The executor container runs `--privileged` because isolate needs namespaces and mounts; it is not published on the host.

## Test

```
(cd backend && go test ./...)
(cd executor && go test ./...)
docker build --target integration -t executor-integration executor && docker run --rm --privileged executor-integration
```
