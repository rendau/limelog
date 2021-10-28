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
