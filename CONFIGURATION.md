# CONFIGURATION

The cicada configuration is specified in the top level of your software project.

# OVERVIEW

cicada configuration uses YAML format.

Many of the keys, such as `debug`, `quiet`, and `lead_months`, are optional. Though the cicada linter may not behave as optimally without appropriate values.

For example, endoflife.date products without an associated entry in `version_queries` may be skipped during scanning, because there would be no known way to collect the live version information for comparison with support timelines.

# EXAMPLE

The Hello World demo project includes an example cicada configuration.

[example/cicada.yaml](example/cicada.yaml)

For more detail, see the [index.go](index.go) structure declaration that defines the configuration object model.

# TROUBLESHOOTING

When cicada is unable to query a semver compatible version string for an application, then it considers the application not installed, and skips over the application.

If you think an application is installed, but cicada is having trouble querying the application version, then enable debug mode to provide additional clues about how cicada collects bills of materials. For example, set `debug` to `true` in your `cicada.yaml` configuration, and/or supply a `-debug` flag to the `cicada`... command.
