# go-gist
## Installation
```
$ go install github.com/ryotarai/go-gist
```

### Usage
```
$ go-gist -h
  -P=false: paste from the clipboard
  -c=false: copy gist URL to the clipboard
  -d="": description of the gist
  -h=false: show help
  -o=false: open the gist in a browser
  -p=false: private gist
  -v=false: show version
```

First, create a new token. (https://github.com/settings/tokens/new)

```
$ export GITHUB_TOKEN="..."
$ go-gist foo.txt
$ cat foo.txt | go-gist
```

