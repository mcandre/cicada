# BUILDTIME REQUIREMENTS

* [ASDF](https://asdf-vm.com/) 0.10 (run `asdf reshim` after provisioning)
* [Go](https://go.dev/) 1.22.5+
* POSIX compatible [make](https://pubs.opengroup.org/onlinepubs/9699919799/utilities/make.html)
* [Node.js](https://nodejs.org/en) 20.10.0+
* [Ruby](https://www.ruby-lang.org/en/)
* [Rust](https://www.rust-lang.org/) 1.75.0+
* POSIX compatible [tar](https://pubs.opengroup.org/onlinepubs/7908799/xcu/tar.html)
* Provision additional dev tools with `make`

## Recommended

* [direnv](https://direnv.net/) 2
* a UNIX environment, such as macOS, Linux, BSD, [WSL](https://learn.microsoft.com/en-us/windows/wsl/), etc.

Non-UNIX environments may produce subtle adverse effects when linting or generating application ports.

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
