FROM golang:1.22.4-bookworm
LABEL author="Konstantin Malikov"
LABEL description="Toolchain for project"

ENV \
    USER=k0st1am \
    TERM=xterm

RUN \
    adduser --disabled-password --gecos '' ${USER} && \
    chown -Rc ${USER}:${USER} "/home/${USER}/" && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
        apt-get --quiet --yes --no-install-recommends install \
            bash-completion chrpath curl dpkg dialog awscli \
            git locales make ssh sudo vim tig && \
    DEBIAN_FRONTEND=noninteractive \
        apt-get clean && \
    echo "${USER} ALL=NOPASSWD: ALL" > /etc/sudoers.d/${USER} &&\

WORKDIR /home/${USER}/project
VOLUME ["/home/${USER}/project"]
CMD ["/bin/bash"]
