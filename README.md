# Kong Plugin

# Run Locally

Use `docker-compose`, based on [their template](https://github.com/Kong/docker-kong/tree/master/compose):
```sh
docker-compose up -d
```

On the first time, you'll need to setup some test service and user:
```sh
./setup-kong.sh
```
