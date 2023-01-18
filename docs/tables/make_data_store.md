# Table: make_data_store

Data Stores are used to store data from scenarios or for transferring data in between individual scenarios or scenario runs.

### Key columns
- Provide a numeric `id` if you want to query for a specific Team.
- Provide a numeric `team_id` to query Data Stores for a specific Team. This can be either set directly in a `where` clause, or specified as part of `join` with another table.

## Examples

### List of all Data Stores in the account

```sql
select
  o.name as organization_name,
  t.name as team_name,
  ds.name as data_store_name
from
  make_organization o
  join
    make_team t
    on t.organization_id = o.id
  join
    make_data_store ds
    on ds.team_id = t.id

-- or a simplified version

select
  id,
  name
from
  make_data_store;
```

### List of all Data Stores including data quota usage

```sql
select
  o.name as organization_name,
  t.name as team_name,
  ds.name as data_store_name,
  round(100.0 / max_size * size, 2) as ds_fill_perc
from
  make_organization o
  join
    make_team t
    on t.organization_id = o.id
  join
    make_data_store ds
    on ds.team_id = t.id
order by
  ds_fill_perc desc;
```

### Usage of Data Store quota for the whole account

```sql
select
  o.name org_name,
  count(ds.id),
  o.license ->> 'dslimit' dslimit,
  round(100.0 / cast(o.license ->> 'dslimit' as int) * count(ds.id), 2) as usage_perc 
from
  make_data_store ds
  join
    make_team t
    on t.id = ds.team_id
  join
    make_organization o
    on o.id = t.organization_id
group by
  o.name,
  o.license
order by
  org_name;
```

### Detail of a Data Store

```sql
select
  id,
  name,
  records,
  size,
  max_size,
  datastructure_id
from
  make_data_store
where
  id = 1;
```
