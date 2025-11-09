import { initElements } from './ui/elements.js';
import { setupControls } from './ui/controls.js';
import * as messages from './ui/messages.js';
import { setupMessageListener } from './logic/network.js';

export class DemoApp {
    static DeBugMode = false;

    constructor(protocolId = "com.devapp", pingInterval = 5000) {
        this.DEBUG_MODE = DemoApp.DeBugMode;
        this.protocolId = protocolId;
        this.pingInterval = pingInterval;
        this.pingTimer = null;
        this.isWaitingForResponse = false;
        this.currentRequestId = 0;
        this.messageElements = new Map();
        this.state = {
            messages: [],
            isReady: false,
        };

        this.uiHandlers = {
            // messages
            onAddMessage: (msg, sender, id) => messages.addMessage(this, msg, sender, id),
            // window
            onSetWaitingState: (set) => messages.setWaitingState(this, set),
        };
    }

    onInit = (sessionId, userAddresses) => {
        initElements(this);
        setupControls(this);
        setupMessageListener(this);
        this.startPinging();
    };

    destroy = () => {
        this.stopPinging();
        document.removeEventListener("visibilitychange", () => handleVisibilityChange(this));
        document.removeEventListener("click", this.globalClickHandler);
    };


    // Ping logic
    pingQuIXI = () => {
        const timestamp = window.SpixiTools?.getTimestamp?.() ;
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
}
