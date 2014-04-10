# Admin

These are miscellaneous administrative commands which can be run. There are a
few others in other sections, such as [monitors](/doc/mon.md).

# Commands

## alistentome
**requires admin: true**

For the most part an internal command, it will probably never be called
manually. It's used so that one node can tell another node that the other should
call `mlocal` on it. The other node will expect that `alistentome` will be
called repeatedly, or else it will close the `mlocal` connection.

Example:

```json
> {"cmd":"alistentome","secret":"<hmac-sha1>"}
< {"return":"OK"}
```

## aignoreme
**requres admin: true**

Used by one node to tell another to stop calling `mlocal` on it.

Example:

```json
> {"cmd":"aignoreme","secret":"<hmac-sha1>"}
< {"return":"OK"}
```

## aglobalsecrets
**requires admin: true**

Returns the list of global secrets in effect for a node.

Example:

```json
> {"cmd":"aglobalsecrets","secret":"<hmac-sha1>"}
< {"return":["fee","fye","foh","fum"]}
```

## asecretsset
**requires admin: true**

Sets per-key secrets for a particular key (see [admin][admin]). Overwrites any
previously set list of secrets. Can be an empty list to remove all keys.

Example:

```json
> {"cmd":"asecretsset","key":"foo","args":["fee","fye","foh","fum"],"secret":"<hmac-sha1>"}
< {"return":"OK"}
```

## asecretsadd
**requires admin: true**

Adds per-key secrets for a particular key (see [admin][admin]).

```json
> {"cmd":"asecretsadd","key":"foo","args":["fee","fye"],"secret":"<hmac-sha1>"}
< {"return":"OK"}
```

## asecretsrem
**requires admin: true**

Removes per-key secrets for a particular key (see [admin][admin]).

```json
> {"cmd":"asecretsrem","key":"foo","args":["for","fum"],"secret":"<hmac-sha1>"}
< {"return":"OK"}
```

## asecrets
**requires admin: true**

Returns the list of currently active per-key secrets for a particular key (see
[admin][admin]).

```json
> {"cmd":"asecrets","key":"foo","secret":"<hmac-sha1>"}
< {"return":["fee","fye","foh","fum"]}
```

[admin]: /doc/admin.md
