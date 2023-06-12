PWD := $(shell pwd)

SCHEMA_PATH := ./schema/schema.graphql
RDF_PATH := rdf-subset
NOTEBOOK_PATH := notebook/graph-analysis-and-visualization
DOCKER_COMPOSE_FILE := docker-compose.yml
DOCKER_NETWORK := $(shell basename $(PWD))_default


.PHONY: clean init load_data install_dependencies start_notebook

all: clean init load_data install_dependencies start_notebook

install_docker_pipenv:
	if [ "$$(uname)" = "Darwin" ]; then \
		brew install docker docker-compose pipenv; \
		pipenv --python 3.9 install; \
	elif [ "$$(uname)" = "Linux" ]; then \
		sudo apt-get update -y; \
		sudo apt-get install -y docker.io docker-compose python3-pip python3-dev; \
		sudo pip3 install pipenv; \
	fi

clean:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) down
	@docker-compose -f $(DOCKER_COMPOSE_FILE) stop
	@sudo rm -rf ~/dgraph/vlg
	@sudo rm -rf out

init:
	@docker-compose -f $(DOCKER_COMPOSE_FILE) up -d
	@bash -c "sleep 20"
	@curl -Ss --data-binary '@$(SCHEMA_PATH)' alpha:8080/admin/schema


load_data:
	@docker run --network $(DOCKER_NETWORK) -v $(PWD):/home dgraph/dgraph:latest dgraph live -f /home/$(RDF_PATH) --alpha vlg_alpha:9080 --zero vlg_zero:5080


setup_virtual_env:
	@cd $(NOTEBOOK_PATH) && pipenv install

start_notebook:
	@cd $(NOTEBOOK_PATH) && pipenv run jupyter lab notebook.ipynb
