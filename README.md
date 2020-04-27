# openttd_exporter - A prometheus exporter for OpenTTD instances

This is a super simple exporter that exports basic information about OpenTTD games (such as player counts).

It uses [gopenttd](https://github.com/ropenttd/gopenttd) for querying.

## Running

For usage information, run the docker image (`redditopenttd/openttd_exporter`) or the binary `go run github.com/ropenttd/openttd_exporter` with the `--help` flag.