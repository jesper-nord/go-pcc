# go-pcc

## Introduction
`go-pcc` is a CLI to control Panasonic equipment connected through the Panasonic Comfort Cloud (PCC).

## CLI usage
You need to create a configuration file with the below content. The tool tries to load the default file `./go-pcc.yaml` but it can also be passed with the cli flag `-config [filepath]`.
```
username: [your PCC username]
password: [your PCC password]
```

### List devices
List all available Panasonic devices for account
```
$ go-pcc -list
```

### Examples
```
$ go-pcc -status
$ go-pcc -temp 19.5
$ go-pcc -on
$ go-pcc -off
$ go-pcc -mode heat
$ go-pcc -ecomode powerful
$ go-pcc -history week
```

For all available commands, see `go-pcc -help`.

### Logging
Enable debug logging with the `-debug` flag. Disable all logging with the `-suppress` flag.
