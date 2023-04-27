## Graph Analysis and Visualization

This notebook demonstrates querying, graph analysis and visualization, and geo coordinate searching.

Note this notebook requires a running cluster that has been populated with either the full or subset VLG data. If you are loading a subset (recommended on laptops), check out the [live loading section](https://github.com/dgraph-io/vlg/blob/main/notes/3.%20Data%20Importing.md#live-loading-a-subset) in the Notes to get your cluster up and running. Note that the [subset RDF files](https://github.com/dgraph-io/vlg/tree/main/rdf-subset) are checked into this repository, so no need to go through the RDF export phase.

### Optional Logins/API keys

This notebook uses [graphistry](https://graphistry.com) for graph visualization and Google Maps to render results from geo-coordinate searches. Both of these are optional, there's still plenty of useful things in the Notebook.

#### Graphistry
You can register for a free graphistry account [here](https://hub.graphistry.com/).

#### Google Maps
You can create a Google Maps API key [here](https://developers.google.com/maps/documentation/javascript/get-api-key).

### Setting Logins/API keys
The notebook uses a .env file to set the login and API key data. Check out the example.env file for the structure of the file.

Once your cluster is up and loaded, you can run the notebook

```
pipenv install
pipenv run juypter lab notebook.ipynb
```

Note, this notebook and associated package installs has been tested with python 3.8 and 3.11.
