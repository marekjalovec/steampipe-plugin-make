#!/bin/bash
FILE=~/.steampipe/logs/plugin-$(date +"%Y-%m-%d").log

# clear both the screen and the log file
printf '\33c\e[3J'
printf "" > "$FILE"

# build the plugin and let me know when we are ready
if make; then
  printf "\e[32mREADY\e[0m\n"
  tail -f "$FILE"
else
  printf "\e[31mFAIL\e[0m\n"
fi
