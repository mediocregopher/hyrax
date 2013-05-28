# Hyrax

A key-val store which sends out updates in real-time

## Keys

All keys in hyrax actually have two parts to their name: their domain and identifier. Both are strings
of any kind.
```json
{ "domain":"____", "id":"____" }
```

Anyone connected to hyrax has the ability to `get` or `sub` to the value of any key, but only those
who have properly authenticated to the key's domain have the ability to change the key.

## Auth

Hyrax is set up with a list of secret keys. When you send a command which requires authentication to a
particular domain you will also be sending a key in that command which must be the hex string form of
the sha512 output of the domain and one of the secret keys on the server.

This allows an external service (such as an api) to authenticate everything that your connected
clients are allowed to do.

## Syntax

Hyrax is a layer in between the world and redis. As such almost all redis commands are available for usage.
Most redis commands take the form of `command key [ value ... ]`. The translated form would look like:

```json
{ "command":"____", "payload":[ { "domain":"____", "id":"____", "values":[ "____","...." ]} ]}
```

Values can be empty (or ommitted), and the values in it can be either strings or numbers, depending on
what's called for by the command. The payload can contain multiple key/val items as well.

### Command syntax examples

The following are examples of commands (and what they return)

Get:
```json
{ "command":"get", "payload":[ { "domain":"td","id":"tid" } ]}
{ "command":"get", "return":[ "Ohaithar" ]}
```

Mget:
```json
{ "command":"mget", "payload":[ { "domain":"td1","id":"tid1" },
                                { "domain":"td2","id":"tid2" } ]}
{ "command":"mget", "return":[ "Ohaithar",null ]}
```

Set:
```json
{ "command":"set", "payload":[ { "domain":"td","id":"tid","values":["tv"] } ]}
{ "command":"set", "return":[ 1 ]}
```

Mset:
```json
{ "command":"mset", "payload":[ { "domain":"td1","id":"tid1","values":["tv1"] },
                                { "domain":"td2","id":"tid2","values":["tv2"] } ]}
{ "command":"mset", "return":[ 1,1 ]}
```

Getrange:
```json
{ "command":"getrange", "payload":[ { "domain":"td","id":"tid","values":[0,-4]} ]}
```
