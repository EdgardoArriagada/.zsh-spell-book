# .zsh-spell-book

- Easily split your bash aliases and functions into different files and folders like you do with your code in your best projects!
- Create a git repo out of this template and enjoy in all your machines your custom aliases, custom functions, custom greetings, program configurations, and more!
- Never lose again any of your custom shell magic!
- Easy to install and uninstall, just source or remove source from .zshrc script
- Uninstalling leaves no residues!

### Prerequisites

- zsh installed and configured as default shell
- [sd](https://github.com/chmln/sd) for bundling
- some commands utilize external packages like [ripgrep](https://github.com/BurntSushi/ripgrep), [sd](https://github.com/chmln/sd) and [fd](https://github.com/sharkdp/fd)

### Installing

1.- git clone ~/.zsh-spell-book

2.- Bundle the project:

```
source ~/.zsh-spell-book/bundle.zsh
```

3.- source the project

Like this:

```shell
echo "source ~/.zsh-spell-book/result.zsh" > ~/.zshrc
```

### Configuring

1.- You should copy `.env.example` to `.env` at the root of this repository
2.- You should go to `~/.zsh-spell-book/src/setup` and source files according to your os

then, in ~/temp/.zshenv paste the following depending on your OS

```zsh
# MacOS
ZSB_MACOS=1
export PATH=~/.zsh-spell-book/src/charms:$PATH

# the rest of OSs
export PATH=~/.zsh-spell-book/src/charms:$PATH
```

### Development

When developing new aliases and function, it is convenient to do the following modification:

```diff
- source ~/.zsh-spell-book/result.zsh
+ source ~/.zsh-spell-book/development.zsh
```

So you get refreshed code every time you open a new terminal

When finished, bundle project again and change back `development.zsh` -> `result.zsh` in your `~/.zshrc`

### Usage

- Explore the _src/spells_ folder to see a list of aliases and function available, or use that folder as a template to build your custom aliases and functions

- Use the _src/utils_ folder to add functions that can help other functions inside _src_ folder.

- If you want to add temporal (ignored by git) aliases, configurations, functions, variable exports, etc. You can run the _temp_ command and create/edit new files and folders with aliases as you which. Once you finish editing, you can restart your terminal to see your new temporal magic.

- Use the _src/automatic-calls_ folder to add anything you want to be executed and printed to console right after a new terminal is initialized

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
