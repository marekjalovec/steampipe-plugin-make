# Table: make_scenario_log

Scenario Logs allow you to explore past runs of your scenarios.

### Key columns

- Provide a numeric `scenario_id` to query Scenario Logs for a specific Scenario. This can be either set directly in
  a `where` clause, or specified as part of `join` with another table.
- Provide both a string `id` and a numeric `scenario_id` if you want to query for a specific Scenario Log.

## Examples

### List of all Scenario Logs

```sql
select
  *
from
  make_scenario_log
where
  scenario_id = 1;
```

### List all non-auto scenario runs (errors, manual, etc.)

```sql
select
  *
from
  make_scenario_log
where
  scenario_id = 1
  and type != 'auto';
```
