`domusic` tool suite for the music library at <http://svenax.net>
=================================================================

A tool set to handle the bagpipe music archive at <http://svenax.net>.
It is mostly a rewrite of the old `makelily` Python code, but with a lot
of fixes and cleanup.

Dependencies
------------

- You need to install the program in your GOPATH as usual.
- Lilypond must be installed with the command line executable accessible from your shell path.
- Mogrify from ImageMagick must be installed if you want to use the `crop` flag in `make`.

Configuration
-------------

Global configuration is set in in a file called `~/.domusic.yaml`. You need
to create this yourself. It can look something like this.

    root: "/path/to/Bagpipemusic"
    ly-editor: "code -r"
    ly-viewer: "Preview"

Usage
-----

Use `domusic help` and `domusic help [command]` for more information.
