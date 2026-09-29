module github.com/carapace-sh/carapace-bridge/cmd

go 1.24.0

require (
	github.com/carapace-sh/carapace v1.17.0-alpha
	github.com/carapace-sh/carapace-bridge v0.0.0-00010101000000-000000000000
	github.com/carapace-sh/carapace-selfupdate v0.1.0-alpha
	github.com/spf13/cobra v1.10.2
)

require (
	github.com/carapace-sh/carapace-shlex/v2 v2.0.0-alpha // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.10 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
)

replace github.com/carapace-sh/carapace-bridge => ../
