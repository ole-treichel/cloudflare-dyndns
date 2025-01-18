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
    volumes:
        - ./config.json:/config.json
    ## if host network ist preferable
    # network_mode: "host"
    # environment:
    #    PORT: $PORT
    ports:
        - "$PORT:8000"
```

    docker compose up -d

The server exposes one endpoint that you can use as "update url":

`/dyndns?token=$TOKEN&ip=$IP`
`/dyndns?token=SHhilfqflVqtg8wsU7_kkdp4sv90_xOqzn_lut4q&ip=`

http://192.168.178.35/dyndns?token=SHhilfqflVqtg8wsU7_kkdp4sv90_xOqzn_lut4q&zone_id=f01b61b93e2310c1190ad6ff5e2439b8&domain=*.ole.md&domain=ole.md&ip=<ipaddr>

### Params

| Param   | Description                                                                                                               |
|---------|---------------------------------------------------------------------------------------------------------------------------|
| token   | Cloudflare api token. Needs `Zone:Edit` and `DNS:Edit` scopes.                                                            |
| ip      | The ip address to set on the dns record(s).                                                                               |

### Config file

Copy the example config file at `config.example.json` to `config.json` and update accordingly:

| Key     | Description                                                            |
|---------|------------------------------------------------------------------------|
| domains | List of domain records to update (e.g. `*.example.com`, `example.com`) |
| zoneId  | The zone id of the domain, found in the cloudflare dashboard           |

## Development

    go run cmd/cloudflare-dyndns/cloudflare-dyndns.go
