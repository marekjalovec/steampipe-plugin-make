# Table: make_team

Teams are containers that contain scenarios and data accessible only by the members of the team.

## Key columns
- Provide a numeric `id` if you want to query for a specific Team.
- Provide a numeric `organization_id` to query Teams for a specific Organization. This can be either set directly in a `where` clause, or specified as part of `join` with another table. 

## Examples

### List of all Teams in the account

```sql
-- with Organization name embedded
select
    t.id,
    o.name || ' -> ' || t.name as name
from make_organization o
join make_team t on t.organization_id = o.id
order by name

-- or a simplified version with own columns only
select
    t.id,
    t.name
from make_team t
order by t.name
```

```
+------+---------------------------+
| id   | name                      |
+------+---------------------------+
| 1    | Brown Inc. -> Engineering |
| 2    | Brown Inc. -> HR          |
| 3    | Brown Inc. -> Ops         |
| 4    | Brown Inc. -> Sales       |
| ..   | ...                       |
+------+---------------------------+
```

### List of all Teams in an Organization

```sql
select
    t.id,
    t.name
from make_team t
where t.organization_id = 1
order by t.name
```

```
+------+-------------+
| id   | name        |
+------+-------------+
| 1    | Engineering |
| 2    | HR          |
| 3    | Ops         |
| 4    | Sales       |
| ..   | ...         |
+------+-------------+
```

### Detail of a Team

```sql
select
    t.id,
    t.name,
    t.organization_id
from make_team t
where id = 1
```

```
+----+-------------+-----------------+
| id | name        | organization_id |
+----+-------------+-----------------+
| 1  | Engineering | 1               |
+----+-------------+-----------------+
```
