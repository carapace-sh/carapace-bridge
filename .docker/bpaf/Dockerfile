FROM alpine

RUN apk add --no-cache curl libc6-compat rust cargo

RUN curl -so - https://dl.elv.sh/linux-amd64/elvish-v0.21.0.tar.gz | tar -xzvC /usr/local/bin

COPY demo /build
RUN cd /build \
 && cargo build --release \
 && cp target/release/demo /usr/local/bin/demo \
 && rm -rf /build

RUN mkdir -p ~/.config/elvish \
 && echo -e "set paths = [ /carapace-bridge/cmd/carapace-bridge \$@paths ]\neval (carapace-bridge _carapace|slurp)\neval (carapace-bridge bpaf demo|slurp)" > ~/.config/elvish/rc.elv
ENV PATH="/carapace-bridge/cmd/carapace-bridge:$PATH"

ENTRYPOINT /usr/local/bin/elvish
