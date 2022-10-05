# Amberflo Kong Metering Plugin

<p>
    <a href="https://github.com/amberflo/kong-plugin/actions">
        <img alt="CI Status" src="https://github.com/amberflo/kong-plugin/actions/workflows/test.yml/badge.svg?branch=main">
    </a>
</p>

[Amberflo](https://amberflo.io) is the simplest way to integrate metering into your application. [Sign up for free](https://ui.amberflo.io/) to get started.

This is the official Amberflo plugin for Kong. Use it to meter the requests handled by your Kong instance and monetize your APIs. Check out [our docs](https://docs.amberflo.io/docs) to learn more.

```sh
go install github.com/amberflo/kong-plugin-amberflo@latest
```

## :zap: How it Works

This plugin will intercept the requests, detect which customer is making the request, generate a meter event and send it to Amberflo.

Customer detection happens via inspection of the request headers. You can configure Kong to inject the customer id as a header before this plugin runs. For instance, if you use Kong's [Key Authentication](https://docs.konghq.com/hub/kong-inc/key-auth/) plugin, this happens automatically.

To avoid impacting the performance of your gateway, the plugin will batch meter records and send them asynchronously to Amberflo.

## :rocket: Installation

1. Compile and make the binary available to your Kong instance.
    - Make sure your compilation environment is compatible with your Kong environment, otherwise the compile binary won't work.

2. Update your Kong configuration with the now-available metering plugin.
    - For instance. Suppose you place the plugin server binary at `/opt/amberflo/metering` on the Kong server. Then the plugin-related parts of your `kong.conf` file should look like [this one](./kong.conf).
    - For more details on how to configure Kong, check out their docs [here](https://docs.konghq.com/gateway/latest/plugin-development/pluginserver/go/#example-configuration) and [here](https://docs.konghq.com/gateway/latest/reference/configuration/)

3. Enable the plugin
    - Either by adding it to your `kong.yaml` file or making an Admin API request.

## :scroll: Configuration

Please find a sample configuration file [here](./metering.json).

Here's a breakdown of the fields and their respective meanings.

| Name             | Type              | Required? | Default      | Description                                                                 |
|------------------|-------------------|-----------|--------------|-----------------------------------------------------------------------------|
| apiKey           | string            | yes       |              | Your Amberflo API key                                                       |
| meterApiName     | string            | yes       |              | Meter for metering the requests                                             |
| customerHeader   | string            | yes       |              | Header from which to get the Amberflo `customerId`                          |
| intervalSeconds  | int               | no        | `1`          | Send the meter record batch every `x` seconds                               |
| batchSize        | int               | no        | `10`         | Send the meter record batch when it reaches this size                       |
| debug            | bool              | no        | `false`      | Enable debug mode of the Amberflo API client (for development)              |
| methodDimension  | string            | no        |              | Dimension name for the request method                                       |
| hostDimension    | string            | no        |              | Dimension name for the target url host                                      |
| routeDimension   | string            | no        |              | Dimension name for the route name                                           |
| serviceDimension | string            | no        |              | Dimension name for the service name                                         |
| dimensionHeaders | map[string]string | no        |              | Map of "dimension name" to "header name", for inclusion in the meter record |
| replacements     | map[string]string | no        | `{"/": ":"}` | Map of "old" to "new" values for transforming dimension values              |

## :construction_worker: Developers

1. Build the plugin

```sh
make
```

2. Start Kong. We're using `docker-compose`, based on [their template](https://github.com/Kong/docker-kong/tree/master/compose):
```sh
make kong-up
```

3. On the first time, setup a test service, route and user (see [this tutorial](https://docs.konghq.com/gateway/3.0.x/get-started/key-authentication/))
```sh
make kong-init
```

4. Restart the kong server whenever you rebuild the plugin
```sh
make update
```

5. Use the helper scripts to interact with Kong:
```sh
./scripts/admin.sh GET /plugins
./scripts/test.sh GET /mock/requests -H 'x-api-key: super-secret-key' -i
```

6. To update the plugin configuration:
```sh
./scripts/admin.sh PUT /plugins/<plugin-id> -d @metering.json
```

7. To stop and clean-up docker resources
```sh
make kong-down
```
