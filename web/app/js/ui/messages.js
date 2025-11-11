// js/ui/controls.js
function handleSendClick(app) {
    if (app.isWaitingForResponse) {
        alert("CancelRequest Logic")
    } else {
        sendMessage(app);
    }
};

function sendMessage(app) {
    const message = app.messageInput.value.trim();
    if (!message || !app.state.isReady || app.isWaitingForResponse) return;
    const messageId = window.SpixiTools?.getTimestamp?.();
    window.messages.addMessage(app, message, "user", messageId);
    app.messageInput.value = "";
    window.messages.setWaitingState(app, true);
    app.currentRequestId = messageId;
    SpixiAppSdk.sendNetworkProtocolData(
        app.protocolId,
        JSON.stringify({ action: "sendMessage",
            text: message,
            messageId: messageId,
        })
    );
};
function scrollToBottom(app) {
    app.chatContainer.scrollTop = app.chatContainer.scrollHeight;
};
function addMessage(app, text, sender, messageId = null) {
    if (!messageId) {
        messageId = window.SpixiTools?.getTimestamp?.() ;
    }
    const messageDate = new Date(messageId);
    const timeString = messageDate.toLocaleTimeString(
        'en-US',
        { hour: '2-digit', minute: '2-digit' }
    );
    const messageDiv = document.createElement("div");
    messageDiv.className = `message-bubble message-${sender}`;
    messageDiv.dataset.sender = sender;
    messageDiv.dataset.messageId = messageId;
    if (sender === 'system') {
        messageDiv.classList.add('message-system');
    }

    const formattedText = window.processMarkdown(text);
    messageDiv.innerHTML = `
        <div>${formattedText}</div>
        <div class="message-time">${timeString}</div>
    `;
    app.chatContainer.appendChild(messageDiv);
    scrollToBottom(app);
    app.messageElements.set(messageId, messageDiv);
};

function setWaitingState(app, waiting) {
    app.isWaitingForResponse = waiting;
    if (waiting) {
        app.inputFooter.classList.add("waiting");
        app.messageInput.disabled = true;
        app.sendButton.innerHTML = "❌";
        app.sendButton.title = "Отменить запрос";
    } else {
        app.inputFooter.classList.remove("waiting");
        app.messageInput.disabled = false;
        app.messageInput.focus();
        app.sendButton.innerHTML = "✉️";
        app.sendButton.title = "Отправить сообщение";
    }
};
window.messages = {
    addMessage: addMessage,
    setWaitingState: setWaitingState,
    sendMessage: sendMessage,
    handleSendClick: handleSendClick,
};
