# LimeLog

### Env variables:

```
DEBUG: true
LOG_LEVEL: "debug"
HTTP_LISTEN: ":80"
AUTH_PASSWORD: "password"
SESSION_TOKEN: "token"
MONGO_HOST: host # default "localhost:27017"
MONGO_USERNAME: username
MONGO_PASSWORD: password
MONGO_DB_NAME: dbName
MONGO_REPLICA_SET: string # optional
NF_TELEGRAM_BOT_TOKEN: string # optional
NF_TELEGRAM_CHAT_ID: 123 # optional
NF_TELEGRAM_LEVELS: "fatal,error,warn" # optional
INPUT_GELF_ADDR: ":9234"
```

<br/>

### Install `swagger-cli`:

```
dir=$(mktemp -d) 
git clone https://github.com/go-swagger/go-swagger "$dir" 
cd "$dir"
go install ./cmd/swagger
rm -rf "$dir"
```
