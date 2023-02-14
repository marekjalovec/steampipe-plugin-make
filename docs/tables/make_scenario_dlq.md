# Table: make_scenario_dlq

If a scenario terminates unexpectedly because of an error, then the scenario run is discarded. You can set the scenario
to store the failed scenario run as an Incomplete Execution. With that, if an error occurs in your scenario, you can
resolve it manually and avoid losing data.

### Key columns

- Provide a numeric `scenario_id` to query Incomplete Executions for a specific Scenario. This can be either set
  directly in a `where` clause, or specified as part of `join` with another table.

### Caveat

- Be careful when requesting all columns (`*`) or one of the below-mentioned columns without using an `id` in the query.
  To load this data, Steampipe will have to make one extra API request per Incomplete Execution
  returned: `index`, `deleted`, `execution_id`, `scenario_id`, `scenario_name`, `team_id`, `team_name`

## Examples

### List of all Incomplete Executions

```sql
select
  *
from
  make_scenario_dlq
where
  scenario_id = 1;
```

### List all unresolved Incomplete Executions

```sql
select
  id,
  reason,
  attempts,
  created
from
  make_scenario_dlq
where
  scenario_id = 1
  and resolved = false;
```
