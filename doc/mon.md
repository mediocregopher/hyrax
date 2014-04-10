# Monitors

Clients can monitor arbitrary keys for changes. Any command in the cluster which
modifies a key will generate a push message to all clients in the cluster
monitoring that key, with the push message containing the entire command which
was executed.

# Example

In this example Client A is monitoring a key, Client B updates that key, and
Client A subsequently sees a push message for that update.

Client A:

```json
> {"cmd":"madd","key":"foo"}
< {"return":"OK"}
```

Client B:

```json
> {"cmd":"set","key":"foo","args":["bar"],"id":"gopher","secret":"<hmac-sha1>"}
< {"return":"OK"}
```

Client A:

```json
< {"cmd":"set","key":"foo","args":["bar"],"id":"gopher"}
```

# Commands

The following are the commands used to interact with monitors

## madd
**modifies: false**

Adds `key` to the set of keys the client is monitoring.

Example:

```json
> {"cmd":"madd","key":"foo"}
< {"return":"OK"}
```

# Commands

## mrem
**modifies: false**

Removes `key` from the set of keys the client is monitoring.

Example:

```json
> {"cmd":"mrem","key":"foo"}
< {"return":"OK"}
```

## mlocal
**modifies: false**  
**requires admin: true**

Upon calling this, the client will receive all key change events which occur on
the node the client is connected to. In effect, these will be all the key change
events the node is pushing to the nodes listed in the [configuraiton][config].

Example:

```json
> {"cmd":"mlocal","secret":"<hmac-sha1>"}
< {"return":"OK"}
```

## mglobal
**modifies: false**  
**requires admin: true**

Upon calling this, the client will recevie all key change events which occur in
the entire cluster. In effect, these will be all the key change events the node
is pulling from the nodes listed in the [configuration][config].

Example:

```json
> {"cmd":"mglobal","secret":"<hmac-sha1>"}
< {"return":"OK"}
```

[config]: /doc/installconfig.md
