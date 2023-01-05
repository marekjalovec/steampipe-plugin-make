# Table: make_user

Existing Users within your account, and their attributes.

## Key columns
- Provide a numeric `id` if you want to query for a specific User.
- Provide a numeric `organization_id`, or `team_id` to query Users for a specific Organization, or a Team. This can be either set directly in a `where` clause, or specified as part of `join` with another table. 

## Caveat
- Thanks to the structure of the Make API, Users can appear in the response multiple times if you query this table. Use `distinct`, or `group by` to get unique records.

## Examples

### List of all Users in the account

```sql
select distinct
    u.id,
    u.name
from make_user u 
order by u.name
```

```
+------+----------------+
| id   | name           |
+------+----------------+
| 1    | Marty McFly    |
| 2    | Biff Tannen    |
| 3    | Linda McFly    |
| 4    | Dave McFly     |
| 5    | Mr. Strickland |
+------+----------------+
```

### List of all Users in an Organization, or a Team

```sql
-- organization
select distinct 
    u.id,
    u.name
from make_user u
where u.organization_id = 1
order by u.name

-- team
select distinct
    u.id,
    u.name
from make_user u
where u.team_id = 1
order by u.name
```

```
+------+----------------+
| id   | name           |
+------+----------------+
| 1    | Marty McFly    |
| 3    | Linda McFly    |
| 4    | Dave McFly     |
+------+----------------+
```
