# Table: make_connection

For most apps included in Make, it is necessary to create a connection, through which Make will communicate with the given third-party service according to the settings of a specific scenario.

## Key columns
- Provide a numeric `id` if you want to query for a specific Team.
- Provide a numeric `team_id` to query Connections for a specific Team. This can be either set directly in a `where` clause, or specified as part of `join` with another table.

## Caveat
- Be careful when requesting all columns (`*`) or the `scopes` column without using an `id` in the query. To load this detail Steampipe will have to make one extra API request per Connection returned. 

## Examples

### List of all Connections in the account

```sql
select
    o.name as organization_name,
    t.name as team_name,
    c.name as connection_name
from make_organization o
join make_team t on t.organization_id = o.id
join make_connection c on c.team_id = t.id

-- or a simplified version

select
    c.id,
    c.name
from make_connection c
```

```
+-------------------+-------------+-----------------------------+
| organization_name | team_name   | connection_name             |
+-------------------+-------------+-----------------------------+
| Brown Inc.        | Engineering | Airtable | Marty McFly      |
| Brown Inc.        | Engineering | JIRA | Marty McFly          |
| Brown Inc.        | Engineering | GMail | Marty McFly         |
| Brown Inc.        | Engineering | Airtable | Dr. Emmett Brown |
| Brown Inc.        | Engineering | Airtable | George McFly     |
+-------------------+-------------+-----------------------------+
```

### List of all Connections in the account satisfying a condition

```sql
-- with Organization and Team name
select
    o.name as organization_name,
    t.name as team_name,
    c.name as connection_name
from make_organization o
join make_team t on t.organization_id = o.id
join make_connection c on c.team_id = t.id
where c.account_name = 'amazon-lambda'

-- or a simplified version with own columns only
select
    c.id, 
    c.name
from make_connection c
where c.account_name = 'amazon-lambda'
```

```
+-------------------+-------------+-----------------------------+
| organization_name | team_name   | connection_name             |
+-------------------+-------------+-----------------------------+
| Brown Inc.        | Engineering | AWS Lambda | Marty McFly    |
| Brown Inc.        | Engineering | AWS Lambda | Biff Tannen    |
| Brown Inc.        | Engineering | AWS Lambda | Linda McFly    |
| Brown Inc.        | Engineering | AWS Lambda | Dave McFly     |
| Brown Inc.        | Engineering | AWS Lambda | Mr. Strickland |
+-------------------+-------------+-----------------------------+
```

### List of all Connections in a specific Team

```sql
select
    c.id,
    c.name
from make_connection c
where c.team_id = 1
```

```
+----+-----------------------------+
| id | connection_name             |
+----+-----------------------------+
| 1  | Airtable | Marty McFly      |
| 2  | JIRA | Marty McFly          |
| 3  | GMail | Marty McFly         |
| 4  | Airtable | Dr. Emmett Brown |
| 5  | Airtable | George McFly     |
+----+-----------------------------+
```

### Detail of a Connection

```sql
select
    c.id,
    c.name,
    c.account_name,
    c.account_label,
    c.account_type,
    c.package_name,
    c.expire,
    c.metadata,
    c.team_id,
    c.upgradeable,
    c.scoped,
    c.scopes,
    c.editable,
    c.uid
from make_connection c
where c.id = 1
```

```
+----+------------------------+--------------+---------------+--------------+--------------+--------+----------+---------+-------------+--------+--------+----------+-----------+
| id | name                   | account_name | account_label | account_type | package_name | expire | metadata | team_id | upgradeable | scoped | scopes | editable | uid       |
+----+------------------------+--------------+---------------+--------------+--------------+--------+----------+---------+-------------+--------+--------+----------+-----------+
| 1  | Airtable | Marty McFly | airtable     | Airtable      | basic        | <null>       | <null> | ...      | 1       | false       | true   | ...    | false    | 128991824 |
+----+------------------------+--------------+---------------+--------------+--------------+--------+----------+---------+-------------+--------+--------+----------+-----------+
```
