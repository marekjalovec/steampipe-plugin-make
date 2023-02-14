# Table: make_user

Existing Users within your account, and their attributes.

### Key columns

- Provide a numeric `id` if you want to query for a specific User.
- Provide a numeric `organization_id`, or `team_id` to query Users for a specific Organization, or a Team. This can be
  either set directly in a `where` clause, or specified as part of `join` with another table.

### Caveat

- Thanks to the structure of the Make API, Users can appear in the response multiple times if you query this table.
  Use `distinct`, or `group by` to get unique records.

## Examples

### List of all Users in the account

```sql
select distinct
  id,
  name
from
  make_user
order by
  name;
```

### List of all Users in an Organization, or a Team

```sql
-- organization
select distinct
  id,
  name
from
  make_user
where
  organization_id = 1
order by
  name;

-- team
select distinct
  id,
  name
from
  make_user
where
  team_id = 1
order by
  name;
```

### List of all Users who logged in the last 30 days

```sql
select distinct
  id,
  name,
  last_login
from
  make_user
where
  last_login > now() - interval '30 days'
order by
  name;
```
