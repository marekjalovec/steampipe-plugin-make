# Table: make_team_variable

Team Variables are user-set variables you can use in your scenarios.

## Key columns
Provide a numeric `team_id` to query Team Variables for a specific Team. This can be either set directly in a where clause, or specified as part of join with another table.

## Examples

### List of all custom Team Variables available in the account

```sql
select
    o.name organization_name,
    t.name team_name,
    tv.name variable_name,
    tv.value variable_value
from make_organization o
join make_team t on t.organization_id = o.id
join make_team_variable tv on tv.team_id = t.id
where tv.is_system = false
```

```
+-------------------+-------------+-------------------------+--------------------------------------+
| organization_name | team_name   | variable_name           | variable_value                       |
+-------------------+-------------+-------------------------+--------------------------------------+
| Brown Inc.        | Engineering | token_engineering_prod  | a7770de4-0e76-4a77-a94c-e6323c9fa47e |
| Brown Inc.        | Engineering | token_engineering_stage | ee25e884-ece3-4313-8225-b3cb91d12664 |
| Brown Inc.        | Engineering | token_engineering_test  | a5b62a1e-541a-443f-a779-00fb351c3507 |
| Brown Inc.        | HR          | token_hr_prod           | 0d53720c-be0a-4be6-a8a6-1eb37087bf75 |
+-------------------+-------------+-------------------------+--------------------------------------+
```
