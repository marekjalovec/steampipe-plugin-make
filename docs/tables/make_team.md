# Table: table_team

Teams are containers that contain scenarios and data accessible only by the members of the team.

## Examples

### Basic Team info

```sql
select
  id,
  name,
  organization_id
from
  make_team;
```

```
+----+--------+-----------------+
| id | name   | organization_id |
+----+--------+-----------------+
| 1  | Team 1 | 1               |
| 2  | Team 2 | 2               |
+----+--------+-----------------+
```
