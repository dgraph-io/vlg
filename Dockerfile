# Use Docker's official Debian-based Python image at version 3.9
FROM python:3.9-buster

# Install Docker CLI
RUN apt-get update && apt-get install -y apt-transport-https ca-certificates curl gnupg lsb-release
RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null
RUN apt-get update && apt-get install -y docker-ce-cli

# Install Docker Compose using pip
RUN pip3 install docker-compose

# Verify that Docker Compose is installed correctly
RUN docker-compose --version

# Install pipenv
RUN pip3 install pipenv

# Install Git
RUN apt-get install -y git

# Clone the vlg repo from a specific branch
RUN git clone -b dshekhar95/notebook-demo https://github.com/dgraph-io/vlg.git

# Set the working directory
WORKDIR vlg

# Set Docker to use the host's Docker daemon
ENV DOCKER_HOST=unix:///var/run/docker.sock

# Set entry point to bash so you can interact with the container
ENTRYPOINT ["/bin/bash"]
