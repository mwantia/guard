# CoreDNS Guard Plugin

## Configuration Example

```
. {
  guard {
    ttl 120
    url https://raw.githubusercontent.com/ph00lt0/blocklists/master/blocklist.txt adguard 3600
  }

  forward . 1.1.1.1 1.0.0.1
}
```

## Default Configuration

### Default Refresh Frequency

Config: `frequency` \
Value:  `0s`

### Default TTL Answer

Config: `ttl` \
Value:  `14400`

### Default IPv4 Answer

Config: `ipv4` \
Value:  `0.0.0.0`

### Default IPv6 Answer

Config: `ipv6` \
Value:  `::`

### Next Or Failure

Config: `next` \
Value:  `true`

## TO-DO

It seems that adguard supports multiple entries per line that are separated by comma ','

```
4-liga.com,4fansites.de,4players.de,9monate.de,aachener-nachrichten.de,aachener-zeitung.de,abendblatt.de,abendzeitung-muenchen.de,about-drinks.com,abseits-ka.de,airliners.de,ajaxshowtime.com,allgemeine-zeitung.de,antenne.de,arcor.de,areadvd.de,areamobile.de,ariva.de,astronews.com,aussenwirtschaftslupe.de,auszeit.bio,auto-motor-und-sport.de,auto-service.de,autobild.de,autoextrem.de,autopixx.de,autorevue.at,az-online.de,baby-vornamen.de,babyclub.de
```