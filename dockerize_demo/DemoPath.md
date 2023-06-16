
Cell 1:

This cell is importing all the necessary libraries that will be used in this Jupyter notebook.
os and json are part of Python's standard library, useful for interacting with the operating system and working with JSON data, respectively.
pandas is a data manipulation and analysis library, multiprocessing is used for creating parallel processes.
pydgraph is a Python client for interacting with Dgraph database, and GraphqlClient from python_graphql_client is used for sending GraphQL queries to a server.
dotenv is used to load environment variables from a .env file.

Cell 2:

This cell is loading environment variables which include username and password for the Graphistry database, as well as the Google Maps Key.
dotenv is a Python library used to read key-value pairs from a .env file and add them to the environment variables.

Cell 3:

This cell sets up the Dgraph and GraphQL clients.
The Dgraph client connects to the Dgraph service running locally on port 9080 and checks the version of Dgraph.
The GraphQL client is configured to connect to a GraphQL endpoint also running locally on port 8080.
The gql_admin_client is used to check the status of the cluster.

Cell 4:

This cell uses the GraphQL client to execute a query against the Dgraph database to get record counts for various data sources.
The %%time magic command at the top of the cell measures the execution time of the query.
The query checks the number of records from the Paradise Papers, Panama Papers, Bahamas Leaks, Offshore Leaks, Pandora Papers, and a total aggregate.
The results are then printed in a formatted manner.

Cell 5:

This cell defines a couple of functions to manipulate and extract information from a dictionary that is created from the result of a Dgraph query.
The update_node function updates a node if it doesn't exist and appends the value for the provided key.
The extract_dict function is recursive and used to traverse the query result dictionary and extract nodes and edges information.

Cell 6:

This cell contains a GraphQL query string to get information about records from the database.
The query is requesting several fields for each record, including its ID, type, name, and source, as well as relationships such as addresses, officers, intermediaries, and connections.

Cell 7:

This cell creates a function to execute the previous GraphQL query with specific offset and limit and another function to load all nodes and edges from the Dgraph database.
This operation is performed in parallel using multiple threads equal to half the number of available CPU cores, to improve the performance and speed of data retrieval.
The data retrieved includes the record count for each chunk of data and the time taken to load all nodes and edges.
The query function retrieves data from the database in chunks, using the provided offset and limit, while the load_all_nodes_and_edges function manages the concurrent execution of the queries, keeps track of the progress, and uses the extract_dict function to process the data.

Cell 8:

This cell initializes an empty dictionary for nodes and an empty list for edges.
Then it calls the function load_all_nodes_and_edges() from cell 7 with these empty data structures to fill them with nodes and edges data retrieved from the database.
The %%time magic command is used again to measure the execution time of the function.

Cell 9:

This cell imports the networkx library and uses it to create a directed graph, G, from the edges data.
It first converts the list of edges into a Pandas DataFrame, and then uses the from_pandas_edgelist() method of networkx to create the graph.
The nodes in the graph represent unique entities and the edges represent relationships between these entities.
It also prints out 3 random rows from the edges dataframe, the details of the created graph, the network density, and tries to print the diameter of the graph.
If the graph is not strongly connected, calculating the diameter would result in an error, which is handled by the exception block.

Cell 10:

This cell converts the dictionary of nodes into a pandas DataFrame and prints three random rows from it.
Each row in this dataframe represents a unique node or entity, and the columns are the attributes of these nodes.

Cell 11:

This cell finds the top 10 nodes in the graph by degree, which is the number of edges incident on a node.
The higher the degree of a node, the more connections it has.
It then prints the ID, name, type, and degree for each of these nodes.

Cell 12:

This cell calculates the degree centrality for each node in the graph and finds the top 10 nodes by this measure.
Degree centrality is a measure of the importance of a node in a network.
It is based on the concept that nodes with more connections are more central.
The nodes' ID, name, type, and centrality are printed.

Cell 13:

This cell calculates the PageRank of each node and identifies the top 10 nodes by this metric.
PageRank is a measure of the importance of nodes in a network, based on the concept that more important nodes are likely to receive more links from other nodes.
The ID, name, type, and PageRank of each of these nodes are printed.

Cell 14:

This cell calculates all shortest paths between pairs of nodes in the graph.
It then finds and stores all paths that have four or more edges.
The paths are stored in lp_list (which stands for 'long paths list').
The all_pairs_shortest_path function generates a dictionary of dictionaries where the key is a node and the value is a dictionary with keys as other nodes and values as the shortest path between the two nodes.
It then iterates through these paths and checks the length of each path, adding the ones with four or more edges to the list.
Finally, it sorts the list and prints it out.

Cell 15:

This cell is importing the graphistry package and using it to plot the graph we constructed in the previous steps.
It authenticates the Graphistry API using a username and password stored in environment variables and registers the user.
The Graphistry API provides a cloud-based service for graph visualization and analysis.
Then, it binds the nodes and edges to specific identifiers in the dataframe and encodes node colors based on their type.
Finally, it plots the graph.

Cell 16:

In this cell, we're querying the graph database to find the shortest path between two nodes.
The from_node and to_node variables represent the node IDs that mark the beginning and end of the path.
The Dgraph query language (DQL) is used to perform this operation, which queries a graph database using a syntax similar to GraphQL.
We then send the query to the database, load the result into the paths variable, and print the contents.

Cell 17:

This cell visualizes the shortest path we found in the previous cell using ipycytoscape, a widget library for interactive graphs.
It first converts the path into a format that ipycytoscape can understand and recursively finds the edge types.
Then it prints out the processed graph data.

Cell 18:

This cell sets the style of the nodes and edges in the graph visualization based on their types and properties, such as color, font, size, border, and label.
The pformat function from the pprint module is used to pretty-print the styles.

Cell 19:

This cell uses the ipycytoscape library to create a widget that displays the graph.
It adds the graph data to the widget, sets the layout and style of the graph, and finally displays the widget.

Cell 20:

Here, a text search is performed on the graph database using a GraphQL query.
The search is for any entities whose name contains the word "live".
The results are then printed.
The %%time magic command is used to measure the execution time of the cell.

Cell 21:

Here, a query is written to fetch addresses that have a location defined in the database.
The has function in the GraphQL filter is used to check for addresses that have a location.
After executing the query, the response is normalized into a pandas DataFrame.
Next, the function extract_names(l) is defined to extract names from a list and join them together into a single string, separated by commas.
This function is applied to the addressFor column to create a single string of names for each address.
Finally, the DataFrame is manipulated to rename some columns for easier interpretation, and a sample of five entries is displayed.

Cell 22:

In this cell, we utilize Bokeh's gmap function to create a map visualization of US Addresses.

Cell 23:

Here we're writing a query to find addresses near specific points, in this case, Syracuse, NY and Los Angeles.
We use the near function in the GraphQL filter to specify a location and distance, and get addresses that are within that distance.

Cell 24:

In this cell, we define a function is_flagged(node) to check if a node in the graph is flagged or not.
This function is then used in convert_to_cyto_objs(nodes, edges), which converts the nodes and edges into a format that can be read by Cytoscape for visualization.

Cell 25:

This cell runs a Dgraph query using the recurse directive to fetch a graph of nodes and edges starting from a specified node, up to a depth of 5.
The fetched data is then converted into a format suitable for Cytoscape and a graph is displayed.

Cell 26:

This cell executes a command to modify the GraphQL schema to add the flagged predicate to the Record type.

Cell 27:

Here we are making a mutation operation on the graph database, which is triggered by clicking a node in the Cytoscape graph.
It flags the clicked node and prints the mutation result.

Cell 28:

This cell reruns the recurse query to fetch the graph of nodes and edges, but this time also fetching the flagged predicate.
It displays the graph, with the nodes that have been flagged by the user displayed in a different style (defined by the cyto_styles variable).
