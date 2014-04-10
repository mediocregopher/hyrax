# Hyrax

A scalable backend for real-time apps. Provides storage and events about said
storage for clients, as well as authentication.

Hyrax is built as a layer between the backend application logic and the
frontend, allowing both to store state and communicate about said state. With
hyrax you can:

* Retrieve/modify keys
* Have clients receive real-time updates when a key changes
* Have clients receive real-time updates when another client disconnects
* Scale to many nodes, each holding many concurrent client connections
* Control what commands clients are allowed to call, and in what context

Hyrax is still under active development, and things are likely to change in the
near future. But please feel free to poke around the code and play with it, and
let me know (at the email in my profile) if there's anything I can do to make
things clearer.

## Table Of Contents

**Getting started**

* [Overview](/doc/overview.md)
* [Installation/Configuration](/doc/installconfig.md)

**Using hyrax**

* [Basics](/doc/basics.md)
* [Protocols/Syntaxes](/doc/protosyntax.md)
* [Authentication](/doc/auth.md)
* [Clients](/doc/clients.md)

**Commands**

* [Mon](/doc/mon.md) - monitor changes to keys
* [Ekg](/doc/ekg.md) - monitor other clients
* [Admin](/doc/admin.md) - Commands for administering a single hyrax node

See the [redis][redis] page for other available commands

**Backends** (currently only one)

* [Redis][redis]

**Deployment**

* Topology examples
* Scaling

## Contact

My name is Brian Picciano. You can get ahold of me at the email in my github
profile (github.com/mediocregopher).

[redis]: /doc/redis.md
