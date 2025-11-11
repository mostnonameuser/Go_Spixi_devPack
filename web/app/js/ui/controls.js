// js/ui/controls.js
window.setupControls = function (app) {

    app.sendButton?.addEventListener("click", () => {
        app.netHandlers.onDebug(0, "SendButton clicked")
        window.messages.handleSendClick(app);
        app.netHandlers.onDebug(0, "SendButton click fin")
    });

    app.messageInput?.addEventListener("keypress", (e) => {
        if (e.key === "Enter" && !e.shiftKey && !app.isWaitingForResponse) {
            e.preventDefault();
            app.netHandlers.onDebug(0, "SendButton Pressed")
            window.messages.sendMessage(app);
            app.netHandlers.onDebug(0, "SendButton Finished")
        }
    });
    app.backBtn?.addEventListener("click", () => window.SpixiAppSdk.back());
};