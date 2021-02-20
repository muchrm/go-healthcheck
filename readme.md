# Healthcheck CLI Tools

## Install

cd go-healthcheck
$ go install -v .

## How to use

1. Create CSV File and add url per each line
2. Run With Command go-healthcheck {csv-path} for example

```bash
go-healthcheck test.csv
```

## Container support

- build your own container with command

```bash
docker build -f Dockerfile -t $(repo):$(commit) .
```

- run with command

```bash
docker run -v $(pwd)/config.json:/go/src/app/config.json -v $(pwd)/test.csv:/go/src/app/test.csv muchrm/go-healthcheck go-healthcheck test.csv
```

## configuration

- MaxWorker : Worker size to be run healthcheck in parallel

### extra configuration

- AccessToken: Bearer Authorization to be set in header request
- HealcheckReportAPI: end point to be send healthcheck result
- Result body represent at below

```json
{
    "total_websites":number,
    "success":number,
    "failure":number,
    "total_time":number //duration in nano second
}
```
