//main.js

function initDemoApp(sessionId = 'mock-session', userAddresses = 'mock-user') {
    if (window.__DemoApp) {
        window.__DemoApp.destroy();
    }

    const app = new DemoApp();
    app.onInit(sessionId, userAddresses);
    window.__DemoApp = app;

}

window.SpixiAppSdk.onInit = initDemoApp;

if (typeof SpixiAppSdk !== 'undefined') {

    if (location.protocol !== 'http:') {
        window.onload = () => SpixiAppSdk.fireOnLoad();
    }
}