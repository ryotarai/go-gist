# go-gist
## Installation
```
$ go get github.com/ryotarai/go-gist
```

### Pre-built Binary
```
$ wget https://github.com/ryotarai/go-gist/releases/download/v[VERSION]/go-gist-[OS]-[ARCH]
```

* [VERSION]
* [OS]: linux, darwin
* [ARCH]: 386, amd64

## Usage
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

If you are using GitHub Enterprise, you should set `GITHUB_URL`:

```
$ export GITHUB_URL="https://your-ghe"
```


