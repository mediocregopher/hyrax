# Protocols/Syntaxes

Protocols describe the pathway data is transferred to/from hyrax over (e.g.
tcp), while syntaxes are the format that data takes (e.g. json). In all cases
the [basic][basics] data structures being transferred are the same, they
just take different forms.

*Note that there are currently limited options. This is due to not having
written them in yet, and not any significant added complexity they pose to
hyrax. Hyrax is designed to be agnostic to the protocol/syntax used.*

## Protocols

### TCP

There is currently only one supported protocol, and it is tcp. You can specify a
tcp listen endpoint in the [configuration][config] by setting the protocol to
`tcp`.

*Note that for all formats tcp will terminate messages it sends with a newline.
It will expect all received messages to be terminated similarly*

## Syntaxes

### JSON

There is currently only one supported syntax, and it is json. You can specify a
json listen endpoint in the [configuration][config] by setting the format to
`json`. Here are examples of Action and ActionReturn structs (see
[basics][basics]) as json. These have been pretty formatted, when actually
communicating with hyrax remove all whitespace and newlines:

Action:

```json
{
    "cmd":"SET",
    "key":"foo",
    "args":["bar"],
    "id":"mediocregopher",
    "secret":"225711f795d512fef53aef38939813163bae3462"
}
```

ActionReturn (successful return and error):

```json
{
    "return":"OK"
}
{
    "error":"not enough pylons"
}
```

[basics]: /doc/basics.md
[config]: /doc/installconfig.md
