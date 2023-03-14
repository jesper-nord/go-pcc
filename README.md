# go-pcc

## Introduction
`go-pcc` is a CLI to control Panasonic equipment connected through the Panasonic Comfort Cloud (PCC).

## CLI usage
You need to create a configuration file with the below content. The tool tries to load the default file `./gopcc.yaml` but it can also be passed with the cli flag `-config [filepath]`.
```
username: [your PCC username]
password: [your PCC password]
```

### List devices
List all available Panasonic devices for account
```
$ gopcc -list
```

### Examples
```
$ gopcc -status
$ gopcc -temp 19.5
$ gopcc -on
$ gopcc -off
$ gopcc -mode heat
$ gopcc -ecomode powerful
$ gopcc -history week
```

For all available commands, see `gopcc -help`.

### Logging
Enable debug logging with the `-debug` flag. Disable all logging with the `-suppress` flag.
