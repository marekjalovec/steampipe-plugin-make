# Table: make_organization

Organizations are main containers that contain all teams, scenarios, and users.

### Key columns
- Provide a numeric `id` if you want to query for a specific Organization.

## Examples

### List of all Organizations in the account

```sql
select
  id,
  name
from
  make_organization
order by
  name
```

### Detail of an Organization

```sql
select
  id,
  name,
  country_id,
  timezone_id,
  license,
  zone,
  service_name,
  is_paused,
  external_id
from
  make_organization
where
  id = 1;
```
