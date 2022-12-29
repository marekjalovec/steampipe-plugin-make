connection "make" {
  plugin = "marekjalovec/make"

  # Make API token
  # To generate the token visit the API tab in your Profile page in Make.
  # Minimum scopes are connections:read and teams:read.
  # Add anything you want to have access to on top of that.
  api_key = "ecc9f531-xxxx-xxxx-xxxx-e6de6b3ebbe2"

  # Environment URL
  # The environment of Make you work in. This can be the link to your private instance of Make,
  # for example, https://development.make.cloud, or the link to Make
  # with or without the zone, depending on a specific endpoint, for example, https://eu1.make.com.
  environment_url = "https://eu1.make.com"
}
