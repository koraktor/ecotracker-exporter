# EcoTracker Exporter

Prometheus exporter for the [local EcoTracker API](https://support.everhome.cloud/en-US/kb/article/36/ecotracker-local-api).

The exporter queries the EcoTracker local API and exposes the resulting power and energy values at `/metrics`. Its default listening port is `9776`.

## Configuration

The exporter requires the following settings:

| Setting | Default | Description |
| --- | --- | --- |
| `HOST` | none | Host name or IP address of the EcoTracker |
| `PORT` | `80` | Port of the EcoTracker local API |
| `LISTEN_ADDRESS` | `:9776` | Address where Prometheus metrics are exposed; `9776` is the default exporter port |
| `LOG_LEVEL` | `warn` | Log level, for example `debug`, `info`, or `error` |

The EcoTracker must be reachable from the container. The exporter requests the `v1/json` endpoint described in the [EcoTracker Local API documentation](https://support.everhome.cloud/en-US/kb/article/36/ecotracker-local-api).

## Installation with Docker

Pull and run the published image from GitHub Container Registry:

```sh
docker run -d \
  --name ecotracker-exporter \
  --restart unless-stopped \
  -p 9776:9776 \
  -e HOST=<ecotracker-address> \
  ghcr.io/koraktor/ecotracker-exporter:latest
```

Replace `<ecotracker-address>` with the hostname or IP address of your EcoTracker. Prometheus can scrape the exporter at:

```text
http://<docker-host>:9776/metrics
```

To build the image locally instead:

```sh
docker build -t ecotracker-exporter .
docker run -d \
  --name ecotracker-exporter \
  --restart unless-stopped \
  -p 9776:9776 \
  -e HOST=<ecotracker-address> \
  ecotracker-exporter
```

## Installation with Docker Compose

Create a `docker-compose.yml` with the following contents:

```yaml
services:
  ecotracker-exporter:
    image: ghcr.io/koraktor/ecotracker-exporter:latest
    container_name: ecotracker-exporter
    restart: unless-stopped
    environment:
      HOST: <ecotracker-address>
      LOG_LEVEL: warn
    ports:
      - "9776:9776"
```

Replace `<ecotracker-address>` with your EcoTracker‘s host name or IP address, then start the
exporter:

```sh
docker compose up -d
```

Check the container logs with:

```sh
docker compose logs -f ecotracker-exporter
```

## Prometheus

Add the exporter as a scrape target in Prometheus:

```yaml
scrape_configs:
  - job_name: ecotracker
    static_configs:
      - targets:
          - <docker-host>:9776
```

## Metrics

All metrics include the `ecotracker_host` label, identifying the EcoTracker queried by the exporter.

| Metric | Type | Description |
| --- | --- | --- |
| `ecotracker_power` | Gauge | Current power consumption or feed-in |
| `ecotracker_power_avg` | Gauge | Average power consumption or feed-in during the last minute |
| `ecotracker_power_last_update` | Counter | Unix timestamp of the last power-data update at the electricity meter |
| `ecotracker_power_phase1` | Gauge | Current power on phase 1 |
| `ecotracker_power_phase2` | Gauge | Current power on phase 2 |
| `ecotracker_power_phase3` | Gauge | Current power on phase 3 |
| `ecotracker_energy_counter_out` | Counter | Total energy counted as feed-in |
| `ecotracker_energy_counter_in` | Counter | Total energy consumption |
| `ecotracker_energy_counter_in_t1` | Counter | Total energy consumption for rate 1 |
| `ecotracker_energy_counter_in_t2` | Counter | Total energy consumption for rate 2 |
| `ecotracker_up` | Gauge | `1` when the EcoTracker API query succeeds, otherwise `0` |
