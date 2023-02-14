# Table: make_user_team_role

Assigned User Roles with a Team.

### Key columns

- Provide a numeric `id` if you want to query for a specific User Role assignment.

## Examples

### List access roles for a specific User

```sql
select distinct
  u.name as user_name,
  t.name as team_name,
  ur.name as role_name,
  ur.permissions
from
  make_user_team_role utr
  join
    make_team t
    on t.id = utr.team_id
  join
    make_user u
    on u.id = utr.user_id
  join
    make_user_role ur
    on ur.id = utr.users_role_id
where
  utr.user_id = 1;
```
