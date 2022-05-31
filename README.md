# cicada: Long Term Support Analyzer

![cicada molt](cicada.png)

> Shed your skin anew.

# EXAMPLE

```console
$ cicada
warning: end of life for ruby v2.6.8 on 2022-03-31
```

See `cicada -help` for more detail.

# ABOUT

Many software components offer Long Term Support (LTS) releases, which receive security updates, bugfixes, and new features more rapidly than older releases. This is where cicada steps in. cicada helps engineers to identify more non-LTS software components. cicada provides focused, actionable information for developers to implement. So that the larger software system remains robust, mature, and well supported by industry standards.

Following classical UNIX conventions, cicada often emits no output of any kind, and returns zero exit status, in the case where no non-LTS components are identified. To see a list of supported versions that cicada identifies on your machine, you may run `cicada -debug` to show additional logs.

# DOCUMENTATION

https://godoc.org/github.com/mcandre/cicada

# DOWNLOAD

https://github.com/mcandre/cicada/releases

# INSTALL FROM SOURCE

```console
$ go install github.com/mcandre/cicada/cmd/cicada@latest
```

# LICENSE

FreeBSD

# RUNTIME REQUIREMENTS

(None)

# CONTRIBUTING

For more information about developing cicada itself, see [DEVELOPMENT.md](DEVELOPMENT.md).

# CREDITS

* [endoflife.date](https://endoflife.date/) for support lifetime data
