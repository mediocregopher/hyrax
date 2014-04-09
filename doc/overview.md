# Overview

## Introduction

There are two components to hyrax: the hyrax server(s) which clients connect to
and interact with, and the storage backend where state is stored (currently
redis, although others can/will be supported).

![Full hyrax stack](/doc/img/fullstack.png)

Clients send all commands through the hyrax server, which proxies them off to
the storage backend when appropriate. Hyrax also provides facilities so that
clients can monitor keys within the datastore and have notifications pushed to
them when any key is modified.

## Use case

Hyrax's best use case is one in which multiple clients need to communicate with
each-other about a shared state. For example: a simple service where you can set
your current status and receive real-time updates when your friend's statuses
change. You would have a key per-user which gets set and then monitor the keys
of all of your friends for changes.

## Scaling

Hyrax scales up to an arbitrary number of nodes. A lot of services claim this,
hyrax's does it. Nodes are organized in a tree-like structure such that branches
never even need interact with each other. Node failure/maintence is clean
because very few nodes are actually affected by the failure.

![Example of the tree topology](/doc/img/tree-example.png)

A very conscious design choice made by hyrax is to remain agnostic to the
storage backend. This allows for flexibility in application design, but more
importantly allows for flexibility in scaling. There are fundemental differences
in how you scale a message system and a storage system, and with hyrax the two
concerns remain completely separate.
