`domusic` tool suite for music library handling
===============================================

A tool set to handle the music archive at <http://svenax.net>. It can of course
be used for other music archives as well. It is a complete rewrite of the old
`makelily` Python code, but with a lot of fixes and cleanup.

Dependencies
------------

- You need to install the program in your GOPATH as usual.
- Lilypond must be installed with the command line executable accessible from
  your shell path.
- Mogrify from ImageMagick must be installed if you want to use the `crop` flag
  in `make`.

Configuration
-------------

The configuration file is automatically searched in these locations (in order):

1. `.domusic.yaml` or `.domusic` in current directory and parent directories
   (project-specific)
2. `~/.config/domusic/config.yaml` (XDG config directory)
3. `/etc/xdg/domusic/config.yaml` (system-wide XDG config)
4. `~/.domusic.yaml` (legacy location)
5. `~/.domusic` (legacy location)

You can also specify a custom location with `--config /path/to/config.yaml`.

This search order allows you to have project-specific configurations that
override your global settings.

To use the template definitions you have access to a number of variables. The
`example.domusic.yaml` file uses them all as intended.

Usage
-----

Use `domusic help` and `domusic help [command]` for more information.
