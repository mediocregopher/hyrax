# Syntax

Hyrax is built to support many multiple syntaxes, or data formats, although
currently only one (json) is implemented. Regardless, all formats have the same
five potential fields:

* `cmd` - the command being executed
* `key` - the key the command is executed on
* `args` - a list of extra arguments to the command, if needed
* `id` - if the client wants to identify itself for [monitor pushes][mon] or
  [ekgs][ekg] it does so with id
* `secret` - used to [authenticate][auth] that the client can run the command it
  is trying to run

# Usage

When adding a `listen-addr` in the configuration you can use any one of the
following as valid formats (the middle value):

* json

MORE TO COME

## JSON

The json format is more or less a direct translation of the keys given above.
Here's an example of a command formatted with json:

```json
{"cmd":"set","key":"foo","args":["bar"],"id":"gopher","secret":"0ea1e43b0c907b9aea16657bd3e18855a7cf4365"}
```

Any fields which are blank/nil/empty can be ommitted.

Note: If the [protocol][proto] being used is just the tcp protocol, then a
newline must be appended to all json messages, and a newline will be appended to
all that are sent from hyrax.

[mon]: /doc/commands/mon.md
[ekg]: /doc/commands/ekg.md
[auth]: /doc/auth.md
[proto]: /doc/proto.md
