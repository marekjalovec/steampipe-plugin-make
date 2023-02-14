## v0.2.0 [2023-02-14]

_What's new?_

- Tables

  - [make_scenario](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_scenario.md)
  - [make_scenario_dlq](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_scenario_dlq.md)
  - [make_scenario_log](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_scenario_log.md)

_Enhancements_

  - Added column `is_active` to table `make_api_token`.
  - Added column `sso_pending` to table `make_user_organization_role`.

_Bug fixes_

  - Column `make_data_store`.`datastructure_id` was not correctly loaded without using `id` in the condition. 

## v0.1.0 [2023-01-19]

_What's new?_

- Tables

  - [make_api_token](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_api_token.md)
  - [make_connection](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_connection.md)
  - [make_data_store](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_data_store.md)
  - [make_organization](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_organization.md)
  - [make_organization_variable](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_organization_variable.md)
  - [make_team](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_team.md)
  - [make_team_variable](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_team_variable.md)
  - [make_user](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_user.md)
  - [make_user_organization_role](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_user_organization_role.md)
  - [make_user_role](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_user_role.md)
  - [make_user_team_role](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_user_team_role.md)


- Basic functionality

  - Client-side rate limiting
  - Missing API Token scope suggestions in case an API call fails
