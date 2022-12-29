# Table: table_organization

Organizations are main containers that contain all teams, scenarios, and users.

## Examples

### Basic Organization info

```sql
select
  id,
  name
from
  make_organization;
```

```
+----+----------------+
| id | name           |
+----+----------------+
| 1  | Organization 1 |
| 2  | Organization 2 |
+----+----------------+
```
