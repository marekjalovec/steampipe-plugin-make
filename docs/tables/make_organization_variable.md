# Table: make_organization_variable

Organization Variables are user-set variables you can use in your scenarios.

### Key columns
- Provide a numeric `organization_id` to query Organization Variables for a specific Organization. This can be either set directly in a where clause, or specified as part of join with another table.

## Examples

### List of all custom Organization Variables available in the account

```sql
select
  o.name organization_name,
  ov.name variable_name,
  ov.value variable_value
from
  make_organization o
  join
    make_organization_variable ov
    on ov.organization_id = o.id
where
  ov.is_system = false;
```

### List of all API endpoints

```sql
select
  o.name organization_name,
  ov.name variable_name,
  ov.value variable_value 
from
  make_organization o 
  join
    make_organization_variable ov 
    on ov.organization_id = o.id 
where
  ov.is_system = false
  and ov.value like '%api%';
```
