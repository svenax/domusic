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

For the `sync` command, you also need to configure the remote server settings:

    sync-server: "your-server.com"
    sync-user: "your-username"
    sync-path: "/var/www/html/music/"
    sync-ssh-key: "~/.ssh/your_key"  # Optional, defaults to ~/.ssh/id_rsa

You can also set default include and exclude patterns:

    sync-include:
      - "*.pdf"
      - "*.png"
    sync-exclude:
      - "*.tmp"
      - "*.log"
      - ".DS_Store"

See `example.domusic.yaml` for a complete configuration example.

Usage
-----

Use `domusic help` and `domusic help [command]` for more information.
