module github.com/carapace-sh/carapace-bridge/cmd

go 1.24

require (
	github.com/carapace-sh/carapace v1.11.0
	github.com/carapace-sh/carapace-bridge v0.0.0-00010101000000-000000000000
	github.com/carapace-sh/carapace-selfupdate v0.0.5
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/carapace-sh/carapace-shlex v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.9 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/carapace-sh/carapace-bridge => ../
