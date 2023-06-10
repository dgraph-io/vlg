# Use Docker's official Docker in Docker (dind) image
FROM docker:20.10-dind

# Install basic utilities and build dependencies
RUN apk add --no-cache \
    curl \
    git \
    make \
    python3.9 \
    py3-pip \
    bash \
    build-base \
    libffi-dev \
    openssl-dev \
    python3-dev \
    cargo

# Install Docker Compose
RUN pip3 install docker-compose

# Install pipenv
RUN pip3 install pipenv

# Clone the vlg repo from a specific branch
RUN git clone -b dshekhar95/notebook-demo https://github.com/dgraph-io/vlg.git

# Set the working directory
WORKDIR vlg

# Set Docker to use the host's Docker daemon
ENV DOCKER_HOST=unix:///var/run/docker.sock

# Set entry point to bash so you can interact with the container
ENTRYPOINT ["/bin/bash"]
