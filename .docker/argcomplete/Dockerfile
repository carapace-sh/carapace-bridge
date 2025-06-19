FROM bitnami/azure-cli

USER root
RUN curl -so - https://dl.elv.sh/linux-amd64/elvish-v0.21.0.tar.gz | tar -xzvC /usr/local/bin

RUN mkdir -p ~/.config/elvish \
 && echo -e "set paths = [ /carapace-bridge/cmd/carapace-bridge \$@paths ]\neval (carapace-bridge _carapace|slurp)\neval (carapace-bridge argcomplete az|slurp)" > ~/.config/elvish/rc.elv
ENV PATH="/carapace-bridge/cmd/carapace-bridge:$PATH"

ENTRYPOINT /usr/local/bin/elvish
