# Basics

This doc covers the basic data structures you will use to interact with hyrax.
It does not discuss the formatting of the actual structures, as that depends on
the [syntax](/doc/protosyntax.md) you choose.

## Action

The basic unit of communication to hyrax looks like this:

```
Action struct {
    Command string
    Key     string
    Args    []Anything // An array containing values of any type
    Id      string
    Secret  string
}
```

`Command` is a string representing the command. This could be passed back to the
storage backend if it is not a builtin command.

`Key` is the key being acted upon. Other clients can monitor any key, receiving
push notifications when that key is modified.

`Args` is an array containing any further information necessary to carry out the
command.

`Id` is an optional field which identifies the client calling the command in
some way. The only purpose of this is to identify the calling client in the
event of other clients monitoring the command's key.

`Secret` is an optional field which is only necessary if
[authentication](/doc/auth.md)
is enabled and the command being called modifies its key's state or is an
[admin](/doc/admin.md) command.

### Action examples

Here's an example of a SET command (assumes that the backend is
[redis][redis], and that the server has a global key `scroopynoopers`)

```
{
    Command: SET
    Key: foo
    Args: [bar]
    Id: mediocregopher
    Secret: 225711f795d512fef53aef38939813163bae3462
}
```

Here's an example of a GET command (assumes that the backend is [redis][redis]).
Note that `Id` and `Secret` aren't set, because GET does not modify its key so
no secret is needed and no key change event will be generated:

```
{
    Command: GET
    Key: foo
    Args: []
    Id:
    Secret:
}
```

## ActionReturn

Once an Action is sent to hyrax an ActionReturn is always sent back. These take
the following form:

```
ActionReturn struct {
    Error  string
    Return Anything
}
```

If `Error` is set than `Return` will be a null or zero value. Otherwise `Return`
will be an appropriate result for whatever the sent Action was.

## Push messages

Hyrax will also push messages to the client at arbitrary times, assuming the
client is monitoring some key or set of keys. Push messages are an exact copy of
the Action which was performed on a monitored key, with the only exception
being that the `Secret` field will be scrubbed out.

[redis]: /doc/redis.md
