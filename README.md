# TCS stocks reader

Small service for update stock prices, using data from Tinkoff stock exchange.

This is **not**:

- an auto trading bot
- a big interactive client, which can help you to operate with your stocks nice and easy

This is:

- Small web service, which just get the prices and let them being parsed by other financial apps

## How to use

1. Make a copy of config file, fill all fields and place it in `config` folder near bin file
2. Start app in background mode. Also, you can use [systemd unit example](configs/sd_unit.service.example) for this step.
3. (opt.) Add config into you financial app

   As example (works in KMyMoney):

   | Param | Value |
   | :-- | :-- |
   | URL | http(s)://{{ IP:PORT }}/getTicker?name=%1 // or domain name instead of ip:port |
   | Identify by | Symbol |
   | Price | `price: (\d[\.\d]*)` |
   | Date | `date: (\d{2}\.\d{2}\.\d{4})` |
   | Date format | `%d.%m.%y` |

4. Call `/update` (you can use cron or timers for it)

## Handlers description

| Name | Description |
| :-- | :-- |
| `/update` | Request and store current prices of shares from your account |
| `/reload` | Update config from file without service restart (it doesn't affect on ip/port) |
| `/getCache` | Return processed data about prices |
| `/getPortfolio` | Return a clean response from TinkoffAPI |
| `/getTicker?name=TICKER` | Return price for `TICKER` and time, when this price has been got |
| `/getConfig` | Return current config |
| `/stop` | Graceful app stop |


