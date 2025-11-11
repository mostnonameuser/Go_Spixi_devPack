// js/logic/network.js

const sendDebug = (app, lvl, text) => {
    if (!app.debmode){
        console.log("[!!!!] debugDisabled")
        return;
    }

    const levels = {
        0: "[DEBUG]",
        1: "[INFO]",
        2: "[WARN]",
        3: "[ERROR]",
        4: "[ERROR]"
    };

    const mode = levels[lvl] || "[UNKNOWN]";
    app.uiHandlers.onAddMessage(
        `${mode} ${text}`,
        "ai",
    )
    SpixiAppSdk.sendNetworkProtocolData(
        app.protocolId,
        JSON.stringify({
            action: "debug",
            text: `${mode} ${text}`,
            ts: window.SpixiTools?.getTimestamp?.()
        })
    );
};

const setupMessageListener = (app) => {
    SpixiAppSdk.onNetworkProtocolData = (senderAddress, receivedProtocolId, data) => {
        if (receivedProtocolId !== app.protocolId) return;
        try {
            const msg = JSON.parse(data);
            switch (msg.action) {
                case "debug":
                    if (!app.debmode){
                        app.uiHandlers.onAddMessage(
                            "debugMode on",
                            "ai",
                        )
                        app.netHandlers.onDebug(0, "Debug mode on");
                        app.netHandlers.turnOnDebug();
                    }
                    return;
                case "system":
                    app.uiHandlers.onAddMessage(msg.text, "system", msg.messageId);
                    app.netHandlers.onDebug(0, "SystenMessage Received");
                    return;
                case "response":
                    app.netHandlers.onDebug(0, "Response Received");
                    if (msg.messageId === app.currentRequestId) {
                        app.netHandlers.onDebug(0, `Response printing ${msg.message} from ${msg.sender} for ID ${msg.messageId}`);
                        app.uiHandlers.onAddMessage(
                            msg.message,
                            msg.sender === 'ai' ? 'ai' : 'user',
                            msg.messageId
                        )
                    }
                    app.netHandlers.onDebug(0, "Ready to send new message");
                    app.uiHandlers.onSetWaitingState( false );
                    return;
            }
        } catch (err) {
            app.netHandlers.onDebug(4, `[ERROR] Parsing message: ${err.message} Data: ${data}`);
        }
    };
};

window.network = {
    sendDebug: sendDebug,
    setupMessageListener: setupMessageListener,
};