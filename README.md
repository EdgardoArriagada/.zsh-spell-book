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

### Usage

- Explore the _src/spells_ folder to see a list of aliases and function available, or use that folder as a template to build your custom aliases and functions

- Use the _src/utils_ folder to add functions that can help other functions inside _src_ folder.

- You should add your Linux application configurations on _src/configurations_ folder

- If you want to add temporal (ignored by git) aliases, configurations, functions, variable exports, etc. You can run the _temp_ command and create/edit new files and folders with aliases as you which. Once you finish editing, you can restart your terminal or run the `resource` command to see your new temporal magic.

- Use the _src/automatic-calls_ folder to add anything you want to be executed and printed to console right after a new terminal is initialized

### Dynamic prefixing

To use dynamic prefixing, just use the `__${zsb}_` variable to prefix function names. I.E:

```shell
  __${zsb}_example () { ... }
```

- All these functions will be deleted after _main.zsh_ is finished

- To avoid clashes, you can change the prefix string in the _.env_ file at the root of the project modifying the `zsb` global variable. _(see .env.example file)_

- _(Optional)_ All functions that are mean to be used inside scripts only and need to be kept after _main.zsh_ is finished should start with `${zsb}_` instead

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details
