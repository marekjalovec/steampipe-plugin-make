# Table: make_team

Teams are containers that contain scenarios and data accessible only by the members of the team.

### Key columns
- Provide a numeric `id` if you want to query for a specific Team.
- Provide a numeric `organization_id` to query Teams for a specific Organization. This can be either set directly in a `where` clause, or specified as part of `join` with another table. 

## Examples

### List of all Teams in the account

```sql
-- with Organization name embedded
select
  t.id,
  o.name || ' -> ' || t.name as name
from
  make_organization o
  join
    make_team t
    on t.organization_id = o.id
order by
  name;

-- or a simplified version with own columns only
select
  id,
  name
from
  make_team
order by
  name;
```

### List of all Teams in an Organization

```sql
select
  id,
  name
from
  make_team
where
  organization_id = 1
order by
  name;
```

### Detail of a Team

```sql
select
  id,
  name,
  organization_id
from
  make_team
where
  id = 1;
```
