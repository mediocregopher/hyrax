# Hyrax

A key-val store which send out updates in real-time

## Syntax

Sent:
```json
{ "command":"____", "params":{"key":"value"} }
```

Response:
```json
{ "command":"____", "result":{"____":"____"} }
```
All messages to and from hyrax will be appended with a newline to delimit the end of the message.

## Keys/vals

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

## Commands

### get

Any key can be gotten by anyone.

Sent:
```json
{ "command":"get", "params":{ "key":{"domain":"____","id":"____"}}}
```

Received:
```json
{ "command":"get", "result":{ "value":"____"}}
```

Or if there is an error:
Received:
```json
{ "command":"get", "result":{ "error":"____" }}
```

An example of an error would be ```key_dne```

### set

A key can be set to any string of your choosing (and only a string!).

Sent:
```json
{ "command":"set", "params":{ "key":{"domain":"____","id":"____"}, "value":"____", "secret":"____" }}
```

Received:
```json
{ "command":"set", "result":{ "success":true }}
```

Or if there is an error:
Received:
```json
{ "command":"set", "result":{ "success":false, "error":"____" }}
```

