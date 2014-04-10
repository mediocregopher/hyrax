# Auth

Authentication is done on a command-by-command basis. Any hyrax command which
modifies the value of its key must be authenticated.  Additionally, some
commands are marked as [admin][admin] commands and must be authenticated in any
case.

Hyrax's authentication is based around secret keys which are shared with the
hyrax node itself and the actual backend of the application which handles the
application logic. When clients need to perform actions they communicate with
the backend to obtain secret hashes for those actions. These secret hashes
are then used to authenticate with hyrax.

The secret keys are specified on a per-node basis. Hyrax nodes do not
communicate with each other about keys.

*Note that this document only applies if either `use-global-auth` or
`use-key-auth` (or both) is set to true in the [configuration][config]*

## Method

A secret is an arbitrary string of characters. For every command a set of
potential secrets is determined (from the global pool and the per-key pool) and
the hyrax checks that the command authenticates with one of those secrets.
Authentication is done by running the following algorithm, and checking if the
algorithm matches the `Secret` field on the command:

```
HexEncode(HmacSHA1(secret + command + key + id))
```

Where `secret` is one of the secrets from the set. The set is ordered in the
same way as the following sections.

## Global secrets

Global secrets are defined in the [configuration][config] of every node. These
are always checked, regardless of the key being acted upon. They are useful for
backend, trusted processes to use.

## Per-key secrets

It is also possible to set secrets on individual keys, using the ASECRETADD and
associated [admin][admin] commands. The per-key secrets for a key will only be
checked when that key is being acted upon. These are useful when you want to
create revokable permissions for individual clients.

[config]: /doc/installconfig.md
[admin]: /doc/admin.md
