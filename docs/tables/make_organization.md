# Table: table_organization

Organizations are main containers that contain all teams, scenarios, and users.

## Key columns
- Provide a numeric `id` if you want to query for a specific Organization.

## Examples

### List of all Organizations in the account

```sql
select
    o.id,
    o.name
from make_organization o
order by o.name
```

```
+-----+------------+
| id  | name       |
+-----+------------+
| 1   | Acme Corp. |
| ... | ...        |
+-----+------------+
```

### Detail of an Organization

```sql
select
    o.id,
    o.name,
    o.country_id,
    o.timezone_id,
    o.license,
    o.zone,
    o.service_name,
    o.is_paused,
    o.external_id
from make_organization o
where o.id = 1
```

```
+----+------------+------------+-------------+---------+--------------+--------------+-----------+-------------+
| id | name       | country_id | timezone_id | license | zone         | service_name | is_paused | external_id |   
+----+------------+------------+-------------+---------+--------------+--------------+-----------+-------------+
| 1  | Acme Corp. | 158        | 114         | ...     | eu1.make.com | default      | false     | <null>      |
+----+------------+------------+-------------+---------+--------------+--------------+-----------+-------------+
