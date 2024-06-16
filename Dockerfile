FROM golang:1.21
LABEL author="Konstantin Malikov"
LABEL description="Toolchain for project"

# Define ARG after FROM to indicate values coming from build arguments are part of the build stage.
# From https://stackoverflow.com/questions/31198835/can-we-pass-env-variables-through-cmd-line-while-building-a-docker-image-through
ARG DOCKER_USER

# Create the environment variables and assign the values from the build arguments.
ENV \
    USER=${DOCKER_USER}

RUN \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive \
        apt-get --quiet --yes --no-install-recommends install \
            bash-completion chrpath curl dpkg dialog awscli \
            git locales make ssh sudo vim tig adduser mc net-tools && \
    adduser --disabled-password --gecos '' ${USER} && \
    chown -Rc ${USER}:${USER} "/home/${USER}/" && \
    DEBIAN_FRONTEND=noninteractive \
        apt-get clean && \
    rm -rf /var/cache/* /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    echo "${USER} ALL=NOPASSWD: ALL" > /etc/sudoers.d/${USER} && \
    git config --global --add safe.directory /home/${USER}/project

WORKDIR /home/${USER}/project
# Use ${USER} user inside docker - another method to set user
USER ${USER}
CMD ["/bin/bash"]
