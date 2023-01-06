# [Make.com](https://www.make.com/en) Plugin for [Steampipe](https://steampipe.io)

Use SQL to query your Make Scenarios, Connections, Variables, Users, and more.

- Plugin documentation: [Table definitions & examples](https://github.com/marekjalovec/steampipe-plugin-make/tree/main/docs/tables)
- Steampipe Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/marekjalovec/steampipe-plugin-make/issues)

## Quick start

The plugin isn't published yet, follow the **Developing** section for installing the plugin.

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

1] Clone:

```sh
git clone https://github.com/marekjalovec/steampipe-plugin-make
cd steampipe-plugin-make
```

2] Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

3] Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/make.spc
```

4] Increase logging level (optional, but very helpful when debugging)

```
export STEAMPIPE_LOG_LEVEL=INFO
```

5] Try it!

```
steampipe query
> .inspect make
> select c.id, c.name from make_connection c
```

6] Further reading:

- [Make API documentation](https://www.make.com/en/api-documentation)
- [Table definitions & examples](https://github.com/marekjalovec/steampipe-plugin-make/tree/main/docs/tables)
- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [Steampipe contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and their [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md), which apply to this plugin as well. All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-aws/blob/main/LICENSE).
