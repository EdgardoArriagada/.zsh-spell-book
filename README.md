# .zsh-spell-book

- Easily split your bash aliases and functions into different files and folders like you do with your code in your best projects!
- Create a git repo out of this template and enjoy in all your machines your custom aliases, custom functions, custom greetings, program configurations, and more!
- Never lose again any of your custom shell magic!
- Easy to install and uninstall, just source or remove source from .zshrc script
- Uninstalling leaves no residues!

### Prerequisites

- zsh installed and configured as default shell

### Installing

Just source the _main.zsh_ file at the end of your _.zshrc_ file

Like this:

```shell
 source ~/.zsh-spell-book/main.zsh
```

### Configuring

1.- You should copy `.env.example` to `.env` at the root of this repository
2.- You should go to `~/.zsh-spell-book/src/setup` and source files according to your os

then, in ~/temp/.zshenv paste the following depending on your OS

```zsh
# MacOS
ZSB_MACOS=1

# empty for the rest of OSs
```

### Usage

- Explore the _src/spells_ folder to see a list of aliases and function available, or use that folder as a template to build your custom aliases and functions

- Use the _src/utils_ folder to add functions that can help other functions inside _src_ folder.

- You should add your Linux application configurations on _src/configurations_ folder

- If you want to add temporal (ignored by git) aliases, configurations, functions, variable exports, etc. You can run the _temp_ command and create/edit new files and folders with aliases as you which. Once you finish editing, you can restart your terminal to see your new temporal magic.

- Use the _src/automatic-calls_ folder to add anything you want to be executed and printed to console right after a new terminal is initialized

- User the _setup_ folder to setup (once per installation) some programs to use configurations stored in this project, then use the `source` command to source these scripts manually

### Dynamic prefixing

To use dynamic prefixing, just use the `__${zsb}.` variable to prefix function names. I.E:

```shell
  __${zsb}.example () { ... }
```

- All these functions will be deleted after _main.zsh_ is finished

- To avoid clashes, you can change the prefix string in the _.env_ file at the root of the project modifying the `zsb` global variable. _(see .env.example file)_

- _(Optional)_ All functions that are mean to be used inside scripts only and need to be kept after _main.zsh_ is finished should start with `${zsb}.` instead

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
