# Ports

## Named Ports

```plaintext
In other Go projects or documentation you might sometimes see network addresses written using named ports like ":http" or ":http-alt" instead of a number. If you use a named port then Go will attempt to look up the relevant port number from your /etc/services file when starting the server, or will return an error if a match can’t be found.
```
