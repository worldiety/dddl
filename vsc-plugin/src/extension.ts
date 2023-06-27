// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from "vscode";
import * as os from "os";
import * as fs from "fs";
import {LanguageClient, LanguageClientOptions, ServerOptions, TransportKind} from "vscode-languageclient/node";

let client: LanguageClient;

async function sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

export async function activate(context: vscode.ExtensionContext) {

    // Select correct language server binary for this platform.
    let platform = `${os.platform()}-${os.arch()}`;
    if (platform === "darwin-arm64") {
        // Workaround while we cannot build native binary on github actions.
        platform = "darwin-x64";
    }
    let binPath = context.asAbsolutePath(`out/bin/dddlsp-${platform}`);
    if (!fs.existsSync(binPath)) {
        vscode.window.showErrorMessage(`wdyspec-support has no binary for platform "${platform}" and will not work. Contact the developer to fix this.`);
        return;
    }

    let serverOptions: ServerOptions = {
        command: binPath,
        transport: TransportKind.stdio,
    };
    let clientOptions: LanguageClientOptions = {
        documentSelector: [{scheme: "file", language: "ddd"}],
    };
    client = new LanguageClient(
        "ddd-language-server",
        "worldiety ddd Language Server",
        serverOptions,
        clientOptions
    );
    await client.start();

   // await sleep(1000);


    // Request an XML preview from the language server and show that result in a new editor.
    context.subscriptions.push(vscode.commands.registerCommand("dddl.exportAsciiDoc", () => {
        let doc = "file://" + vscode.window.activeTextEditor?.document.uri.fsPath;
        if (doc) {
            client.sendRequest("custom/exportAsciiDoc", doc).then((resp) => {
                vscode.workspace.openTextDocument({
                    content: String(resp),
                    language: "asciidoc"
                }).then((document) => {
                    vscode.window.showTextDocument(document);
                });
            });
        }
    }));


    //preview???



    client.onNotification("custom/newAsyncPreviewHtml",resp => {
        console.log("shall async preview")

        if (typeof resp === "string") {
            if (PreviewPanel.currentPanel!=null){
                PreviewPanel.currentPanel._html = resp
                PreviewPanel.currentPanel._update() // do not capture focus
            }
        } else {
            console.log("cannot handle resp, not a string");
        }
    });

    context.subscriptions.push(vscode.commands.registerCommand("wdyspec.previewHTML", () => {
        let doc = "file://" + vscode.window.activeTextEditor?.document.uri.fsPath;
        if (doc) {
            PreviewPanel.createOrShow(context.extensionUri, "<p>Einen Moment bitte...</p>");

            let tailwindUri = PreviewPanel.currentPanel?._tailwindUri;
            client.sendRequest("custom/previewHTML", {doc:doc,tailwindUri:tailwindUri?.toString()}).then((resp) => {
                console.log("shall preview", resp);

                if (typeof resp === "string") {
                    PreviewPanel.createOrShow(context.extensionUri, resp);
                } else {
                    console.log("cannot handle resp, not a string");
                }
            });
        }
    }));

    if (vscode.window.registerWebviewPanelSerializer) {
        // Make sure we register a serializer in activation event
        vscode.window.registerWebviewPanelSerializer(PreviewPanel.viewType, {
            async deserializeWebviewPanel(webviewPanel: vscode.WebviewPanel, state: any) {
                console.log(`Got state: ${state}`);
                // Reset the webview options so we use latest uri for `localResourceRoots`.
                webviewPanel.webview.options = getWebviewOptions(context.extensionUri);
                PreviewPanel.revive(webviewPanel, context.extensionUri);

                // PreviewPanel.currentPanel._html = state TODO how to do that, when to save, broken example?
            }
        });
    }
}

function getViewColumn(sideBySide: boolean): vscode.ViewColumn | undefined {
    const active = vscode.window.activeTextEditor;
    if (!active) {
        return vscode.ViewColumn.One;
    }

    if (!sideBySide) {
        return active.viewColumn;
    }

    switch (active.viewColumn) {
        case vscode.ViewColumn.One:
            return vscode.ViewColumn.Two;
        case vscode.ViewColumn.Two:
            return vscode.ViewColumn.Three;
    }

    return active.viewColumn;
}

// Shut down language server and close preview panels when extension is deactivated
export function deactivate(): Thenable<void> | undefined {
    if (!client) {
        return undefined;
    } else {
        return client.stop();
    }
}


//================

const cats = {
    'Coding Cat': 'https://media.giphy.com/media/JIX9t2j0ZTN9S/giphy.gif',
    'Compiling Cat': 'https://media.giphy.com/media/mlvseq9yvZhba/giphy.gif',
    'Testing Cat': 'https://media.giphy.com/media/3oriO0OEd9QIDdllqo/giphy.gif'
};


function getWebviewOptions(extensionUri: vscode.Uri): vscode.WebviewOptions {
    return {
        // Enable javascript in the webview
        enableScripts: true,

        // And restrict the webview to only loading content from our extension's `media` directory.
        //localResourceRoots: [vscode.Uri.joinPath(extensionUri, 'media')]
        localResourceRoots: [
            extensionUri
        ]
    };
}

/**
 * Manages cat coding webview panels
 */
class PreviewPanel {
    /**
     * Track the currently panel. Only allow a single panel to exist at a time.
     */
    public static currentPanel: PreviewPanel | undefined;

    public static readonly viewType = 'catCoding';

    private readonly _panel: vscode.WebviewPanel;
    private readonly _extensionUri: vscode.Uri;
    public readonly _tailwindUri: vscode.Uri;
    public _html: string;
    private _disposables: vscode.Disposable[] = [];

    public static createOrShow(extensionUri: vscode.Uri, html: string) {
        const column = getViewColumn(true);

        // If we already have a panel, show it.
        if (PreviewPanel.currentPanel) {
            PreviewPanel.currentPanel._html = html
            PreviewPanel.currentPanel._update()
            PreviewPanel.currentPanel._panel.reveal(column);
            return;
        }

        // Otherwise, create a new panel.
        const panel = vscode.window.createWebviewPanel(
            PreviewPanel.viewType,
            'wdy dddl preview',
            vscode.ViewColumn.Two,
            getWebviewOptions(extensionUri),
        );

        PreviewPanel.currentPanel = new PreviewPanel(panel, extensionUri, html);
    }

    public static revive(panel: vscode.WebviewPanel, extensionUri: vscode.Uri) {
        //PreviewPanel.currentPanel = new PreviewPanel(panel, extensionUri, "TODO cannot revive from nothing");
        PreviewPanel.createOrShow(extensionUri, "<p>Einen Moment bitte...</p>");

        let tailwindUri = PreviewPanel.currentPanel?._tailwindUri;
        client.sendRequest("custom/webViewParams",{TailwindUri:tailwindUri?.toString()}).catch(e=>console.log(e))

    }

    private constructor(panel: vscode.WebviewPanel, extensionUri: vscode.Uri, html: string) {
        this._panel = panel;
        this._extensionUri = extensionUri;
        this._html = html;

        this._tailwindUri = panel.webview.asWebviewUri(vscode.Uri.joinPath(this._extensionUri, 'media', 'tailwind.js'));

        // Set the webview's initial html content
        this._update();

        // Listen for when the panel is disposed
        // This happens when the user closes the panel or when the panel is closed programmatically
        this._panel.onDidDispose(() => this.dispose(), null, this._disposables);

        // Update the content based on view changes
        this._panel.onDidChangeViewState(
            e => {
                if (this._panel.visible) {
                    this._update();
                }
            },
            null,
            this._disposables
        );

        // Handle messages from the webview
        this._panel.webview.onDidReceiveMessage(
            message => {
                switch (message.command) {
                    case 'alert':
                        vscode.window.showErrorMessage(message.text);
                        return;
                }
            },
            null,
            this._disposables
        );
    }

    public doRefactor() {
        // Send a message to the webview webview.
        // You can send any JSON serializable data.
        this._panel.webview.postMessage({command: 'refactor'});
    }

    public dispose() {
        PreviewPanel.currentPanel = undefined;

        // Clean up our resources
        this._panel.dispose();

        while (this._disposables.length) {
            const x = this._disposables.pop();
            if (x) {
                x.dispose();
            }
        }
    }

    public _update() {
        const webview = this._panel.webview;

		webview.html=this._html

       // webview.html = this._getHtmlForWebview(webview,"abc")
        // Vary the webview's content based on where it is located in the editor.
        /*switch (this._panel.viewColumn) {
            case vscode.ViewColumn.Two:
                this._updateForCat(webview, 'Compiling Cat');
                return;

            case vscode.ViewColumn.Three:
                this._updateForCat(webview, 'Testing Cat');
                return;

            case vscode.ViewColumn.One:
            default:
                this._updateForCat(webview, 'Coding Cat');
                return;
        }*/
    }


    private _updateForCat(webview: vscode.Webview, catName: keyof typeof cats) {
        this._panel.title = catName;
        this._panel.webview.html = this._getHtmlForWebview(webview, cats[catName]);
    }

    private _getHtmlForWebview(webview: vscode.Webview, catGifPath: string) {
        // Local path to main script run in the webview
        const scriptPathOnDisk = vscode.Uri.joinPath(this._extensionUri, 'media', 'main.js');

        // And the uri we use to load this script in the webview
        const scriptUri = webview.asWebviewUri(scriptPathOnDisk);

        // Local path to css styles
        const styleResetPath = vscode.Uri.joinPath(this._extensionUri, 'media', 'reset.css');
        const stylesPathMainPath = vscode.Uri.joinPath(this._extensionUri, 'media', 'vscode.css');

        // Uri to load styles into webview
        const stylesResetUri = webview.asWebviewUri(styleResetPath);
        const stylesMainUri = webview.asWebviewUri(stylesPathMainPath);

        const tailwindUri = webview.asWebviewUri(vscode.Uri.joinPath(this._extensionUri, 'media', 'tailwind.js'));

        // Use a nonce to only allow specific scripts to be run
        const nonce = getNonce();

        return `<!DOCTYPE html>
			<html lang="en">
			<head>
				<meta charset="UTF-8">

				<!--
					Use a content security policy to only allow loading images from https or from our extension directory,
					and only allow scripts that have a specific nonce.
				-->
				<meta http-equiv="Content-Security-Policy" content="default-src 'none'; style-src 'self' 'unsafe-inline'; img-src ${webview.cspSource} https:; script-src 'nonce-${nonce}';">

				<meta name="viewport" content="width=device-width, initial-scale=1.0">

			
				 <script src="${tailwindUri}" type="text/javascript" nonce="${nonce}"></script>

				<title>Cat Coding</title>
			</head>
			<body>
				<img src="${catGifPath}" width="300" />
				<h1 class="bg-yellow-500" id="lines-of-code-counter">0asdadsadas</h1>
			
                <p>${tailwindUri}</p>
                 <p>${this._extensionUri}</p>
				<script nonce="${nonce}" src="${scriptUri}"></script>
			</body>
			</html>`;
    }
}

function getNonce() {
    let text = '';
    const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    for (let i = 0; i < 32; i++) {
        text += possible.charAt(Math.floor(Math.random() * possible.length));
    }
    return text;
}