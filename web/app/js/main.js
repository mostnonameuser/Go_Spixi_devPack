//main.js
import { DemoApp } from './DemoApp.js';

console.log("DemoApp class defined");

function initDemoApp(sessionId = 'mock-session', userAddresses = 'mock-user') {
    if (window.__DemoApp) {
        window.__DemoApp.destroy();
    }
    const app = new DemoApp();
    app.onInit(sessionId, userAddresses);
    window.__DemoApp = app;
}

window.SpixiAppSdk.onInit = initDemoApp;

console.log("SpixiAppSdk.onInit assigned:", window.SpixiAppSdk.onInit);

if (typeof SpixiAppSdk !== 'undefined') {
    if (location.protocol !== 'http:') {
        window.onload = () => SpixiAppSdk.fireOnLoad();
    }
}