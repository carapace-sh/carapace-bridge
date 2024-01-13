FROM ghcr.io/rsteube/carapace:latest

RUN apt-get update \
 &&  apt-get install -y bash-completion \
                        npm \
                        pip

# argcomplete
RUN pip install --break-system-packages azure-cli

# click
RUN pip install --break-system-packages td-watson

# cobra
RUN curl -Lo /usr/local/bin/minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64 \
 && chmod +x /usr/local/bin/minikube

# complete
RUN curl --proto '=https' --tlsv1.2 -fsSL https://get.opentofu.org/install-opentofu.sh | sh -s -- --install-method deb

# inshellisense
RUN npm install -g @microsoft/inshellisense
RUN cd /usr/local/lib/node_modules/@microsoft/inshellisense/ \
 && npm i @withfig/autocomplete@2.648.2

# kingpin
RUN curl https://goteleport.com/static/install.sh | bash -s 14.3.3

# urvavecli
RUN curl -Lo /usr/local/bin/tea https://dl.gitea.com/tea/0.9.2/tea-0.9.2-linux-amd64 \
 && chmod +x /usr/local/bin/tea

# yargs
RUN npm install -g @angular/cli
