# `theme`

This is my bare minimum theme switcher. It simply copies files and runs hooks &mdash; that's it.
It does not generate color palletes from wallpaper like other theme switchers usually do.

## Installation

The simplest way is to install it using Go:

```shell
go install github.com/shelepuginivan/theme@latest
```

You can also clone this repository and install `theme` system-wide by running:

```shell
git clone "https://github.com/shelepuginivan/theme.git"
cd theme
sudo make install
```

I do not maintain any packages, although prebuilt binaries may be available in the future.

## Example

I use `theme` myself, so you can have a look at my [dotfiles](https://github.com/shelepuginivan/dotfiles/tree/main/.config/theme) for an example.

## Usage

Themes are stored as directories located in **prefix**, which is `$XDG_CONFIG_HOME/theme` (`~/.config/theme` in most cases).

```
$XDG_CONFIG_HOME/theme/
├── my_awesome_theme/
│   ├── copy.json
│   ├── run.d/
│   │   ├── hook_1.sh
│   │   ├── hook_2.sh
│   │   └── hook_3.sh
│   └── other files...
└── other themes...
```

The directory name is the name of the theme.

Nesting is supported, so you can use subdirectories to organize your themes (e.g. `waybar/solid-top` is a valid theme name).
Or you can use another prefix instead, by passing a `-p` flag (see below).

### CLI

```
A very simple theme switcher

Usage:
  theme [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  list        List available themes
  random      Set a random available theme
  set         Set a theme

Flags:
  -h, --help            help for theme
  -p, --prefix string   Directory where themes are stored (default "/home/user/.config/theme")

Use "theme [command] --help" for more information about a command.
```

### `copy.json`

This is a special file where you can specify which files should be copied and where they should be copied to.
Directories are copied recursively and merged. Conflicting files are replaced.

Files are copied concurrently, both in different entries of `copy.json` and within a single entry when copying a directory.

> [!TIP]
> For convenience, `theme` expands paths specified in this file as follows:
> - `@` is replaced with path to the theme (in the example above, `$XDG_CONFIG_HOME/my_awesome_theme`)
> - `~` is replaced with user home directory (`$HOME`)
> - All environment variables are expanded.

Below is an example `copy.json` file.

```json
[
    {
        "src": "@/files/alacritty.toml",
        "dst": "$HOME/.config/alacritty/colors.toml"
    },
    {
        "src": "@/files/dunst",
        "dst": "~/.config/dunst"
    },
    {
        "src": "@/files/gtk.css",
        "dst": "$HOME/.cache/gtk.css"
    }
]
```

### `run.d`

This directory contains executable files (hooks) that run after files are copied.
A common use case is restarting programs to update their configurations:

```shell
#!/bin/sh

if systemctl --user --quiet is-active dunst; then
    systemctl --user restart dunst.service
fi
```

Hooks are executed concurrently.

> [!TIP]
> The files in the `run.d` directory run with the working directory of the theme:
> ```shell
> #!/bin/sh
>
> swww img "./files/wall.png" # this path is resolved to $XDG_CONFIG_HOME/my_awesome_theme/files/wall.png
> ```
