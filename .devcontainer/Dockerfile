# See here for image contents: https://github.com/microsoft/vscode-dev-containers/tree/v0.238.0/containers/go/.devcontainer/base.Dockerfile

# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG VARIANT="1.18-bullseye"
FROM mcr.microsoft.com/vscode/devcontainers/go:0-${VARIANT}

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10
ARG NODE_VERSION="none"
RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends python3 python3-pip python3-dev \
    && pip install pre-commit

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode
RUN go install mvdan.cc/gofumpt@latest \
	&& go install honnef.co/go/tools/cmd/staticcheck@latest \
	&& go install github.com/kisielk/errcheck@latest \
    && go install golang.org/x/tools/cmd/goimports@latest \
    && go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

# Configure ZSH autocompletion
USER root

RUN chsh -s /usr/bin/zsh
# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
