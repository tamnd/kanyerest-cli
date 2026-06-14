---
title: "Resource URIs"
description: "Use kanyerest as a database/sql-style driver so a host program can address kanyerest as kanyerest:// URIs."
weight: 20
---

`kanyerest` is a command line, but the `kanyerest` Go package is also a
small driver that makes kanyerest addressable as a resource URI. A host
program registers it the way a program registers a database driver with
`database/sql`, then dereferences `kanyerest://` URIs without knowing
anything about how kanyerest is fetched.

The host that does this today is [ant](https://github.com/tamnd/ant), a single
binary that puts one URI namespace over a family of site tools. The examples
below use `ant`; any program that links the package gets the same behaviour.

## Mounting the driver

A host enables the driver with one blank import, exactly like `import _
"github.com/lib/pq"`:

```go
import _ "github.com/tamnd/kanyerest-cli/kanyerest"
```

The package's `init` registers a domain with the scheme `kanyerest` for the
host `kanyerest.com`. The standalone `kanyerest` binary does not change.

## Addressing records

A URI is `scheme://authority/id`. The scaffold ships one type:

| URI                              | What it is                              |
| -------------------------------- | --------------------------------------- |
| `kanyerest://page/<path>`    | a page, keyed by its path on kanyerest.com |

```bash
ant get kanyerest://page/<path>    # the page record
ant cat kanyerest://page/<path>    # just the body text
ant url kanyerest://page/<path>    # the live https URL
ant resolve https://kanyerest.com/<path> # a pasted link, back to its URI
```

As you add resolver operations in `kanyerest/domain.go`, each new `URIType`
becomes another addressable authority here, with no extra wiring. See
[add a command](/guides/adding-a-command/).

## Walking the graph

`ls` lists the members of a collection, and every member is itself an
addressable URI, so a host can follow the graph and write it to disk:

```bash
ant ls     kanyerest://page/<path>             # the pages this one links to
ant export kanyerest://page/<path> --follow 1 --to ./data
```

The example `links` op emits page stubs, so each listed member is a
`kanyerest://page/` URI in its own right. When you model edges between your
real records with `kit:"link"` tags, `ant export --follow` and `ant graph` walk
those edges too, across tools when a link points at another site's scheme.

## Why this is the same code

The driver and the binary share one definition per operation. A resolver op
answers both `kanyerest page` on the command line and `ant get
kanyerest://page/...` through a host, from the same handler and the same
client. There is no second implementation to keep in step.
