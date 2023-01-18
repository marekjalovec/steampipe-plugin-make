# Table: make_connection

For most apps included in Make, it is necessary to create a connection, through which Make will communicate with the given third-party service according to the settings of a specific scenario.

### Key columns
- Provide a numeric `id` if you want to query for a specific Team.
- Provide a numeric `team_id` to query Connections for a specific Team. This can be either set directly in a `where` clause, or specified as part of `join` with another table.

### Caveat
- Be careful when requesting all columns (`*`) or the `scopes` column without using an `id` in the query. To load this detail Steampipe will have to make one extra API request per Connection returned. 

## Examples

### List of all Connections in the account

```sql
select
  o.name as organization_name,
  t.name as team_name,
  c.name as connection_name
from
  make_organization o
  join
    make_team t
    on t.organization_id = o.id
  join
    make_connection c
    on c.team_id = t.id;

-- or a simplified version

select
  id,
  name
from
  make_connection;
```

### List of all Connections for a particular Make App

```sql
-- with Organization and Team name
select
  o.name as organization_name,
  t.name as team_name,
  c.name as connection_name
from
  make_organization o
  join
    make_team t
    on t.organization_id = o.id
  join
    make_connection c
    on c.team_id = t.id
where
  c.account_name = 'amazon-lambda';

-- or a simplified version with own columns only
select
  id,
  name
from
  make_connection
where
  account_name = 'amazon-lambda';
```

### List of all Connections owned by a specific Team

```sql
select
  id,
  name
from
  make_connection
where
  team_id = 1;
```

### List of all Connections owned by a specific User

```sql
select
   t.name as team_name,
   c.id,
   c.name,
   c.account_label,
   c.expire
from
   make_connection c
   join
      make_team t
      on t.id = c.team_id
where
   c.metadata ->> 'type' = 'email'
   and c.metadata ->> 'value' = 'marty@mcfly.family';
```

### Detail of a Connection

```sql
select
  id,
  name,
  account_name,
  account_label,
  account_type,
  package_name,
  expire,
  metadata,
  team_id,
  upgradeable,
  scoped,
  scopes,
  editable,
  uid
from
  make_connection
where
  id = 1;
```
