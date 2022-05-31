# cicada: Long Term Support Analyzer

![cicada molt](cicada.png)

> Shed your skin anew.

# SUMMARY

cicada scans your computer for outdated components.

# EXAMPLE

```console
$ cicada
warning: end of life for ruby v2.6.8 on 2022-03-31
```

See `cicada -help` for more detail.

# ABOUT

Many software components offer Long Term Support (LTS) releases, which receive security updates, bugfixes, and new features more rapidly than older releases. Unfortunately, it is often left up to the developer to opt into LTS releases. That is not an easy proposition, because software tends to grow in complexity over time. The dependency tree tends to get bigger and bigger. Meaning the risk of accidentally consuming a dead package is high. And the likelihood of spotting a dead package is low.

This is where cicada steps in. cicada helps engineers to identify more non-LTS software components. In an entirely automated fashion. cicada provides focused, actionable information for developers to implement. So that the larger software system remains robust, mature, and well supported by industry standards.

For enterprise systems, a variety of tools are available to identify specific vulnerabilities. But by then, it may be too late to migrate and keep up with cyberattacks. This is why cicada takes a more aggressive approach. We want to ward off entire major versions of software that do not receive ongoing patches. So that we can be one step ahead of attackers.

Following classical UNIX conventions, cicada often emits no output of any kind, and returns zero exit status, in the case where no non-LTS components are identified. To see a list of supported versions that cicada identifies on your machine, you may run `cicada -debug` to show additional logs. cicada is designed to integrate easily into full linter suites for CI/CD usage. Especially for Docker image builds. With a simple tweak to your base images, you can unlock more productive, more predictable SDLC workflows.

cicada is future compatible: Any software components targeting developmental, git HEAD versions should not trigger false alarms. Maybe you're already consuming clang tip, for example, which is well ahead of the pack. For the rest of us, we can relax and enjoy stable, LTS releases. Rest well on islands of stability. Let cicada guide you to successful operation.

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
