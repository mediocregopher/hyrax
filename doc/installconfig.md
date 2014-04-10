# Installation

Hyrax consists of two distinct parts. The actual hyrax node, and the chosen
storage backend. The installation of your storage backend is not covered as it
depends on which you choose (although currently there is only one).

Hyrax itself is a simple static binary, which can either be compiled or
downloaded.

## Download

Check the [releases][releases] tab for downloadable binaries. Each binary in a
release corresponds to different systems and architectures.

## Compile

Compilation will require having [goat][goat] on your `$PATH`. You can follow
these steps to compile the server:

```bash
> git clone https://github.com/mediocregopher/hyrax.git
> cd hyrax
> goat deps
> goat build server/hyrax-server.go
```

At this point you should have a `hyrax-server` file in the root of the project.

# Configuration

Hyrax can take in either command-line options or a configuration file, or both
(command-line options supercede the configuration file).

## Command line

You can see all available parameters by calling:

```
hyrax-server --help
```

## Configuration file

You can load in a config file by calling:

```
hyrax-server --config <configfile>
```

and you can get an example configuration (pre-filled with default values and
options) by calling:

```
hyrax-server --example
```

(append ` > hyrax.conf` to save the example configuration to a file)

## Parameters

*Note: Many parameters take the form of a "listen endpoint". An endpoint is
described by the form `<protocol>::<syntax>::<listen address>`, where `protocol`
and `syntax` are one of those mentioned in
[Protocols/Syntaxes](/doc/protosyntax.md).*

```
* Can be updated without restarting using a USR1 signal
```

* `storage-info` - The actual form this takes will depend on the storage backend
  used. Consult the doc page for the backend you're using for the exact format
  (currently there is only [redis](/doc/redis.md)).

* `listen-endpoint` - Specifies that this hyrax node should listen for clients
  at this listen endpoint, using the endpoint's specified protocol and format.
  Can be specified 0 or more times.

* `push-to-endpoint` * - A listen endpoint to push key change events that happen
  on this actual node to. This should be the "root" node in the topology (see
  [Topology Examples][topology] for more details). For a single node setup, this
  would loopback to the current node. Can be specified 0 or more times.

* `pull-from-endpoint` * - A listen endpoint to pull key change events that
  happen on other nodes from. See [Topology Examples][topology] for more
  details. Can be specified 0 or more times.

* `interaction-secret` * - If [Authentication][auth] is enabled on any other
  hyrax nodes, this must be set to one of their global secret keys so this node
  can interact with them.

* `my-endpoint` - If this node is interacting with other nodes in any way, this
  is the listen endpoint those nodes should connect to to interact with this
  node.

* `log-level` * - The minimum log level to output. Can be `debug`, `info`,
  `warn`, `error`, or `fatal`.

* `log-file` * - The place to log all files to. Can also be set to `stdout` or
  `stderr`. Note that hyrax does not do any sort of log rotation.

* `use-global-auth` * - Whether or not to use the set of global secrets for
  authentication of client commands. See [Authentication][auth].

* `secret` * - A secret which will be used for [authentication][auth], assuming
  `use-global-auth` is set to true. Can be specified 0 or more times.

* `use-key-auth` * - Whether or not to check each key a client is modifying for
  a set of secrets to [authenticate][auth] against.

[releases]: https://github.com/mediocregopher/hyrax/releases
[goat]: https://github.com/mediocregopher/goat
[topology]: /doc/topology-examples.md
[auth]: /doc/auth.md
