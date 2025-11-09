// /js/ui/controls.js

import { handleSendClick, sendMessage } from './messages.js';

export const setupControls = (app) => {

    app.sendButton?.addEventListener("click", () => {
        if (app.DEBUG_MODE) {
            alert("Send button pressed");
        }
        handleSendClick(app);
    });

    app.messageInput?.addEventListener("keypress", (e) => {
        if (e.key === "Enter" && !e.shiftKey && !app.isWaitingForResponse) {
            e.preventDefault();
            sendMessage(app);
        }
    });

    app.backBtn?.addEventListener("click", () => SpixiAppSdk.back());
};
