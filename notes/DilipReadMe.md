# create the network
docker network create dgraph_network

# start the dgraph zero container
docker run -d --network dgraph_network --name vlg_zero -p 5080:5080 dgraph/dgraph:latest dgraph zero --my=vlg_zero:5080

# start the dgraph bulk container
docker run --network dgraph_network -v $PWD:/home dgraph/dgraph:latest dgraph bulk -s /dev/null -g /home/schema/schema.graphql -f /home/rdf-subset --zero vlg_zero:5080 --out /home/out
