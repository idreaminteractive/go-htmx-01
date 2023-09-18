FROM gitpod/workspace-full:2023-08-30-14-07-38

RUN sudo apt-get update && sudo apt-get install -y fuse3 sqlite3 docker-buildx-plugin

RUN curl -L https://fly.io/install.sh | sh
ENV FLYCTL_INSTALL="/home/gitpod/.fly"
ENV PATH="$FLYCTL_INSTALL/bin:$PATH"

# install doppler locally.
RUN (curl -Ls --tlsv1.2 --proto "=https" --retry 3 https://cli.doppler.com/install.sh || wget -t 3 -qO- https://cli.doppler.com/install.sh) | sudo sh

RUN sudo mkdir /litefs
RUN sudo chown -R gitpod /litefs
# install air
RUN go install github.com/cosmtrek/air@latest

# install static analysis (not sure if we need this?)
RUN go install honnef.co/go/tools/cmd/staticcheck@latest

# install task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b ~/.local/bin

# install goose
RUN  go install github.com/pressly/goose/v3/cmd/goose@latest

RUN echo 'alias home="cd ${GITPOD_REPO_ROOT}"' | tee -a ~/.bashrc ~/.zshrc

