# vlg
Dgraph's Very Large Graph (WIP) Project

## Introduction
This work-in-progress project aims to create a public, large graph to showcase Dgraph's power and ease of use.

## Why
The current graph implementations used in benchmarking and "the tour" do not use new features and capabilities of Dgraph. And they exclude some of the difficult design problems that developers face when building production applications with Dgraph.

## What
Proposed is a publicly-available graph that represents Amazon's IMDB (Internet Movie Database). This dataset consists of millions of titles of movies and television series, along with cast and crew details. The dataset is updated daily. In addition, we'll scrape the Twitter API for title mentions to form new edges in the graph.

## How
We'll build the graph using practices common in production graphs:

* Use of GraphQL SDL for the schema
* Use of the bulk loader for initial loading
* Use of a (nightly) batch updater to load new titles and tweets into the graph
* Deploy the graph in Dgraph Cloud (TDB, maybe instead a k8s cluster to illustrate how that's best accomplished)
* We'll create a Docker container (separate repo) with a subset of the data that developers can use to quickly start with Dgraph

## Goals

1. Publicize the existence of the VLG to showcase Dgraph's performance
2. Inform developers of best practices in schema and query design

