# Go Spixi DevPack

> **The first development toolkit for Ixian Spixi MiniApps**  
> Build, test, and debug your MiniApp **in the browser** — then deploy to Spixi **without changing a single line of code**.

Created by **[mostnonameuser](https://github.com/mostnonameuser)**.

---

## Features

- **Dual-mode runtime**: works in **browser** (for dev) and **Spixi** (for prod)
- **Go backend**: handles Ixian messages [QuIXI+MQTT] + WebSocket API
- **Live browser testing**: no Spixi rebuild needed for UI tweaks
- **Modular frontend**: clean separation of UI, logic, and styles
- **MIT Licensed**: fully compatible with [Ixian Core](https://github.com/ProjectIxian/Ixian-Core)

---

## Quick Start

### 1. Clone & prepare
```bash
git clone https://github.com/mostnonameuser/Go_Spixi_devPack.git
cd Go_Spixi_devPack
```
### 2. Place your MiniApp frontend
```
web/app/
├── index.html
├── css/
├── js/
└── .../
```
Your app will use the same SDK calls (SpixiAppSdk.sendNetworkProtocolData, etc.) in both browser and Spixi.

### 3. Run the dev server
```bash
cd cmd/quixi 
go run .
```
### 4. Develop in browser → deploy to Spixi
   * Test UI/UX, streaming, chat logic in browser
   * When ready, zip web/app/ and install as Spixi MiniApp
   * Zero code changes required

### Project Structure
```
Go_Spixi_devPack/
├── cmd/quixi/            # Go entry point
├── internal/
│   ├── config/           # Configuration
│   ├── devserver/        # DevWeb  
│   └── nwtwork/          # Ixian message processing / WS API
├── web/app/              # ← YOUR MINIAPP GOES HERE
│   ├── index.html
│   └── js/
└── go.mod
```
### Disclaimer

This project is an unofficial development toolkit for Ixian Spixi MiniApps.
It is not affiliated with, maintained, or endorsed by Ixian Ltd.
Ixian, Spixi, and related trademarks are property of their respective owners.

This software does not include any Ixian proprietary code.

### License

Distributed under the MIT License. See LICENSE for details.

Compatible with Ixian Core (also MIT Licensed).

### "The best way to predict the future is to build it." — Go Spixi DevPack
