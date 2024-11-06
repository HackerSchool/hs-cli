# hs-cli
# Usage
```
NAME:
   hscli - CLI client for the HackerSchool API

USAGE:
   hscli [global options] command [command options]

VERSION:
   0.0.1

COMMANDS:
   mgetall      retrieve all members
   mget         retrieve information of a member
   mcreate      create member
   mupdate      update member information
   mdelete      delete member from the database
   mprojects    get the projects a member is in
   mlogo        get member logo
   mtags        get member tags
   maddproject  add a project to a member
   maddlogo     upload member logo
   maddtag      add member tag
   mdeltag      delete member tag
   pgetall      retrieve all projects
   pget         retrieve information of a project
   pcreate      create a new project
   pupdate      update information of a project
   pdelete      delete project from the database
   pmembers     get members in a project
   plogo        get project logo
   paddmember   add member to a project
   login        login to the API, saving the cookie to the cookiejar
   logout       logoout off the API, clearing the session
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -f value      path to config file
   --root value, -r value        API root url        (overwrites file and HS_ROOT environment configs)
   --user value, -u value        username            (overwrites file and HS_USER environment configs)
   --password value, -p value    user password       (overwrites file and HS_PASSWORD environment configs)
   --cookie-jar value, -c value  cookie jar path     (overwrites file and HS_COOKIEJAR environment configs)
   --debug, -d                   log debug information to the console (default: false)
   --help, -h                    show help
   --version, -v                 print the version
```

# Configuration
If configuration values are not provided as CLI arguments the program will first attempt to load them from the `--config` option and, after that, look for the environment variables and overwrite any options set by the configuration file. This way the order of preference is CLI args > environment > configuration file.

Example `config.yaml` file:
```yml
root: https://api.hackerschool.dev
user: username 
password: password
cookiejar: ./cookiejar.json
```
Example `.env` file:
```sh
export HS_ROOT="https://api.hackerschool.dev"
export HS_USER="username"
export HS_PASSWORD="password"
export HS_HS_COOKIERJAR="/home/my/cookiejar.json"
```

# Examples
## Configuration
Example using CLI options:
```sh
hscli -r https://api.hackerschool.dev -u username -p password -c cookiejar.json -d login 
```
Example using `config.yaml` file:
```sh
hscli --config config.yaml -d login
```

Example using `.env` file:
```sh
source .env; hscli -d login
```
The `source` command only needs to be ran once per shell session.

## Command Arguments 
For commands that expect a payload, such as `mcreate`, the `[<file>]` argument is optional, if ommited, the program will attempt to read the payload from standard input. This allows for some flexibility, e.g, the two following examples accomplish the same:
```sh
hscli mcreate member.json 
```
```sh
cat member.json | hscli mcreate
```

## Exit Codes 
The program returns `1` for API errors and `2` for other errors, e.g, "no connection to host", etc.
This can be leveraged for scripting.
```bash
hscli -d mget username > /dev/null
if [ $? -eq 0 ]; then
    hscli -d mgetlogo username > logos/username.png
    echo "Logo saved sufessfully!"
elif [ $? -eq 1 ]; then
    echo "API Error!"
else
    echo "Error!"
fi
```
This sript will save the logo of a user if it exists.

## Output 
The program writes raw JSON to `stdout` and error and log messages to `stderr`. Because of this it's recommended to make use of other programs such as `jq`.
```sh
hscli mgetall | jq
```
```sh
hscli maddtag dev_tag.json | jq
```
```sh
hscli pmembers proj_name | jq
```
```sh
cat updated_project.json | hscli update proj_name | jq
```