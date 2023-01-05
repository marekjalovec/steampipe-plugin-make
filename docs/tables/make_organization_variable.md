# Table: make_organization_variable

Organization Variables are user-set variables you can use in your scenarios.

## Key columns
Provide a numeric `organization_id` to query Organization Variables for a specific Organization. This can be either set directly in a where clause, or specified as part of join with another table.

## Examples

### List of all custom Organization Variables available in the account

```sql
select 
    o.name organization_name, 
    ov.name variable_name, 
    ov.value variable_value 
from make_organization o 
join make_organization_variable ov on ov.organization_id = o.id 
where ov.is_system = false
```

```
+-------------------+-------------------------+--------------------------------------+
| organization_name | variable_name           | variable_value                       |
+-------------------+-------------------------+--------------------------------------+
| Acme Corp.        | token_engineering_prod  | a7770de4-0e76-4a77-a94c-e6323c9fa47e |
| Acme Corp.        | token_engineering_stage | ee25e884-ece3-4313-8225-b3cb91d12664 |
| Acme Corp.        | token_engineering_test  | a5b62a1e-541a-443f-a779-00fb351c3507 |
| Acme Corp.        | token_hr_prod           | 0d53720c-be0a-4be6-a8a6-1eb37087bf75 |
+-------------------+-------------------------+--------------------------------------+
```
