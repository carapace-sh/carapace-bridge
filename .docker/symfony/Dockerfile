FROM php:8.3-cli-alpine

RUN apk add --no-cache curl libc6-compat

RUN curl -so - https://dl.elv.sh/linux-amd64/elvish-v0.21.0.tar.gz | tar -xzvC /usr/local/bin

COPY --from=composer:2 /usr/bin/composer /usr/local/bin/composer
COPY demo /usr/local/bin/demo
RUN chmod +x /usr/local/bin/demo \
 && mkdir -p /opt/demo \
 && echo '{}' > /opt/demo/composer.json \
 && composer --working-dir=/opt/demo require symfony/console --no-interaction

RUN mkdir -p ~/.config/elvish \
 && echo -e "set paths = [ /carapace-bridge/cmd/carapace-bridge \$@paths ]\neval (carapace-bridge _carapace|slurp)\neval (carapace-bridge symfony demo|slurp)" > ~/.config/elvish/rc.elv
ENV PATH="/carapace-bridge/cmd/carapace-bridge:$PATH"

ENTRYPOINT /usr/local/bin/elvish
