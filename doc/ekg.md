# EKG

An EKG is an entity which is identified by a key. When clients monitor an EKG
key, they will receive an alert whenever a client adds itself to the EKG, is
disconnected from hyrax after adding itself to the EKG, or manually removes
itself from the EKG. When interacting with an EKG, clients specify an id for
themselves so they can be identified by other clients.

# Example

Client A:

```json
> {"cmd":"madd","key":"foo"}
< {"return":1}
> {"cmd":"madd","key":"bar"}
< {"return":1}
```

Client B:

```json
> {"cmd":"eadd","key":"foo","id":"mediocre","secret":"<hmac-sha1>"}
< {"return":1}
> {"cmd":"eadd","key":"bar","id":"gopher","secret":"<hmac-sha1>"}
< {"return":1}
```

Client A:

```json
< {"cmd":"eadd","key":"foo","id":"mediocre"}
< {"cmd":"eadd","key":"bar","id":"gopher"}
```

Client B:

```json
> {"cmd":"erem","key":"foo","id":"mediocre","secret":"<hmac-sha1>"}
< {"return":1}
<disconnect>
```

Client A:

```json
< {"cmd":"erem","key":"foo","id":"mediocre"}
< {"cmd":"eclose","key":"bar","id":"gopher"}
```

# Commands

The following are the commands used to interact with EKGs:

## eadd
**modifies: true**

Adds the id given by the client to the EKG named by `key`, creating the key if
it didn't previously exist. A client adding itself to an EKG twice has no effect
except to overwrite the `id` sent in the first `eadd`.

Example:

```json
> {"cmd":"eadd","key":"foo","id":"gopher","secret":"<hmac-sha1>"}
< {"return":"OK"}
```

## erem
**modifies: true**

Removes the client from the EKG named by `key`. It is not necessary to specify
`id`, but you may still want to if any clients [monitoring][mon] the EKG want to
know who is calling `erem`.

Example:

```json
> {"cmd":"erem","key":"foo","secret":"<hmac-sha1>"}
< {"return":"OK"}
```

## emembers
**modifies: false**

Returns the list of clients who are currently connected to hyrax and added to
the EKG named by `key`, or empty list if the EKG is not present.

**Note that this only returns the `id`s of clients connected to the current
node. To retrieve the EKG's list of clients across the cluster you will have to
ask every node in the cluster**

Example:

```json
> {"cmd":"emembers","key":"foo"}
< {"return":["mediocre","gopher"]}
```

## ecard
**modifies: false**

Returns the number of clients who are currently connected to hyrax and added to
the EKG named by `key`, or `0` if the EKG is not present.

**Note that this only returns the count of clients connected to the current
node. To retrieve the count of clients in the EKG across the cluster you will
have to ask every node in the cluster**

Example:

```json
> {"cmd":"ecard","key":"foo"}
< {"return":2}
```

# Caveats

EKGs are in a weird state where they are not actually stored in the backend
datastore, but still share the same key-space. As a result it is possible to do
normal commands against the datastore on a key which is already an EKG, and the
commands will be completed as if the EKG doesn't exist (and vice-versa).

Essentially, the same key can have two different values. The only reason this is
the case is because I haven't thought of an elegant way to fix it. Do not rely
on this behavior!!!

[mon]: /doc/mon.md
