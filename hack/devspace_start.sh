#!/bin/bash
set +e 
unset GOPATH
export PATH="$(go env GOPATH)/bin:$PATH"

AIR_CMD="air -c air.toml"
BUILD_CMD="GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -gcflags 'all=-N -l' -o dist/cnfuzz-debug src/main.go"
DLV_CMD="dlv --listen=:2345 --headless=true --api-version=2 exec dist/cnfuzz-debug"

echo "# Working directory
# . or absolute path, please note that the directories following must be under root.
root = \".\"
tmp_dir = \"tmp\"

[build]
  cmd = \"${BUILD_CMD}\"
  bin = \"dist/cnfuzz-debug\"
  full_bin = \"while :; do ${DLV_CMD}; sleep 1; done\"
  include_ext = [\"go\"]
  exclude_dir = [\"dist\", \"docs\", \"charts\"]
  include_dir = [\"src\"]
  exclude_file = []
  exclude_regex = [\"_test.go\"]
  exclude_unchanged = true
  follow_symlink = true
  log = \"air.log\"
  delay = 1000 # ms
  stop_on_error = true
  send_interrupt = false
  kill_delay = 500 # ms
  # args_bin = []

[color]
  main = \"magenta\"
  watcher = \"cyan\"
  build = \"yellow\"
  runner = \"green\"

[misc]
  # Delete tmp directory on exit
  clean_on_exit = true
" > air.toml


COLOR_BLUE="\033[0;94m"
COLOR_GREEN="\033[0;92m"
COLOR_CYAN="\033[0;36m"
COLOR_RESET="\033[0m"

# Print useful output for user
echo -e "${COLOR_BLUE}
 .--. .-..-..---.
: .--': \`: :: .--'
: :   : .\` :: \`;.-..-..---. .---.
: :__ : :. :: : : :; :\`-'_.'\`-'_.'
\`.__.':_;:_;:_; \`.__.'\`.___;\`.___;
${COLOR_RESET}

Welcome to your development container!
- Files are synced from your host machine
- Port 2345 is forwarded

Some useful commands:
- Run \`${COLOR_CYAN}${BUILD_CMD}${COLOR_RESET}\` to build a binary usable for debugging
- Run \`${COLOR_CYAN}${DLV_CMD}${COLOR_RESET}\` to start debugging the binary from a remote debugger
- Run \`${COLOR_CYAN}${AIR_CMD}${COLOR_RESET}\` to live reload the debugger on changes made on your host machine
"
# echo ${COLOR_GREEN}TIP:${COLOR_RESET} check the bash history to find these commands (arrow up)
# history -s $DLV_CMD
# history -s $BUILD_CMD
# history -s $AIR_CMD

# Set terminal prompt
export PS1="\[${COLOR_BLUE}\]devspace\[${COLOR_RESET}\] ./\W \[${COLOR_BLUE}\]\\$\[${COLOR_RESET}\] "
if [ -z "$BASH" ]; then export PS1="$ "; fi

bash --norc
