# Table: make_api_token

API tokens of the currently authenticated user.

## Examples

### List of all my API Tokens

```sql
select
  token,
  label,
  scope,
  created
from
  make_api_token;
```

### Find out which API Token is being used without opening `make.spc`

```sql
select
  token,
  label
from
  make_api_token
where
  is_active = true;
```
