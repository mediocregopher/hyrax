# Hyrax

A clustered key-value store, built as a thin wrapper on top of redis, with
real-time updates on key changes and client events.

Hyrax gives you redis for web (and other) clients, with redis-cluster-like
scaling and more functionality. Here's some things you can do:

* Have clients receive real-time updates when a key changes [(here)][mon]
* Have clients receive real-time updates when another client disconnects
  [(here)][ekg]
* Scale to many nodes, and hold many concurrent client connections
* Control what commands clients are allowed to call, and in what context
  [(here)][auth]
* Easily access hyrax from a browser

Hyrax is still under active development, and things are likely to change in the
near future. But please feel free to poke around the code and play with it, and
let me know (at the email in my profile) if there's anything I can do to make
things clearer.

## Table Of Contents

* [Syntax][syntax] - The language(s) used to communicate with hyrax
* [Protocols][proto] - The channels over which hyrax can communicate
* [Commands][commands] - Commands which hyrax supports
* [Authentication][auth] - Restricting client behavior

## Contact

My name is Brian Picciano. You can get ahold of me at the email in my github
profile (github.com/mediocregopher).

[mon]: /doc/commands/mon.md
[ekg]: /doc/commands/ekg.md
[auth]: /doc/auth.md
[syntax]: /doc/syntax.md
[proto]: /doc/proto.md
[commands]: /doc/commands/commands.md

