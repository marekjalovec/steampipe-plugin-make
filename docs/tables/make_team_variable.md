# Table: make_team_variable

Team Variables are user-set variables you can use in your scenarios.

### Key columns
- Provide a numeric `team_id` to query Team Variables for a specific Team. This can be either set directly in a where clause, or specified as part of join with another table.

## Examples

### List of all custom Team Variables available in the account

```sql
select
  o.name organization_name,
  t.name team_name,
  tv.name variable_name,
  tv.value variable_value 
from
  make_organization o 
  join
    make_team t 
    on t.organization_id = o.id 
  join
    make_team_variable tv 
    on tv.team_id = t.id 
where
  tv.is_system = false
```

### List of all Bearer tokens

```sql
select
  o.name organization_name,
  t.name team_name,
  tv.name variable_name,
  tv.value variable_value 
from
  make_organization o 
  join
    make_team t 
    on t.organization_id = o.id 
  join
    make_team_variable tv 
    on tv.team_id = t.id 
where
  tv.is_system = false
  and tv.value like 'Bearer %'
```
