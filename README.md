# cloudsql-proxy-is-up

Built on top of (https://github.com/monzo/envoy-preflight)[https://github.com/monzo/envoy-preflight].

`cloudsql-proxy-is-up` is a simple wrapper application which makes it easier to run applications which depend on Cloud SQL Proxy as a sidecar container for Cloud SQL access. It ensures that your application doesn't start until Cloud SQL Proxy is ready, and that Cloud SQL Proxy shuts down when the application exits. It is best used as a prefix to your existing Docker entrypoint. It executes any argument passed to it, doing a simple path lookup:
```
cloudsql-proxy-is-up echo "hi"
cloudsql-proxy-is-up /bin/ls -a
```

The `cloudsql-proxy-is-up` wrapper won't do anything special unless you provide at least the `CLOUDSQL_PROXY_API` environment variable.  This makes, _e.g._, local development of your app easy.

If you do provide the `CLOUDSQL_PROXY_API` environment variable, `cloudsql-proxy-is-up`
will poll the proxy indefinitely with backoff, waiting for Cloud SQL Proxy to report itself as live. Only then will it execute the command provided as an argument, so that your app can immediately start accessing the Cloud SQL.

All signals are passed to the underlying application. Be warned that `SIGKILL` cannot be passed, so this can leave behind a orphaned process.

When the application exits, as long as it does so with exit code 0, `cloudsql-proxy-is-up` will instruct Cloud SQL Proxy to shut down immediately.

## Environment variables

| Variable              | Purpose                                                                                                                                                                                                                                                                                                                                  |
|-----------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `CLOUDSQL_PROXY_API`     | This is the path to Cloud SQL Proxy port, in the format `http://127.0.0.1:3306`. If provided, `cloudsql-proxy-is-up` will poll this port until it opens. If provided and local (`127.0.0.1` or `localhost`), then Cloud SQL Proxy will be instructed to shut down if the application exits cleanly. |
| `START_WITHOUT_CLOUDSQL_PROXY_API` | If provided and set to `true`, `cloudsql-proxy-is-up` will not wait for Cloud SQL Proxy to be LIVE before starting the main application. However, it will still instruct Cloud SQL Proxy to exit.                                                                                                                                                                 |
