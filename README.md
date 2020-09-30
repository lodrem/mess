Mess
====
[![Building](https://img.shields.io/github/workflow/status/luncj/mess/Go?style=flat-square)](https://github.com/luncj/mess/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/luncj/mess?style=flat-square)](https://goreportcard.com/report/github.com/luncj/mess)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/luncj/mess)
[![Releases](https://img.shields.io/github/release/luncj/mess/all.svg?style=flat-square)](https://github.com/luncj/mess/releases)
[![LICENSE](https://img.shields.io/github/license/luncj/mess.svg?style=flat-square)](https://github.com/luncj/mess/blob/master/LICENSE)

A toolkit to generate random DML SQL for MySQL.

## Goals
- generate multiple DML(Insert/Update/Delete) SQL by defined schema
- generate for MySQL
- support increment generation(generate DML based on the previous generated data)

## Usage

```shell
./mess generate -h
Generate SQL with random data

Usage:
  mess generate [flags]

Flags:
      --dml string             DML (default "insert")
  -h, --help                   help for generate
      --metadata-path string   Path of metadata for increment generating (default "./metadata.json")
      --num-rows uint          Number of rows (default 1000)
      --output-path string     Path of generated SQL (default "output.sql")
      --schema-path string     Path of schema definition (default "./schema.json")
```

## Examples

- check [Makefile](./Makefile)
- define [schema.json](./examples/schema.json)

Generate `insert` SQL to `./output.sql`
```shell
make run-example DML=insert
```

Generate `update` SQL based on the data inserted before to `./output.sql`
```shell
make run-example DML=update
```

Generate `delete` SQL based on the data before to `./output.sql`
```shell
make run-example DML=delete
```

## License
MIT
