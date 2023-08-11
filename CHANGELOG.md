## v0.4.1 [2023-08-11]

_What's new?_

- Dependencies updated to latest versions

## v0.4.0 [2023-04-19]

_What's new?_

- Tables

  - [make_hook](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_hook.md)

_Breaking changes_

- `make_scenario`.`is_linked` renamed to `is_enabled` (following SDK changes).

## v0.3.1 [2023-04-19]

_Dependencies_

  - Recompiled plugin with steampipe-plugin-sdk v5.3.0
  - Dependencies updated to latest versions

## v0.3.0 [2023-03-22]

_What's new?_

- Tables

  - [make_function](https://github.com/marekjalovec/steampipe-plugin-make/blob/master/docs/tables/make_function.md)

_Breaking changes_

- `make_organization`.`licence` updated to match current documentation.
- `make_data_store`.`datastructure_id` renamed to `data_structure_id`.

## v0.2.1 [2023-03-22]

_Bug fixes_

- Fixed selecting of a specific Organization by ID

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
