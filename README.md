# fvtt-packs

It's **NOT** usable for now, it's a completely WIP project for me, to improve my knowledge on Go.

---

Manage your compendium packs with ease. You can unpack the LevelDB to get human-readable files, modify them
then pack them again into LevelDB files.

You can call this utility from everywhere. Go into your system/module directory then launch the appropriate command!

For example:

`fvtt-packs unpack`

This will unpack all the LevelDB which are in the packs directory into a _packs_sources directory, containing
human-readable files.

`fvtt-packs pack`

This will pack all the human-readable files inside _packs_sources directory into the packs directory.

---

Flags can be used to customize the tools.

Usage:

`fvtt-packs [command]`

Available Commands:

* `help` Help about any command
* `unpack` Unpack LevelDB into human-readable files

Flags:

* `-h`, `--help` help for fvtt-packs

Use `fvtt-packs [command] --help` for more information about a command.
