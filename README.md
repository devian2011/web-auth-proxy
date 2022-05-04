# Authorization proxy server

[**CHANGELOG**](./changelog.md)
[**TODO**](./TODO.md)

## Configuration

### General config

> Configuration example (all available options)

```json
{
  "logs": {
    "file": {
      "path": "./logs",
      "log_level": [
        "info",
        "warning",
        "error",
        "critical"
      ]
    },
    "console": {
      "level": [
        "debug",
        "info",
        "warning",
        "error",
        "critical",
        "fatal"
      ]
    }
  },
  "auth": {
    "providers": [
      {
        "type": "db",
        "code": "base_db_auth",
        "is_active": true,
        "driver": "postgres",
        "dsn": "host=$host port=$port user=$user password=$password dbname=$db sslmode=disable"
      }
    ]
  },
  "proxy": {
    "port": "8080",
    "points_dir": "points",
    "tls": {
      "cert_file": "./config/certs/proxy.rsa.crt",
      "key_file": "./config/certs/proxy.rsa.key"
    },
    "points": []
  },
  "admin": {
    "port": "8081",
    "tls": {
      "cert_file": "./config/certs/admin.rsa.crt",
      "key_file": "./config/certs/admin.rsa.key"
    }
  }
}
```

### Point config

> Point example config (all available options)

```json
{
  "code": "google_proxy",
  "token_header": "X-User-Token",
  "match": {
    "host": "^.*$",
    "path": "^\/google\/.*$"
  },
  "destination": "https://google.com",
  "providers": [],
  "proxy_headers": [
    "*"
  ],
  "proxy_appended_data_header_name": "PR-USER-DATA",
  "is_header_crypt": false,
  "header_crypt_master_pass": "1234567890",
  "secret": "1312412314123"
}
```

