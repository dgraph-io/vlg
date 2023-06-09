sudo rm -rf ~/dgraph/vlg
sudo rm -rf out/
docker stop vlg-zero-1
docker rm vlg-zero-1
docker stop vlg_zero
docker rm vlg_zero
docker stop vlg_alpha
docker rm vlg_alpha
# create the network
docker network create dgraph_network
echo "XX1"
# start the dgraph zero container
docker run -d --network dgraph_network --name vlg_zero -p 5080:5080 dgraph/dgraph:latest dgraph zero --my=vlg_zero:5080
sleep 5
echo "XX2"
docker ps
# start the dgraph bulk container
docker run -d --network dgraph_network -v $PWD:/home dgraph/dgraph:latest dgraph bulk -s /dev/null -g /home/schema/schema.graphql -f /home/rdf-subset --zero vlg_zero:5080 --out /home/out
sleep 20
echo "XX3"
docker ps
sleep 5
sudo cp -r out/0/p ~/dgraph/vlg/
sleep 5
# docker stop vlg_zero
# docker rm vlg_zero
docker ps
sleep 10
docker compose up alpha -d
#curl -Ss --data-binary '@./schema/schema.graphql' http://localhost:8080/admin/schema
sleep 5
cd notebook/graph-analysis-and-visualization

pipenv install

pipenv run jupyter lab notebook.ipynb
