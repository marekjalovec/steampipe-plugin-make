# Table: make_user_organization_role

Assigned User Roles with an Organization.

### Key columns
- Provide a numeric `id` if you want to query for a specific User Role assignment.

## Examples

### List access roles for a specific User

```sql
select distinct
  u.name as user_name,
  o.name as organization_name,
  ur.name as role_name,
  ur.permissions
from
  make_user_organization_role uor
  join
    make_organization o
    on o.id = uor.organization_id
  join
    make_user u
    on u.id = uor.user_id
  join
    make_user_role ur
    on ur.id = uor.users_role_id
where
  uor.user_id = 1;
```
