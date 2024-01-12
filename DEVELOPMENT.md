# BUILDTIME REQUIREMENTS

* [ASDF](https://asdf-vm.com/) 0.10 (run `asdf reshim` after provisioning)
* [Go](https://golang.org/) 1.21.5+
* [Node.js](https://nodejs.org/en) 20.10.0+
* [Ruby](https://www.ruby-lang.org/en/)
* [Rust](https://www.rust-lang.org/) 1.75.0+
* a POSIX compliant [make](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html) implementation (e.g. GNU make, BSD make, etc.)
* Provision additional dev tools with `make`

## Recommended

* [direnv](https://direnv.net/) 2

# SECURITY AUDIT

```console
$ mage audit
```

# INSTALL

```console
$ mage install
```

# UNINSTALL

```console
$ mage uninstall
```

# LINT

```console
$ mage lint
```

# TEST: UNIT + INTEGRATION

```console
$ mage test
```

# PORT

```console
$ mage port
```
