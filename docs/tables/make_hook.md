# Table: make_hook

Webhooks allow you to send data to Make over HTTP. Webhooks create a URL that you can call from an external app or
service, or from another Make scenario. Use webhooks to trigger the execution of scenarios.

### Key columns

- Provide a numeric `id` if you want to query for a specific Hook.

## Examples

### List of all Hooks

```sql
select
  *
from
  make_hook;
```

### List all Hooks not attached to a Scenario

```sql
select
  *
from
  make_hook
where
  scenario_id is null;
```

### List all Hooks with queued executions

```sql
select
  *
from
  make_hook
where
  queue_count > 0;
```

### List all Hooks tied to active scenarios and with dangerously full queues

```sql
select
  h.id,
  h.name,
  h.queue_count,
  h.queue_limit,
  h.team_id,
  s.id as scenario_id,
  s.name as scenario_name
from
  make_hook h
  join
    make_scenario s
    on s.id = h.scenario_id
where
  s.is_enabled = true
  and ((100.0 / h.queue_limit) * h.queue_count) > 80;
```
