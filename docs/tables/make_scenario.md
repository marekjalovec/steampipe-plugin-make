# Table: make_scenario

Scenarios allow you to create and run automation tasks. A scenario consists of a series of modules that indicate how
data should be transferred and transformed between apps or services.

### Key columns

- Provide a numeric `id` if you want to query for a specific Scenario.

## Examples

### List of all Scenarios

```sql
select
  *
from
  make_scenario;
```

### List all Scenarios using Webhooks as an entry-point

```sql
select
  *
from
  make_scenario
where
  hook_id is not null;
```

### List all Scenarios using Airtable

```sql
select
  *
from
  make_scenario
where
  used_packages::jsonb ? 'airtable';
```

### List of all Scenarios edited today

```sql
select
  *
from
  make_scenario
where
  last_edit > now() - interval '1 days';
```

### List of all Scenarios created by a specific user

```sql
select
  *
from
  make_scenario
where
  created_by_user ->> 'email' = 'marty@mcfly.family';
```

### List of all Scenarios paused due to an error

```sql
select
  *
from
  make_scenario
where
  is_waiting = true;
```

### List of all Scenarios with incomplete executions

```sql
select
  *
from
  make_scenario
where
  dlq_count > 0;
```
