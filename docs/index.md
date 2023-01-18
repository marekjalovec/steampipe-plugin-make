---
organization: marekjalovec
category: ["saas"]
icon_url: "/images/plugins/marekjalovec/make.svg"
brand_color: "#6D01CC"
display_name: "Make"
short_name: "make"
description: "Make plugin for exploring your automations in depth."
og_description: "Query Make with SQL! Open source CLI. No DB required."
og_image: "/images/plugins/marekjalovec/make-social-graphic.png"
---

# Make + Steampipe

[Make](https://www.make.com) allows you to work with everything from tasks and workflows to apps and systems, to build and automate anything in one powerful visual platform.

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

For example:

```sql
select distinct
  id,
  name,
  email,
  last_login
from
  make_user;
```

```
+------+----------------+-----------------------------+---------------------------+
| id   | name           | email                       | last_login                |
+------+----------------+-----------------------------+---------------------------+
| 1    | Marty McFly    | marty@mcfly.family          | 2015-10-21T08:13:06+01:00 |
| 2    | Biff Tannen    | biff.tannen@maddog.rich     | 2015-10-21T08:22:03+01:00 |
| 3    | Linda McFly    | linda@mcfly.family          | 2015-10-21T09:47:13+01:00 |
| 4    | Dave McFly     | dave@mcfly.family           | 2015-10-21T10:51:23+01:00 |
| 5    | Mr. Strickland | director@hillvalleyhigh.edu | 1983-05-22T13:37:00+01:00 |
+------+----------------+-----------------------------+---------------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/marekjalovec/make/tables)**

## Get started

### Install

Download and install the latest Make.com plugin:

```bash
steampipe plugin install marekjalovec/make
```

### Credentials

| Item | Description |
| - | - |
| Credentials | [Get your API token](https://www.make.com/en/api-documentation/authentication-token). |
| Resolution | Credentials explicitly set in a steampipe config file (`~/.steampipe/config/make.spc`). |

### Configuration

Installing the latest Make.com plugin will create a config file (`~/.steampipe/config/make.spc`) with a single connection named `make`:

```hcl
connection "make" {
  plugin = "marekjalovec/make"

  # Make API token
  # To generate the token visit the API tab in your Profile page in Make.
  # Minimum scopes are connections:read and teams:read, as these objects cascade to other API calls.
  # If you add user:read, the plugin can give you useful hints about missing scopes.
  # Add anything you want to have access to on top of that.
  api_token = "ecc9f531-xxxx-xxxx-xxxx-e6de6b3ebbe2"

  # Environment URL
  # The environment of Make you work in. This can be the link to your private instance of Make,
  # for example, https://development.make.cloud, or the link to Make
  # with or without the zone, depending on a specific endpoint, for example, https://eu1.make.com.
  environment_url = "https://eu1.make.com"

  # Rate limiting
  # Make API limits the number of requests you can send to the Make API.
  # Make sets the rate limits based on your organization plan:
  # - Core: 60 per minute
  # - Pro: 120 per minute
  # - Teams: 240 per minute
  # - Enterprise: 1 000 per minute
  # We recommend to set a value below (or at most at) 80% of your total limit.
  # The default value is 50 if you don't override it here.
  # rate_limit = 50
}
```

## Get involved

- Open source: https://github.com/marekjalovec/steampipe-plugin-make
- Community: [Slack Channel](https://steampipe.io/community/join)
