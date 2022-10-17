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
