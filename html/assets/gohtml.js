let lastHotReloadToken = null;

/**
 * hotReload long polls the given url and causes a refresh if required.
 * @param {string} url
 */
function hotReload(url) {
    const myurl = url;
    let curry = function () {
        // console.trace(); is this a bug? looks like accumulating stack infinietly
        fetch(myurl).then(function (response) {
            return response.text()
        }).then(function (str) {
            if (lastHotReloadToken !== str) {
                console.log("last token:", lastHotReloadToken, "new token", str);
                lastHotReloadToken = str;
                refresh();
            }
            setTimeout(curry, 1000); // try reconnect with a delay
            console.log("retry", myurl)
        }).catch(function (err) {
            awaitService(myurl).then(()=>{
                refresh();
                curry();
            });
        })
    }


    curry();
}

/**
 *
 * @param {string} url
 * @returns {Promise<boolean>}
 */
function awaitService(url) {
    return fetch(url + "?ping").then(r => r.status === 200).catch(e => {
        return wait(500).then(()=>awaitService(url))
    })
}

function wait(delay){
    return new Promise((resolve) => setTimeout(resolve, delay));
}


/**
 * Refresh sends the entire state without an event, to issue a re-rendering.
 * This is especially interesting for hot-reloading, but can also be used to implement an
 * automatic update mechanism (e.g. through polling or web sockets).
 */
function refresh() {
    send("!refresh", null)
}

/**
 * send submits the message to the server, which responds either with a new page or an artificial redirect response.
 * @param {string} msgType
 * @param {object|HTMLFormElement} msgData
 */
function send(msgType, msgData) {
    // We keep the page state from which this page has been rendered within a header>meta tag.
    // Sessions are shared between browser tabs and are therefore not suited.
    // Localstorage may work, but when cleared this may easily cause any sort of invalid states and transitions.
    // It is obvious, that this is never be a security issue, but depending on browser behavior, this may
    // lead to massive UX confusions.
    let state = document.getElementById("_state").content;
    console.log(state);

    sendWithState(state, msgType, msgData)
}


/**
 * send submits the message to the server, which responds either with a new page or an artificial redirect response.
 *     - headers are automatically replayed by the browser on redirects, which does not make any sense
 *     for our message based transition model across pages
 *     - also it may be dangerous because proxies or gateways may throw our headers away
 *     evtType := request.Header.Get("X-WDY-EventType")
 *
 *     this approach sucks really bad, but the fetch api is entirely broken for 3xx codes
 * @param {object} state
 * @param {string} msgType
 * @param {object|HTMLFormElement} msgData
 */
function sendWithState(state, msgType, msgData) {
    if (msgType === "") {
        console.log("no msgtype => trigger GET update without state and data", window.location)
        handleFetch(fetch(window.location, {
            method: "GET",
        }))

        return
    }

    let data;

    if (msgData instanceof HTMLFormElement) {
        // if called inline, we have a global "event", which must be prevented for forms, otherwise it will
        // be sent twice and automatic one is invalid, due to missing state and msg information
        if (event != null) {
            event.preventDefault();
        }

        data = new FormData(msgData);
        console.log("sending form", data)
    } else {
        data = new FormData();
        data.set("_eventData", JSON.stringify(msgData));
    }

    data.set("_state", state);
    data.set("_eventType", msgType);

    console.log("sending POST to", window.location)
    handleFetch(fetch(window.location, {
        method: "POST",
        body: data,
    }))
}

/**
 *
 * @param {Promise<Response>} promise
 */
function handleFetch(promise) {
    promise.then(async function (response) {
        let content = await response.text()
        let contentType = response.headers.get("Content-Type")
        return {contentType: contentType, content: content, url: response.url}
    }).then(function (msgRes) {
        console.log("2");
        if (msgRes.contentType === "application/json") {
            console.log("3")
            // this is the virtual redirect response
            let obj = JSON.parse(msgRes.content);
            let target = obj["target"]; // must not be null
            let navDir = obj["navDir"];
            let state = obj["state"]; // may be null
            let msgType = obj["msgType"]; // may be null
            let msgData = obj["msgData"]; // may be null

            // put the new target onto the browser navigation stack
            switch (navDir) {
                //TODO iterate backwards to clean up the state
                default:
                    history.pushState({url: target}, "", target);
            }

            sendWithState(state, msgType, msgData);
            return null
        }

        console.log("4")
        let doc = new DOMParser().parseFromString(msgRes.content, "text/html")
        Idiomorph.morph(document.documentElement, doc.documentElement, {head: {style: 'morph'}})


    }).catch(function (err) {
        console.warn('failed to fetch', err);
        showFatalError("failed processing state message", err);
    });
}

/**
 *
 * @param {string} caption
 * @param {string} text
 */
function showFatalError(caption, text) {
    const fatalErrTpl = `<div role="alert">
  <div id="ErrCaption" class="bg-red-500 text-white font-bold rounded-t px-4 py-2">
    Danger
  </div>
  <div class="border border-t-0 border-red-400 rounded-b bg-red-100 px-4 py-3 text-red-700">
    <p id="ErrText">Something not ideal might be happening.</p>
  </div>
</div>`

    document.body.innerHTML = fatalErrTpl
    document.getElementById("ErrCaption").innerText = caption
    document.getElementById("ErrText").innerText = text
}


history.pushState({url: location.href}, "", location.href);
console.log("pushed initial history", location.href);


window.addEventListener("popstate", (event) => {
    console.log("pop state");
    if (event.state) {
        fetch(event.state.url,).then(function (response) {

            return response.text();
        }).then(function (html) {
            let doc = new DOMParser().parseFromString(html, "text/html")
            Idiomorph.morph(document.documentElement, doc.documentElement, {head: {style: 'morph'}})

        }).catch(function (err) {
            console.warn('failed to fetch', err);
            showFatalError("failed processing state message", err);
        });
    }
})

console.log("hg (htmlgo) is here");