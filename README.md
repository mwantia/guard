# coredns-guard

```
. {
  guard {
    url https://raw.githubusercontent.com/ph00lt0/blocklists/master/blocklist.txt adguard 3600
  }

  forward . 1.1.1.1 1.0.0.1
}
```