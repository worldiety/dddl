# dddl

For more information and example, visit the [Marketplace](https://marketplace.visualstudio.com/items?itemName=worldiety.dddl).

## dddc
You can install the dddl standalone compiler as follows. 

```bash
go install github.com/worldiety/dddl/cmd/dddc@latest

# within the working directory of the  *.ddd files project
# e.g. to generate a standalone html file
dddc -format=html -out=index.html

# or generate go code (by default places always in mod-root/internal/<context name>
dddc -format=go
```

## Plugin for Visual Studio Code
A VS code extension to add support for [DDDL](https://github.com/worldiety/dddl). 
It will feature syntax highlighting, linting, live preview and more.
See [changelog](vsc-plugin/CHANGELOG.md)

### Development
You need [VS Code](https://code.visualstudio.com/) to try this extension. Run `npm install`, then open this project in vscode and select `Run > Start Debugging` or press `F5`. You need to have a go compiler and npm installed.

It might be useful to see the extensions output, which you can show by opening `View > Output` and selecting it in the dropdown on the right. `DYML Language Server` and `Log (Extension Host)` might be of interest here.

To package the application into a `.vsix` file run `make`. You might need to `npm install -g vsce` first. This package can then be installed in vscode by opening the overflow menu in the extension tab, and selecting `Install from VSIX`.
