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
the sha512 output of the concatenation of the domain and one of the secret keys on the server.

This allows an external service (such as an api) to authenticate everything that your connected
clients are allowed to do.

## Redis command syntax

Hyrax is a layer in between the world and redis. As such almost all redis commands are available for usage.
Most redis commands take the form of `command key [ value ... ]`. The translated form would look like:

```json
{ "command":"____", "payload":[ { "domain":"____", "id":"____", "secret":"____", "values":[ "____","...." ]} ]}
```

Values can be empty (or ommitted), and the values in it can be either strings or numbers, depending on
what's called for by the command. Secrets can also be ommitted if the command doesn't actually alter
anything. The payload can contain multiple key/val items as well.

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
{ "command":"set", "payload":[ { "domain":"td","id":"tid","secret":"lotsahex","values":["tv"] } ]}
{ "command":"set", "return":[ 1 ]}
```

Mset:
```json
{ "command":"mset", "payload":[ { "domain":"td1","id":"tid1","secret":"lotsahex1","values":["tv1"] },
                                { "domain":"td2","id":"tid2","secret":"lotsahex2","values":["tv2"] } ]}
{ "command":"mset", "return":[ 1,1 ]}
```

Getrange:
```json
{ "command":"getrange", "payload":[ { "domain":"td","id":"tid","values":[0,-4]} ]}
```

## Non-redis command syntax

There are other commands that don't have a direct correlation to redis as well. They are documented
here

### Monitor

Monitor a list of keys, and receive updates when those keys change. The command's immediate return includes
the current values of those keys, or an empty value if the key hasn't been set. There is a separate monitor command
for each redis data-type:
* `mon`: For normal strings (`get`,`set`,etc...)
* `hmon`: For hashes (`hget`,`hset`,etc...)
* `lmon`: For lists (`lindex`,`lset`,etc...)
* `smon`: For sets (`sadd`,`sismember`,etc...)
* `zmon`: For sorted sets (`zadd`,`zismemberq,etc...)
* `emon`: For an ekg (see below)

There is also a generic monitor, which doesn't return the current value of the key, but will keep you updated
on the changes to the key: `amon`

Push messages about keys that you're monitoring will merely contain the command used to update them. For example:
```json
{ "command":"mon-push", "return":{ "key":{"domain":"td","id":"tid"}, "values":["whatever"], "command":"set" }}
```

Here's some examples of the individual *mon* commands and what they return (note, for all these examples the second
key hasn't been set yet):

mon:
```json
{ "command":"mon", "payload":[ { "domain":"td1","id":"tid1" }, { "domain":"td2","id":"tid2" } ]}
{ "command":"mon", "return":[ "foo", null ]}
```

hmon:
```json
{ "command":"hmon", "payload":[ { "domain":"td1","id":"tid1" }, { "domain":"td2","id":"tid2" } ]}
{ "command":"hmon", "return":[ { "a":"foo","b":"bar","c":"baz"}, {} ]}
```

lmon:
```json
{ "command":"lmon", "payload":[ { "domain":"td1","id":"tid1" }, { "domain":"td2","id":"tid2" } ]}
{ "command":"lmon", "return":[ ["a","b","c"], [] ]}
```


smon:
```json
{ "command":"smon", "payload":[ { "domain":"td1","id":"tid1" }, { "domain":"td2","id":"tid2" } ]}
{ "command":"smon", "return":[ ["a","b","c"], [] ]}
```

zmon (the return is a map of values to their weight as an integer. It's not exactly pretty since it
basically unorders the set, but that's json for you):
```json
{ "command":"zmon", "payload":[ { "domain":"td1","id":"tid1" }, { "domain":"td2","id":"tid2" } ]}
{ "command":"zmon", "return":[ {"a":1,"b":2,"c":3}, {} ]}
```

amon:
```json
{ "command":"amon", "payload":[ { "domain":"td1","id":"tid1" }, { "domain":"td2","id":"tid2" } ]}
{ "command":"amon", "return":[ 1,1 ] }
```

emon:
```json
{ "command":"emon", "payload":[ { "domain":"td1","id":"tid1" }, { "domain":"td2","id":"tid2" } ]}
{ "command":"emon", "return":[ ["mathew","mark","luke","john"], [] ]}
```

For all of the above commands (except `amon`) if you try to monitor a key the contains a different
type then the one associated with your *mon command the return for that key will be an empty value.

### EKG

EKG's are the little heartbeat monitors that doctors use to track your heart rhythm and stuff. In hyrax an EKG is
a value that will alert others monitoring that value that you have disconnected or removed your ekg. There are
four ekg-altering commands:
* `eadd`: Add yourself to an ekg. Anyone monitoring the ekg will see that you have been added. If your connection
          becomes disconnected anyone monitoring the ekg will receive a notice about that as well.
* `eaddq`: Quietly add yourself to the ekg, anyone monitoring the ekg will NOT see that you've added yourself.
           Disconnecting will still alert anyone monitoring that you've done so.
* `erem`: Remove yourself from an ekg. Anyone monitoring the ekg will see that you've removed yourself, and will
          no longer receive an update in the event of a disconnect.
* `eremq`: Quietly remove yourself from an ekg, anyone monitoring will NOT see that you've removed yourself, and
           disconnecting will no longer generate an alert.

The syntax for ekg-altering commands is different from previous ones due the need to identify yourself to an
ekg value. The syntax looks like this:
```json
{ "command":"eadd", "payload":[ { "domain":"td1","id":"tid1","name":"joseph","secret":"lotsahex" } ]}
{ "command":"eadd", "result": [ 1 ]}
```

(If the ekg already has that value added by another connection the return will be `0` and there will be no change)

The `name` field has been added as a sort of identifier for the connection. For ekg commands the secret is the hash
of the concatenation of the domain, the name, and one of the secret keys on the server.

When monitored the ekg will send out push alerts when a connection either removes a value from an ekg or when it
disconnects. If the same connection has added multiple different values to the same ekg, that ekg will generate a
separate push alert for each value added. Here's an example of what a monitored ekg would send out on a disconnect:

```json
{ "command":"mon-push", "return":{ "key":{"domain":"td","id":"tid","name":"joseph"}, "values":[], "command":"disconnect" }}
```
