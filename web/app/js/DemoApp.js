// js/DemoApp.js

window.DemoApp = class DemoApp {
    constructor(protocolId = "com.devapp", pingInterval = 30000) {
        this.debmode = false;
        this.protocolId = protocolId;
        this.pingInterval = pingInterval;
        this.pingTimer = null;
        this.isWaitingForResponse = false;
        this.currentRequestId = 0;
        this.messageElements = new Map();
        this.state = {
            messages: [],
            isReady: true,
        };

        this.uiHandlers = {
            // messages
            onAddMessage: (msg, sender, id) => window.messages.addMessage(this, msg, sender, id),
            // window
            onSetWaitingState: (set) => window.messages.setWaitingState(this, set),
        };
        this.netHandlers = {
            onDebug: (lvl, text) => window.network.sendDebug(this, lvl, text),
            turnOnDebug: () => this.debmode = true,
        }
    }

    onInit = (sessionId, userAddresses) => {
        window.initElements(this);
        window.setupControls(this);
        window.network.setupMessageListener(this);
        this.startPinging();
    };

    destroy = () => {
        this.stopPinging();
    };

    pingQuIXI = () => {
        const timestamp = window.SpixiTools?.getTimestamp?.();
        window.SpixiAppSdk?.sendNetworkProtocolData(
            this.protocolId,
            JSON.stringify({ action: "ping", ts: timestamp })
        );
    };

    startPinging = () => {
        if (!this.pingTimer) {
            this.pingQuIXI();
            this.pingTimer = setInterval(() => this.pingQuIXI(), this.pingInterval);
        }
    };

    stopPinging = () => {
        if (this.pingTimer) {
            clearInterval(this.pingTimer);
            this.pingTimer = null;
        }
    };
};

