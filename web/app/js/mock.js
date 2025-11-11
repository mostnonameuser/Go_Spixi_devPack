// mock.js
const DEV_SERVER_WS_URL = 'ws://localhost:8888/ws';
(function () {
    if (location.protocol !== 'http:' ) {
        return;
    }
    var isSpixiEnvironment = false;
    var ws = null; 
    function logMock(message) {
        if (!isSpixiEnvironment) {
            if (typeof console !== 'undefined' && console.log) {
                console.log('[SpixiMock] ' + message);
            } else {
                alert('[SpixiMock] ' + message);
            }
        } else {
            // Can use ALERT here
        }
    }

    if (!isSpixiEnvironment) {
        logMock('SpixiSDK mock loaded in browser environment.');
        connectToDevServer();
        
        // --- SpixiTools mock ---
        window.SpixiTools = {
            version: 0.1,

            base64ToBytes: function (base64) {
                const binString = atob(base64);
                return Uint8Array.from(binString, (m) => m.codePointAt(0));
            },

            executeUiCommand: function (cmd) {
                try {
                    const decodedArgs = [];
                    for (let i = 1; i < arguments.length; i++) {
                        const bytes = SpixiTools.base64ToBytes(arguments[i]);
                        decodedArgs.push(new TextDecoder().decode(bytes));
                    }
                    if (typeof cmd === 'function') {
                        cmd.apply(null, decodedArgs);
                    } else {
                        logMock('executeUiCommand: cmd is not a function');
                    }
                } catch (e) {
                    const alertMessage = "Cmd: " + cmd + "\nArguments: " + decodedArgs.join(", ") + "\nError: " + e + "\nStack: " + e.stack;
                    alert(alertMessage);
                }
            },

            unescapeParameter: function (str) {
                return str.replace(/>/g, ">")
                    .replace(/</g, "<")
                    .replace(/&#92;/g, "\\")
                    .replace(/&#39;/g, "'")
                    .replace(/&#34;/g, "\"");
            },

            escapeParameter: function (str) {
                return str.replace(/&/g, "&amp;")
                    .replace(/</g, "<")
                    .replace(/>/g, ">")
                    .replace(/"/g, "&quot;")
                    .replace(/'/g, "&#039;");
            },

            getTimestamp: function () {
                return Math.round(Date.now());
            }
        };

        window.executeUiCommand = function (cmd) {
            SpixiTools.executeUiCommand.apply(null, arguments);
        };
        function connectToDevServer() {
            ws = new WebSocket(DEV_SERVER_WS_URL);

            ws.onopen = () => {
                console.log('[Mock] Connected to dev server');
            };

            ws.onmessage = (event) => {
                
                try {
                    const msg = JSON.parse(event.data);
                    if (typeof SpixiAppSdk.onNetworkProtocolData === 'function') {
                        SpixiAppSdk.onNetworkProtocolData("mock-peer", msg.protocolId, msg.data);
                    }
                } catch (e) {
                    console.error('[Mock] Invalid message from server:', event.data, e);
                }
            };

            ws.onerror = (err) => {
                console.error('[Mock] WebSocket error:', err);
            };

            ws.onclose = () => {
                console.log('[Mock] Dev server disconnected');
            };
        }

        // --- SpixiAppSdk mock ---
        const mockStorage = {};

        window.SpixiAppSdk = {
            version: 0.3,
            date: "2025-07-31",

            fireOnLoad: function () {
                logMock("fireOnLoad() called");
                // Эмулируем инициализацию через таймаут
                setTimeout(() => {
                    if (typeof SpixiAppSdk.onInit === 'function') {
                        SpixiAppSdk.onInit("mock-session-id", "mock-address");
                    }
                }, 100);
            },

            back: function () {
                logMock("back() called");
                history.back();
            },

            sendNetworkData: function (data) {
                logMock("sendNetworkData: " + data);
                // Эмулируем получение этого же сообщения как входящего
                setTimeout(() => {
                    if (typeof SpixiAppSdk.onNetworkData === 'function') {
                        SpixiAppSdk.onNetworkData("mock-peer", data);
                    }
                }, 50);
            },

            sendNetworkProtocolData: function (protocolId, data) {
                logMock("sendNetworkProtocolData: " + protocolId + " = " + data);
                console.log('[Mock] Sending to dev server:', protocolId, data);
                if (ws && ws.readyState === WebSocket.OPEN) {
                    ws.send(JSON.stringify({ protocolId, data }));
                } else {
                    console.warn('[Mock] WebSocket not ready, message dropped');
                }
            
            },

            getStorageData: function (key) {
                const value = mockStorage[key] || null;
                logMock("getStorageData: " + key + " => " + value);
                setTimeout(() => {
                    if (typeof SpixiAppSdk.onStorageData === 'function') {
                        SpixiAppSdk.onStorageData(key, value);
                    }
                }, 10);
            },

            setStorageData: function (key, value) {
                mockStorage[key] = value;
                logMock("setStorageData: " + key + " = " + value);
            },

            spixiAction: function (actionData) {
                logMock("spixiAction: " + actionData);
            },

            onInit: function(sessionId, userAddresses) {
                logMock("onInit handler should be overridden!");
            },
            onStorageData: function (key, value) {
                logMock("onStorageData(" + key + ", " + value + ")");
            },
            onNetworkData: function (senderAddress, data) {
                logMock("onNetworkData(" + senderAddress + ", " + data + ")");
            },
            onNetworkProtocolData: function (senderAddress, protocolId, data) {
                logMock("onNetworkProtocolData(" + senderAddress + ", " + protocolId + ", " + data + ")");
            },
            onRequestAccept: function (data) {
                logMock("onRequestAccept(" + data + ")");
            },
            onRequestReject: function (data) {
                logMock("onRequestReject(" + data + ")");
            },
            onAppEndSession: function (data) {
                logMock("onAppEndSession(" + data + ")");
            }
        };

        if (document.readyState === 'loading') {
            document.addEventListener('DOMContentLoaded', () => {
                logMock("DOM ready — triggering SpixiAppSdk.onInit");
                setTimeout(() => {
                    if (typeof SpixiAppSdk.onInit === 'function') {
                        SpixiAppSdk.onInit("mock-session-id", "mock-user-address");
                    } else {
                        alert("[MOCK ERROR] SpixiAppSdk.onInit is not a function!");
                    }
                }, 50); // небольшая задержка для гарантии
            });
        } else {
            SpixiAppSdk.fireOnLoad();
        }

    } else {
        logMock('Running inside Spixi environment – mock not applied.');
    }
})();
