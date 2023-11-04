# Cloudflare Dynamic DNS

## Usage

##### Docker
    docker run -d -p $PORT:8000 ghcr.io/ole-treichel/cloudflare-dyndns:latest

##### Docker Compose

```
version: '3.3'

services:
  cloudflare_dyndns:
    restart: always
    image: "ghcr.io/ole-treichel/cloudflare-dyndns:latest"
    ## if host network ist preferable
    # network_mode: "host"
    # environment:
    #    PORT: $PORT
    ports:
        - "$PORT:8000"
```

    docker compose up -d

The server exposes one endpoint that you can use as "update url":

`/dyndns?token=$TOKEN&zone_id=$ZONE_ID&domain=$DOMAIN&domain=$DOMAIN&ip=$IP`

### Params

| Param   | Description                                                                                                               |
|---------|---------------------------------------------------------------------------------------------------------------------------|
| token   | Cloudflare api token. Needs `Zone:Edit` and `DNS:Edit` scopes.                                                            |
| zone_id | Cloudflare [zone id](https://developers.cloudflare.com/fundamentals/setup/find-account-and-zone-ids/).                    |
| domain  | The domain name you wish to update (e.g. `example.com`). Can be repeated to update multiple domains at once.              |
| ip      | The ip address to set on the dns record(s).                                                                               |

## Development

    go run cmd/cloudflare-dyndns/cloudflare-dyndns.go

