// The module 'vscode' contains the VS Code extensibility API
// Import the module and reference it with the alias vscode in your code below
import * as vscode from "vscode";
import * as os from "os";
import * as fs from "fs";
import { LanguageClient, LanguageClientOptions, ServerOptions, TransportKind } from "vscode-languageclient/node";
import puppeteer, { PDFOptions } from "puppeteer";
import * as express from "express";

let client: LanguageClient;

// async function sleep(ms: number) {
//     return new Promise(resolve => setTimeout(resolve, ms));
// }

export async function activate(context: vscode.ExtensionContext) {

    // Select correct language server binary for this platform.
    let platform = `${os.platform()}-${os.arch()}`;
    if (platform === "darwin-arm64") {
        // Workaround while we cannot build native binary on github actions.
        platform = "darwin-x64";
    }
    const binPath = context.asAbsolutePath(`out/bin/dddlsp-${platform}`);
    if (!fs.existsSync(binPath)) {
        vscode.window.showErrorMessage(`wdyspec-support has no binary for platform "${platform}" and will not work. Contact the developer to fix this.`);
        return;
    }

    const serverOptions: ServerOptions = {
        command: binPath,
        transport: TransportKind.stdio,
    };
    const clientOptions: LanguageClientOptions = {
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
        const doc = "file://" + vscode.window.activeTextEditor?.document.uri.fsPath;
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
        // initialize puppeteer browser
        const browser = await puppeteer.launch({ headless: 'new' });
        const page = await browser.newPage();

        // set output path for the generated pdf
        const outputFolderPath = vscode.workspace.workspaceFolders !== undefined ?
            vscode.workspace.workspaceFolders[0].uri :
            vscode.Uri.file("");

        // start server to serve markdown images
        const app = express();
        const port = 3000;
        app.use(express.static(outputFolderPath.fsPath));
        const server = app.listen(port, function(){
            console.log(`Webserver listening at http://localhost:${port}`);
        });

        let isCreated = true;
        try {
            // get styled html file
            const html = String(await client.sendRequest("custom/ExportHTML", null));

            fs.writeFileSync(outputFolderPath.fsPath + '/index.html', html);

            await page.goto(`http://localhost:${port}/index.html`);

            // styling for header and footer templates
            const css = "<style>span { font-size:10px; margin: 0px 5px; }</style>";

            // print pdf with options
            const options: PDFOptions = {
                path: vscode.Uri.joinPath(outputFolderPath, "index.pdf").fsPath,
                margin: {
                    top: "50px",
                    bottom: "50px",
                },
                displayHeaderFooter: true,
                headerTemplate: `${css}<span class="date" style="margin-left: 5%;"></span>`,
                footerTemplate: `${css}<span style="display: flex; justify-content: flex-end; width: 95%;">Seite <span class="pageNumber"></span> von <span class="totalPages"></span></span>`,
                printBackground: true,
                timeout: 300000,
            }
            const ProgressOptions: vscode.ProgressOptions = {
                location: vscode.ProgressLocation.Notification
            }
            const task = (progress: vscode.Progress<{
                message?: string | undefined;
                increment?: number | undefined;
            }>, token: vscode.CancellationToken) => {
                progress.report({
                    message: "PDF wird generiert",
                });
                token.onCancellationRequested(async () => {
                    console.log("PDF printing cancelled");
                });
                return page.pdf(options);
            };
            await vscode.window.withProgress(ProgressOptions, task)
        } catch (error) {
            console.log(error);
            isCreated = false;
        } finally {
            await browser.close();
            server.close(function() {
                console.log("Stopping webserver.")
            });
        }

        isCreated ?
            vscode.window.showInformationMessage("wdy: PDF erfolgreich erstellt") :
            vscode.window.showErrorMessage("wdy: PDF konnte nicht erstellt werden");
    }));

    context.subscriptions.push(vscode.commands.registerCommand("ddd.GenerateGo", () => {

        client.sendRequest("custom/GenerateGo", null).then(() => {

        });
    }));

    //preview???


    client.onNotification("custom/newAsyncPreviewHtml", resp => {
        console.log("shall async preview")

        if (typeof resp === "string") {
            if (PreviewPanel.currentPanel != null) {
                if (resp==="lastPreviewParams missing"){
                    const tailwindUri = PreviewPanel.currentPanel?._tailwindUri;
                    const webviewPrefix =  PreviewPanel.currentPanel?._webviewPrefixUri;
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
        const doc = "file://" + vscode.window.activeTextEditor?.document.uri.fsPath;
        if (doc) {
            PreviewPanel.createOrShow(context.extensionUri, "<p>Einen Moment bitte...</p>");

            const tailwindUri = PreviewPanel.currentPanel?._tailwindUri;
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
            async deserializeWebviewPanel(webviewPanel: vscode.WebviewPanel, state: unknown) {
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

    console.log(extensionUri)

    return {
        // Enable javascript in the webview
        enableScripts: true,

        // And restrict the webview to only loading content from our extension's `media` directory.
        localResourceRoots: [vscode.Uri.joinPath(extensionUri, 'media'), localWorkspaceUri]
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
    public readonly _webviewPrefixUri: vscode.Uri;

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

        const tailwindUri = PreviewPanel.currentPanel?._tailwindUri;
        const webviewPrefix =  PreviewPanel.currentPanel?._webviewPrefixUri;
        client.sendRequest("custom/webViewParams", {WebviewWorkspacePrefix:webviewPrefix?.toString(),TailwindUri: tailwindUri?.toString()}).catch(e => console.log(e))

    }

    private constructor(panel: vscode.WebviewPanel, extensionUri: vscode.Uri, html: string) {
        this._panel = panel;
        this._extensionUri = extensionUri;
        this._html = html;

        this._tailwindUri = panel.webview.asWebviewUri(vscode.Uri.joinPath(this._extensionUri, 'media', 'tailwind.js'));
        if (vscode.workspace.workspaceFolders !== undefined) {
            this._webviewPrefixUri = panel.webview.asWebviewUri(vscode.workspace.workspaceFolders[0].uri);
        } else {
            this._webviewPrefixUri = vscode.Uri.file("");
        }
        
        

        // Set the webview's initial html content
        this._update();

        // Listen for when the panel is disposed
        // This happens when the user closes the panel or when the panel is closed programmatically
        this._panel.onDidDispose(() => this.dispose(), null, this._disposables);

        // Update the content based on view changes
        this._panel.onDidChangeViewState(
            () => {
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

