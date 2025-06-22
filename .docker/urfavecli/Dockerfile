FROM golang

RUN curl -so - https://dl.elv.sh/linux-amd64/elvish-v0.21.0.tar.gz | tar -xzvC /usr/local/bin

ADD greet /greet

RUN cd /greet && go install

RUN mkdir -p ~/.config/elvish \
 && echo "set paths = [ /carapace-bridge/cmd/carapace-bridge \$@paths ]\neval (carapace-bridge _carapace|slurp)\neval (carapace-bridge urfavecli greet|slurp)" > ~/.config/elvish/rc.elv
ENV PATH="/carapace-bridge/cmd/carapace-bridge:$PATH"

ENTRYPOINT /usr/local/bin/elvish
