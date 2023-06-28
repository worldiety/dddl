# dyml-vscode

A VS code extension to add support for [DYML](https://github.com/golangee/dyml). It will feature syntax highlighting and syntax checking.

Currently WIP.

## Development
You need [VS Code](https://code.visualstudio.com/) to try this extension. Run `npm install`, then open this project in vscode and select `Run > Start Debugging` or press `F5`. You need to have a go compiler and npm installed.

It might be useful to see the extensions output, which you can show by opening `View > Output` and selecting it in the dropdown on the right. `DYML Language Server` and `Log (Extension Host)` might be of interest here.

To package the application into a `.vsix` file run `make`. You might need to `npm install -g vsce` first. This package can then be installed in vscode by opening the overflow menu in the extension tab, and selecting `Install from VSIX`.