docker compose down
docker compose stop
sudo rm -rf ~/dgraph/vlg
sudo rm -rf out
docker compose up -d
sleep 20
curl -Ss --data-binary '@./schema/schema.graphql' http://localhost:8080/admin/schema
sleep 20
docker run -d --network vlg_default -v $PWD:/home dgraph/dgraph:latest dgraph live -f /home/rdf-subset --alpha vlg_alpha:9080 --zero vlg_zero:5080
sleep 30
cd notebook/graph-analysis-and-visualization

pipenv install

pipenv run jupyter lab notebook.ipynb

