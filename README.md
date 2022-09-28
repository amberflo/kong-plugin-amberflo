# Kong Plugin

# How to Use

1. Build the plugin

```sh
make
```

2. Start Kong. We're using `docker-compose`, based on [their template](https://github.com/Kong/docker-kong/tree/master/compose):
```sh
docker-compose up -d
```

3. On the first time, setup a test service, route and user (see [this tutorial](https://docs.konghq.com/gateway/3.0.x/get-started/key-authentication/))
```sh
./setup-kong.sh
```

4. Restart the kong server whenever you rebuild the plugin

5. Use the helper scripts to interact with Kong:
```
./admin.sh POST /plugins -d '{"name": "metering"}'
./test.sh GET /mock/requests -H 'x-api-key: super-secret-key' -i
```
