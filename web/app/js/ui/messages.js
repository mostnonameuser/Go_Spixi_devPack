import { processMarkdown } from '../utils/markdown.js';

export const handleSendClick = (app) => {
    if (app.isWaitingForResponse) {
        alert("CancelRequest Logic")
    } else {
        sendMessage(app);
    }
};

export const sendMessage = (app) => {
    const message = app.messageInput.value.trim();
    if (!message || !app.state.isReady || app.isWaitingForResponse) return;
    const messageId = window.SpixiTools?.getTimestamp?.();
    addMessage(app, message, "user", messageId);
    app.messageInput.value = "";
    setWaitingState(app, true);
    app.currentRequestId = messageId;
    SpixiAppSdk.sendNetworkProtocolData(
        app.protocolId,
        JSON.stringify({ action: "sendMessage",
            text: message,
            messageId: messageId,
            chatId: app.state.chatId,
            chat: app.state.chatName
        })
    );
};

export const scrollToBottom = (app) => {
    app.chatContainer.scrollTop = app.chatContainer.scrollHeight;
};

export const addMessage = (app, text, sender, messageId = null) => {
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

    const formattedText = processMarkdown(text);
    messageDiv.innerHTML = `
        <div>${formattedText}</div>
        <div class="message-time">${timeString}</div>
    `;
    app.chatContainer.appendChild(messageDiv);
    scrollToBottom(app);
    app.messageElements.set(messageId, messageDiv);
};

export const setWaitingState = (app, waiting) => {
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

