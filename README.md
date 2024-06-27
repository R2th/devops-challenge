# Truong Duc Thang

## Challenge 1:

- solution: https://go.dev/play/p/Tsxly7R6gX_t

## Challenge 2:

- [solution](https://github.com/R2th/devops-challenge/tree/main/chall2)

### Checking the results

Visiting http://localhost:9115/probe?id=1

### Local Build

```bash
make
```

### Building with Docker

```bash
docker build -t exporter .
```

### Local run

```bash
make run
```

### Prometheus Configuration

Example config:

```yaml
scrape_configs:
  - job_name: "exporter"
    metrics_path: /probe
    params:
      id: [1]
    static_configs:
      - targets:
          - http://prometheus.io # Target to probe with http.
          - https://prometheus.io # Target to probe with https.
          - http://example.com:8080 # Target to probe with http on port 8080.
```

## Challenge 3:

- [solution](https://github.com/R2th/devops-challenge/tree/main/chall3)

| Username | Team1      | Team2      | Team3      | Team4      | Team5      | Team6      |
| -------- | ---------- | ---------- | ---------- | ---------- | ---------- | ---------- |
| User1    | member     | maintainer | member     | member     | maintainer | member     |
| User2    | maintainer | member     | member     | maintainer | member     | member     |
| User3    | member     | member     | maintainer | member     | member     | member     |
| User4    |            | member     | member     | member     | member     | maintainer |
| User5    | member     |            | member     | member     |            | member     |
| User6    | member     | member     | member     | member     | member     |            |
| User7    | member     | member     | member     | member     | member     | member     |
| User8    | member     | member     |            | member     | member     | member     |
| User9    | member     | member     |            | member     | member     | member     |
| User10   | member     | member     |            | member     | member     | member     |
