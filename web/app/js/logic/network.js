// /js/logic/network.js

export const setupMessageListener = (app) => {
    SpixiAppSdk.onNetworkProtocolData = (senderAddress, receivedProtocolId, data) => {
        if (receivedProtocolId !== app.protocolId) return;
        try {
            const msg = JSON.parse(data);
            switch (msg.action) {
                case "system":
                    app.uiHandlers.onAddMessage(msg.text, "system", msg.messageId);
                    break;
                case "response":
                    if (app.DEBUG_MODE) {
                        // alert(`[DEBUG] Final: ` + JSON.stringify(msg));
                        console.log(`[DEBUG] Final: ` + JSON.stringify(msg));
                    }
                    if (msg.messageId === app.currentRequestId) {
                        app.uiHandlers.onAddMessage(
                            msg.message,
                            msg.sender === 'ai' ? 'ai' : 'user',
                            msg.messageId
                        )
                    }
                    app.uiHandlers.onSetWaitingState( false );
                    break;
            }
        } catch (err) {
            if (app.DEBUG_MODE) {
                alert(`[ERROR] Parsing message:
                ${err.message}
                Data: ${data}`);
            }
        }
    };
};