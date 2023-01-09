# vlg
Dgraph's Very Large Graph (WIP) Project

## Introduction
This work-in-progress project aims to create a public, large graph to showcase Dgraph's power and ease of use.

## Why
Unlike our competitors, Dgraph does not maintain a large, publicly available graph. Further, the current graph implementations used in benchmarking and "the tour" do not use new features and capabilities of Dgraph. And they exclude some of the difficult design problems that developers face when building production applications with Dgraph.

## What
Proposed is a publicly-available graph that presents the [Offshore Data Leaks](https://offshoreleaks.icij.org/) dataset from the International Consortium of Investigative Journalists (ICIJ), which compromise millions of entities encompassing the Panama Papers, Pandora Papers, and others.

## How
We'll build the graph using practices common in production graphs:

* Analysis of the data, schema design, import and update strategies
* Create the GraphQL SDL for the schema and document the reasons for decisions made
* Use of the bulk loader for initial loading
* Use of a batch updater (maybe several languages?) to load new data into the graph
* Deploy the graph to HA production (Dgraph Cloud or maybe instead a k8s cluster to illustrate how that's best accomplished)
* Creation of a Docker container (separate repo) with a subset of the data that developers can use to quickly start with Dgraph (one that's more complex than the 'million-movie' dataset).
* Creation of a simple UI that allows for search and visualization of the graph

This repo will contain folders for each step, and will include detailed documentation descibing the 'whys' and 'hows' behind the step. This README will be replaced my an _Introduction_ page.

## Goals

1. Publicize the existence of the VLG to **showcase Dgraph's performance**
2. Use the repo to **instruct developers** on best practices in schema, ETL and query design, as well as best deployment practices
3. Serve as a test fixture for upcoming releases of Dgraph

## Future work

1. Extend the graph to allow user accounts. Extend the sample UI to allow accounts to save searches and make personal notes associated with graph entities.
