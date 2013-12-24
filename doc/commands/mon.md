# Monitors

Monitors can be used to receive alerts when a key in the hyrax cluster is
altered by a client. When monitoring a key and another client performs a command
which modifies that key, hyrax will push to the monitoring client the verbatim
command the modifying client sent (sans the `secret` field).

# Commands

The following are the commands used to add/remove monitors:

## madd
**modifies: false**

Adds `key` to the set of keys the client is monitoring. It will be pushed any
commands other clients perform which modify that key.

Example:

```json
> {"cmd":"madd","key":"foo"}
< {"return":1}
```

## mrem
**modifies: false**

Removes `key` from the set of keys the client is monitoring. The client won't
receive any more notifications about `key`.

Example:

```json
> {"cmd":"mrem","key":"foo"}
< {"return":1}
```

# Example

Client A:

```json
> {"cmd":"madd","key":"foo"}
< {"return":1}
```

Client B:

```json
> {"cmd":"set","key":"foo","args":["bar"],"id":"gopher","secret":"<hmac-sha1>"}
< {"return":"OK"}
```

Client A:

```json
< {"cmd":"set","key":"foo","args":["bar"],"id":"ohai"}
```
