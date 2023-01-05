# Table: make_data_store

Data Stores are used to store data from scenarios or for transferring data in between individual scenarios or scenario runs.

## Key columns
- Provide a numeric `id` if you want to query for a specific Team.
- Provide a numeric `team_id` to query Data Stores for a specific Team. This can be either set directly in a `where` clause, or specified as part of `join` with another table.

## Examples

### List of all Data Stores in the account

```sql
select
    o.name as organization_name,
    t.name as team_name,
    ds.name as data_store_name
from make_organization o
join make_team t on t.organization_id = o.id
join make_data_store ds on ds.team_id = t.id

-- or a simplified version

select
    ds.id,
    ds.name
from make_data_store ds
```

```
+--------------------+------------------+---------------------+
| organization_name  | team_name        | data_store_name     |
+--------------------+------------------+---------------------+
| Brown Inc.         | Engineering      | Lyon Estates        |
| Brown Inc.         | Engineering      | Lone Pine Mall      |
| Brown Inc.         | Engineering      | Western Auto Stores |
| ...                | ...              | ...                 |
+--------------------+------------------+---------------------+
```

### List of all Data Stores including data quota usage

```sql
select
    o.name as organization_name,
    t.name as team_name,
    ds.name as data_store_name,
    round(100.0 / max_size * size, 2) as ds_fill_perc
from make_organization o
join make_team t on t.organization_id = o.id
join make_data_store ds on ds.team_id = t.id
order by ds_fill_perc desc
```

```
+--------------------+------------------+---------------------+--------------+
| organization_name  | team_name        | data_store_name     | ds_fill_perc |
+--------------------+------------------+---------------------+--------------+
| Brown Inc.         | Engineering      | Lyon Estates        | 60.84        |
| Brown Inc.         | Engineering      | Lone Pine Mall      | 14.31        |
| Brown Inc.         | Engineering      | Western Auto Stores | 8.5          |
| ...                | ...              | ...                 | ...          |
+--------------------+------------------+---------------------+--------------+
```

### Usage of Data Store quota for the whole account

```sql
select
    o.name org_name,
    count(ds.id),
    o.license ->> 'dslimit' dslimit,
    round(100.0 / cast(o.license ->> 'dslimit' as int) * count(ds.id), 2) as usage_perc
from make_data_store ds
join make_team t on t.id = ds.team_id
join make_organization o on o.id = t.organization_id
group by o.name, o.license
order by org_name
```

```
+--------------------+-------+---------+------------+
| org_name           | count | dslimit | usage_perc |
+--------------------+-------+---------+------------+
| Brown Inc.         | 100   | 200     | 50         |
+--------------------+-------+---------+------------+
```

### Detail of a Data Store

```sql
select
    ds.id,
    ds.name,
    ds.records,
    ds.size,
    ds.max_size,
    ds.datastructure_id
from make_data_store ds
where ds.id = 1
```

```
+----+--------------+---------+---------+----------+------------------+
| id | name         | records | size    | max_size | datastructure_id |
+----+--------------+---------+---------+----------+------------------+
| 1  | Lyon Estates | 25448   | 2551815 | 4194304  | <null>           |
+----+--------------+---------+---------+----------+------------------+
```
