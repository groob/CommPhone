Web scraper/CLI utility to print and search phone extensions for your ogranization.

# Example Usage:

```
$ ./commphone FirstName
First Last 182
$ ./commphone -line LastName
First Last 203-123-4567
$ ./commphone -all
First Last 182
First Last 181
...
First Last 100
```

# Configuration:
configured via environment variables

```bash
export COMMPORTAL_USER=2031234567
export COMMPORTAL_PASS=2222
# don't add a trailing slash to the end of the url
export COMMPORTAL_url=http://myloginurl
```

