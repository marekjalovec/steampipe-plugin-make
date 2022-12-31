# Table: make_connection

For most apps included in Make, it is necessary to create a connection, through which Make will communicate with the given third-party service according to the settings of a specific scenario.

## Examples

### Basic Connection info

```sql
select
  id,
  name
from
  make_connection;
```

```
+----+--------------+
| id | name         |
+----+--------------+
| 1  | Connection 1 |
| 2  | Connection 2 |
+----+--------------+
```
