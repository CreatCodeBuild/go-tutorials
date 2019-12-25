FROM gitpod/workspace-full

USER root
RUN apt-get update && \
    apt-get install -y zsh && \
    chsh -s $(which zsh)
