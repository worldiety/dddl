// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from "vscode";
import * as os from "os";
import * as fs from "fs";
import { LanguageClient, LanguageClientOptions, ServerOptions, TransportKind } from "vscode-languageclient/node";
import { Builder } from 'selenium-webdriver';
import * as chrome from 'selenium-webdriver/chrome';

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


    context.subscriptions.push(vscode.commands.registerCommand("ddd.ExportHTML", () => {
        client.sendRequest("custom/ExportHTML", null).then((resp) => {
            vscode.workspace.openTextDocument({
                content: String(resp),
                language: "html"
            }).then((document) => {
                vscode.window.showTextDocument(document);
            });
        });
    }));

    context.subscriptions.push(vscode.commands.registerCommand("ddd.ExportPDF", async () => {
        const filePath = vscode.window.activeTextEditor?.document.uri.fsPath;
        let doc = "file://" + filePath;
        if (doc) {
            const options = new chrome.Options();
            // options.addArguments('--headless');
            options.addArguments('--kiosk-printing');
    
            const driver = await new Builder().forBrowser('chrome').setChromeOptions(options).build();
    
            try {
                await driver.get(doc);
                await driver.executeScript('window.print();');
            } catch (error) {
                console.log(error);
                vscode.window.showErrorMessage('PDF could not be generated')
            } finally {
                driver.quit();
            }
            vscode.window.showInformationMessage(`PDF was generated in your default downloads folder`)
        }
    }));

    context.subscriptions.push(vscode.commands.registerCommand("ddd.GenerateGo", () => {

        client.sendRequest("custom/GenerateGo", null).then((resp) => {

        });
    }));

    //preview???


    client.onNotification("custom/newAsyncPreviewHtml", resp => {
        console.log("shall async preview")

        if (typeof resp === "string") {
            if (PreviewPanel.currentPanel != null) {
                if (resp==="lastPreviewParams missing"){
                    let tailwindUri = PreviewPanel.currentPanel?._tailwindUri;
                    let webviewPrefix =  PreviewPanel.currentPanel?._webviewPrefixUri;
                    client.sendRequest("custom/webViewParams", {WebviewWorkspacePrefix:webviewPrefix?.toString(),TailwindUri: tailwindUri?.toString()}).catch(e => console.log(e))

                }

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
            client.sendRequest("custom/previewHTML", {doc: doc, tailwindUri: tailwindUri?.toString()}).then((resp) => {
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



function getWebviewOptions(extensionUri: vscode.Uri): vscode.WebviewOptions {
    let localWorkspaceUri = vscode.Uri.file("");

    if(vscode.workspace.workspaceFolders !== undefined) {
        localWorkspaceUri = vscode.workspace.workspaceFolders[0].uri;
    }
    

    return {
        // Enable javascript in the webview
        enableScripts: true,

        // And restrict the webview to only loading content from our extension's `media` directory.
        //localResourceRoots: [vscode.Uri.joinPath(extensionUri, 'media')]


       
        
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
    public readonly _webviewPrefixUri:vscode.Uri;

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
        PreviewPanel.currentPanel = new PreviewPanel(panel, extensionUri, "...");
        //PreviewPanel.createOrShow(extensionUri, "<p>Einen Moment bitte...</p>");

        let tailwindUri = PreviewPanel.currentPanel?._tailwindUri;
        let webviewPrefix =  PreviewPanel.currentPanel?._webviewPrefixUri;
        client.sendRequest("custom/webViewParams", {WebviewWorkspacePrefix:webviewPrefix?.toString(),TailwindUri: tailwindUri?.toString()}).catch(e => console.log(e))

    }

    private constructor(panel: vscode.WebviewPanel, extensionUri: vscode.Uri, html: string) {
        this._panel = panel;
        this._extensionUri = extensionUri;
        this._html = html;

        this._tailwindUri = panel.webview.asWebviewUri(vscode.Uri.joinPath(this._extensionUri, 'media', 'tailwind.js'));
        if(vscode.workspace.workspaceFolders !== undefined) {
            this._webviewPrefixUri = panel.webview.asWebviewUri(vscode.workspace.workspaceFolders[0].uri);
        }else{
            this._webviewPrefixUri=vscode.Uri.file("")
        }
        
        

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
        webview.html = this._html
    }


}

