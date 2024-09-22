# blacklog-exporter: A tiny log shipping latency exporter

A simple blackbox exporter generating latency metrics from shipped logs

## Introduction

This service consumes logs from a specific Apache Kafka topic. It calculates the latency of the incoming log by subtracting the value of the `timestamp` field of logs from the current time and exposes metrics in Prometheus format.

The generated logs are supposed to be in JSON format with the following schema:

```json
{
    "timestamp": "<RFC 3339 LIKE 2024-09-21T00:17:27>",
    "hostname": "<HOSTNAME OF THE SOURCE OF LOGS>"
}
```

The hostname field helps to inspect unhealthy log sources and is added as a label for exposed metrics.

## Usage

### Local

- Clone the project

  ``` shell
  git clone https://github.com/debMan/blacklog-exporter
  cd blacklog-exporter
  ```

- Ensure you have a running Kafka instance or cluster:

  ``` shell
  # Install Kafka on your own way or run on docker:
  docker compose up -d

  # Create the required topic if it does not exist
  docker compose  exec -it broker kafka-topics --bootstrap-server localhost:9092 --create --topic blacklogs
  ```

- Create the configuration file if required else, leave it with default values else

  ``` shell
  cp config.example.yaml config.yaml
  go run ./cmd/blacklog-exporter
  # go run ./cmd/blacklog-exporter -c <PATH_TO_CONFIG>
  ```

### Container

``` shell
# local build
git clone https://github.com/debMan/blacklog-exporter
docker image build -t blacklog-exporter blacklog-exporter
docker container run -d blacklog-exporter:latest

# or from Docker Hub
docker container run -d idebman/blacklog-exporter
```

## Configuration

Use the `config.yaml` file or environment variables to configure *blacklog-exporter*. The `config.example.yaml` can help you; default values are added. Refer to [Configurations references](#configurations-references).

## Configurations references

The path to the config file can be set by `-c` or `--config` command line argument. The default path to the config file is `./config.yaml`.

| `yaml` key              | type     | env | default              | description                        |
| ----------------------- | -------- | --- | -------------------- | ---------------------------------- |
| logger.level            | string   | `-` | `info`               | Log level                          |
| kafka.auto_commit       | bool     | `-` | `true`               | Kafka auto commit                  |
| kafka.auto_offset_reset | string   | `-` | `earliest`           | Auto offset reset policy for kafka |
| kafka.bootstrap_servers | []string | `-` | `["localhost:9092"]` | TODO                               |
| kafka.debug             | bool     | `-` | `false               | TODO                               |
| kafka.group_id          | string   | `-` | `blacklog-exporter   | TODO                               |
| kafka.sasl_mechanisms   | string   | `-` | `""`                 | TODO                               |
| kafka.sasl_password     | string   | `-` | `""`                 | TODO                               |
| kafka.sasl_username     | string   | `-` | `""`                 | TODO                               |
| kafka.security_protocol | string   | `-` | `""`                 | TODO                               |
| kafka.topics            | []string | `-` | `["blacklogs"]`      | TODO                               |
| metric.enabled          | bool     | `-` | `true`               | TODO                               |
| metric.server.address   | string   | `-` | `":8080"`            | TODO                               |
| metric.server.path      | string   | `-` | `"metrics"`          | TODO                               |

## Missing improvements / TODO

- [ ] Add topic name to metrics
- [ ] Add configurable timeouts for kafka connection
- [ ] Handle commit after consuming if auto-commit is false
- [ ] Add healthchecks for kafka connection
