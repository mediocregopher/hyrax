# Auth

Any hyrax command which modifies the value of its key must be authenticated.

Hyrax has a set of global secret keys in its memory, as well as the possiblity
of secret keys for individual keys in the cluster. An authenticated command is
one where the `secret` field in the command is filled with the hmac-sha1 of the
`cmd`, `key`, and `id` fields of the command concatentated together, with a
valid secret (either a global one or a key-specific one) as the secret for the
hmac.

This secret allows a backend service to authenticate what untrusted clients are
allowed to do. The backend service would enumerate what keys it expects a client
to interact with, what commands it expects the client to use, and what id the
client is able to use for those keys, and generates the appropriate secrets
which it then gives the client.

Alternatively, if you don't care about authentication you can simply give the
clients one of the secret keys and let them generate their own secrets.

# Example

A client wants to perform the following command:

```json
{"cmd":"set","key":"foo","args":["bar"],"id":"gopher"}
```

And one of the global secret keys is `toy story 2 was okay`. The backend would
generate the secret as follows:

```
# hmac-sha1(secret,data)
hmac-sha1("toy story 2 was okay", "setfoogopher")
```

and would get back `0ea1e43b0c907b9aea16657bd3e18855a7cf4365`. It would then
give this to the client, and the client could call the command by doing:

```json
{"cmd":"set","key":"foo","args":["bar"],"id":"gopher","secret":"0ea1e43b0c907b9aea16657bd3e18855a7cf4365"}
```
