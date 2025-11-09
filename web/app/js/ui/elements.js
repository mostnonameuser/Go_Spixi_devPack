export const initElements = (app) => {
    if (location.protocol === 'http:') {
        app.state.isReady = true;
    }
    app.chatContainer = document.getElementById("chatContainer");
    app.messageInput = document.getElementById("messageInput");
    app.sendButton = document.getElementById("sendButton");
    app.backBtn = document.getElementById("backBtn");
    app.inputFooter = document.getElementById("inputFooter");
};