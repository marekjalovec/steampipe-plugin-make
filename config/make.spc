connection "make" {
  plugin = "marekjalovec/make"

  # Make API token
  # To generate the token visit the API tab in your Profile page in Make.
  # Minimum scopes are connections:read and teams:read, as these objects cascade to other API calls.
  # If you add user:read, the plugin can give you useful hints about missing scopes.
  # Add anything you want to have access to on top of that.
  api_token = "ecc9f531-xxxx-xxxx-xxxx-e6de6b3ebbe2"

  # Environment URL
  # The environment of Make you work in. This can be the link to your private instance of Make,
  # for example, https://development.make.cloud, or the link to Make
  # with or without the zone, depending on a specific endpoint, for example, https://eu1.make.com.
  environment_url = "https://eu1.make.com"

  # Rate limiting
  # Make API limits the number of requests you can send to the Make API.
  # Make sets the rate limits based on your organization plan:
  # - Core: 60 per minute
  # - Pro: 120 per minute
  # - Teams: 240 per minute
  # - Enterprise: 1 000 per minute
  # We recommend to set a value below (or at most at) 80% of your total limit.
  # The default value is 50 if you don't override it here.
  # rate_limit = 50
}
