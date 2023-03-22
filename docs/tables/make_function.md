# Table: make_function

Functions in Make allow you to transform data. You use functions when mapping data from one module to another. Make offers a variety of built-in functions. On top of that, you can create your own custom functions.

### Key columns
- Provide a numeric `id` if you want to query for a specific Function.
- Provide a numeric `team_id` to query Functions for a specific Team. This can be either set directly in a `where` clause, or specified as part of `join` with another table.

### Caveat
- Be careful when requesting all columns (`*`), `code`, or the `scenarios` column without using an `id` in the query. To load this data, Steampipe will have to make one extra API request per Function returned.

## Examples

### List of all Functions in the account

```sql
select
  f.id,
  f.name,
  f.description,
  t.name as team_name
from
  make_function f
  join
    make_team t
    on t.id = f.team_id;

-- or a simplified version

select
  id,
  name,
  description
from
  make_function;
```

### List of all Functions including arguments and code

```sql
select
    f.id,
    f.name,
    f.description,
    f.args,
    f.code,
    t.name as team_name
from
    make_function f
        join
    make_team t
    on t.id = f.team_id;
```
