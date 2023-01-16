# Table: make_user_role

User roles usable for users within your Organizations and Teams. The column `category` signifies where the role can be used.

### Key columns
- Provide a numeric `id` if you want to query for a specific User Role.

## Examples

### List of all User Roles available in the account

```sql
select
  id,
  name,
  subsidiary,
  category,
  permissions 
from
  make_user_role
```

### List of all Admins in the account

```sql
select distinct
   u.name as user_name,
   o.name as organization_name 
from
   make_user u 
   inner join
      make_user_organization_role uor 
      on uor.user_id = u.id 
   join
      make_organization o 
      on o.id = uor.organization_id 
where
   uor.users_role_id = 12
```

### List access for a specific User

```sql
-- organization
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
  uor.user_id = 1

-- team
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
  utr.user_id = 1
```
