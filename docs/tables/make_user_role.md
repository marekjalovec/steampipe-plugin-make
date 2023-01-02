# Table: make_user_role

User roles usable for users within your Organizations and Teams. The column `category` signifies where the role can be used.

## Key columns
- Provide a numeric `id` if you want to query for a specific User Role.

## Examples

### List of all User Roles available in the account

```sql
select
    ur.id,
    ur.name,
    ur.subsidiary,
    ur.category,
    ur.permissions
from make_user_role ur 
```

```
+----+------------------------+------------+--------------+------------------------------+
| id | name                   | subsidiary | category     | permissions                  |
+----+------------------------+------------+--------------+------------------------------+
| 1  | Owner                  | false      | organization | ["access all teams","orga... |
| 2  | Admin                  | true       | organization | ["organization app instal... |
| 3  | App Developer          | true       | organization | ["organization app instal... |
| 4  | Member                 | true       | organization | ["organization variables...  |
| 5  | Team Monitoring        | true       | team         | ["dlq view","scenario log... |
| 6  | Team Admin             | true       | team         | ["account add","account d... |
| 7  | Team Member            | true       | team         | ["account add","account d... |
+----+------------------------+------------+--------------+------------------------------+
```

### List access for a specific User

```sql
-- organization
select distinct 
    mu.name as user_name,
    mo.name as organization_name,
    ur.name as role_name,
    ur.permissions
from make_user_organization_role uor
join make_organization mo on mo.id = uor.organization_id
join make_user mu on mu.id = uor.user_id
join make_user_role ur on ur.id = uor.users_role_id
where uor.user_id = 1

-- team
select distinct
    mu.name as user_name,
    mt.name as team_name,
    ur.name as role_name,
    ur.permissions
from make_user_team_role utr
join make_team mt on mt.id = utr.team_id
join make_user mu on mu.id = utr.user_id
join make_user_role ur on ur.id = utr.users_role_id
where utr.user_id = 1
```

```
+---------------+------------------+-------------+------------------------------+
| user_name     | team_name        | role_name   | permissions                  |
+---------------+------------------+-------------+------------------------------+
| Marty McFLy   | Engineering      | Team Admin  | ["organization app instal... |
| Marty McFLy   | Playground       | Team Admin  | ["organization app instal... |
| Marty McFLy   | Sales            | Team Member | ["account add","account d... |
+---------------+------------------+-------------+------------------------------+
```
