# Cloudflare Dynamic DNS

## Usage

The server exposes one endpoint that you can use as 'update url'.

`/dyndns?token=$TOKEN&zone_id=$ZONE_ID&domain=$DOMAIN&domain=$DOMAIN&ip=$IP`

### Params

| Param   | Description                                                                                                               |
|---------|---------------------------------------------------------------------------------------------------------------------------|
| token   |  Cloudflare api token. Needs `Zone:Edit` and `DNS:Edit` scopes                                                            |
| zone_id | Cloudflare [zone id](https://developers.cloudflare.com/fundamentals/setup/find-account-and-zone-ids/).                    |
| domain  | The domain name you wish to update (e.g. example.com). Can be repeated to update multiple domain at once.                 |
| ip      | The ip address to set on the dns record(s).                                                                               |
