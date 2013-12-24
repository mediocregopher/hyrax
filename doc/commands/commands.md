# Commands

Hyrax is built as a thin layer on top of a backend data-store (redis). As such,
it supports almost all commands that redis does, as well as a few which it
implements itself.

All commands take in a `cmd` and `key` field. The `key` field is used to route
the command to the appropriate redis node in the cluster.

Commands are broken up into four different sections:

* [direct](/doc/commands/direct.md) - Commands which are passed through to redis
* [monitors](/doc/commands/mon.md) - Used to receive notifications on changes to
  other keys
* [ekgs](/doc/commands/ekg.md) - Used so clients can track the connectivity of
  each other
* [admin](/doc/commands/admin.md) - Commands to administer the hyrax cluster
