![image](https://hub.steampipe.io/images/plugins/marekjalovec/make-social-graphic.png)

# Make plugin for Steampipe

Use SQL to query your Make Scenarios, Connections, Variables, Users, and more.

- **[Get started ->](https://hub.steampipe.io/plugins/marekjalovec/make)**
- Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/marekjalovec/make/tables)
- Community: [Slack Channel](https://steampipe.io/community/join)
- Get involved: [Issues](https://github.com/marekjalovec/steampipe-plugin-make/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):

```shell
steampipe plugin install marekjalovec/make
```

Configure your [credentials](https://hub.steampipe.io/plugins/marekjalovec/make#credentials) and [config file](https://hub.steampipe.io/plugins/marekjalovec/make#configuration).

Run a query:

```sql
select
  distinct id,
  name,
  email,
  last_login
from
  make_user;
```

## Developing

Prerequisites:

- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/marekjalovec/steampipe-plugin-make
cd steampipe-plugin-make
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:

```
make
```

Configure the plugin:

```
cp config/* ~/.steampipe/config
vi ~/.steampipe/config/make.spc
```

Increase logging level (optional, but very helpful when debugging)

```
export STEAMPIPE_LOG_LEVEL=INFO
```

5] Try it!

```
steampipe query
> .inspect make
> select id, name from make_connection
```

6] Further reading:

- [Make API documentation](https://www.make.com/en/api-documentation)
- [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
- [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [Steampipe contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and their [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md), which apply to this plugin as well. All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-aws/blob/main/LICENSE).
