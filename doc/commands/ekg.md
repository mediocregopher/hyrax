# EKG

EKGs are used as a set of clients which are currently connected to hyrax. A
client adds itself to an EKG with the given id, and is automatically removed
from the set when it disconnects for any reason. Other clients monitoring the
EKG will receive a push message with the command being `eclose` when a client is
removed from an EKG due to disconnection.

# Commands

The following are the commands used to interact with EKGs:

## eadd
**modifies: true**

Adds the id given by the client to the EKG named by `key`.

Example:

```json
> {"cmd":"eadd","key":"foo","id":"gopher","secret":"<hmac-sha1>"}
< {"return":1}
```

## erem
**modifies: true**

Removes the id given by the client from the EKG named by `key`.

Example:

```json
> {"cmd":"erem","key":"foo","id":"gopher","secret":"<hmac-sha1>"}
< {"return":1}
```

## emembers
**modifies: false**

Returns the list of clients who are currently connected to hyrax and added to
the EKG named by `key`, or empty list if the EKG is not present.

Example:

```json
> {"cmd":"emembers","key":"foo"}
< {"return":["mediocre","gopher"]}
```

## ecard
**modifies: false**

Returns the number of clients who are currently connected to hyrax and added to
the EKG named by `key`, or `0` if the EKG is not present.

Example:

```json
> {"cmd":"ecard","key":"foo"}
< {"return":2}
```

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

# Caveats

EKGs are in a weird state where they are stored in the backend data-store, but
they are in a different namespace than keys interracted with using
[direct commands][direct], so direct commands like `del` and `expire` can't
effect them.

It also means that you could theoretically use the same key for two
different values, one an EKG and the other through a direct command. The
instance I can think of how to efficiently fix this bug I will, so don't rely on
it for anything!

[direct]: /doc/commands/direct.md
