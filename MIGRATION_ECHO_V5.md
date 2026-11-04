# Echo v5 migration

This update moves ironman to Echo v5.4.4 and echo-jwt v5.0.2. It is a breaking
source change for applications that pass Echo types to ironman. Update all
application handlers and middleware together; Echo v4 and v5 types cannot be mixed.

## Application changes

- Replace `github.com/labstack/echo/v4` with `github.com/labstack/echo/v5` and
  `github.com/labstack/echo-jwt/v4` with `github.com/labstack/echo-jwt/v5`.
- Use `*echo.Context` in handlers, JWT callbacks and calls to `ironman.JSON`.
- Custom HTTP error handlers now receive `(c *echo.Context, err error)`.
- `Context.Response()` returns `http.ResponseWriter`. Use `echo.UnwrapResponse`
  or `echo.ResolveResponseStatus` for commitment/status inspection.
- Echo uses `log/slog`. `ironman` keeps its Zap access logger through its adapter.
- `ironman.Server` uses `echo.StartConfig`, handles SIGINT/SIGTERM, and allows
  ten seconds for graceful shutdown. Normal shutdown is not logged as an error.

## Deployment checks

Build with Go 1.27.3 or a newer supported, security-patched toolchain. The module
minimum is Go 1.26 because the selected crypto/net dependencies require it.

Echo v5 uses the connection peer for RealIP by default. Behind a reverse proxy,
configure `e.IPExtractor` to trust only that proxy's address ranges and have the
edge proxy overwrite incoming forwarding headers. Verify scheme extraction and
HSTS behind the actual load balancer; do not restore blanket header trust.

Attach authorization directly to static mounts. Check encoded separators, dot
segments and filenames through the actual proxy. Do not rely on a sibling guarded
route to protect files served by a broader static root.

If used by an application, review BodyDump capture limits, Decompress limits,
CSRF lookup validation, JSONP callbacks and WebSocket proxy behavior against the
Echo release notes. Test uploads and routing with real application configuration.

This dependency migration does not change JWT signing keys or token extraction
policy. Applications must replace the public default `ironman` signing key before
production use. `TokenOf` and `ParserTokenUnverified` do not authenticate tokens;
never use their results to grant access.

## Verification

```sh
go mod verify
go build ./...
go vet ./...
go test -race -timeout 120s ./...
govulncheck ./...
```

The existing manual server, MongoDB, remote-logo QR and file-logging examples are
behind the `integration` build tag. They require configuration, external services,
or manual shutdown. Run selected examples explicitly in an isolated environment,
for example `go test -tags=integration -run '^TestMongoDbConnectTest$' .`.
The logger example expects an `app.yml` in its working directory (the supplied
example is named `logger/app.yaml`). They are not covered by ordinary unit tests.

Run vulnerability analysis on each final application as well, with its actual
build tags and deployment platform. A clean library scan does not prove that all
downstream configurations are safe. Rebuild and deploy the final binaries to
deliver dependency and standard-library fixes.
