FROM alpine

RUN apk add --no-cache bash curl libc6-compat

RUN curl -so - https://dl.elv.sh/linux-amd64/elvish-v0.21.0.tar.gz | tar -xzvC /usr/local/bin

RUN curl https://cdn.teleport.dev/install.sh | bash -s 17.5.2 oss

RUN mkdir -p ~/.config/elvish \
 && echo -e "set paths = [ /carapace-bridge/cmd/carapace-bridge \$@paths ]\neval (carapace-bridge _carapace|slurp)\neval (carapace-bridge kingpin tsh|slurp)" > ~/.config/elvish/rc.elv
ENV PATH="/carapace-bridge/cmd/carapace-bridge:$PATH"

ENTRYPOINT /usr/local/bin/elvish
