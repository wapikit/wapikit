<div align="center">
<br />
<p align="center">
<a href="https://wapijs.co"><img src="https://mintlify.s3-us-west-1.amazonaws.com/wapikit/logo/dark.svg" alt="@wapijs/wapi.js" height="100" /></a>
</p>
<br />
</div>

<p align="center">
<a href='https://join.slack.com/t/wapikit/shared_invite/zt-2kl7eg29s-4DfP9lFwojQg_yCcyW_w6Q'><img alt="Join Slack Community" src="https://img.shields.io/badge/slack%20community-join-green"/></a>
<a href='https://twitter.com/wapikit'><img alt="Follow WapiKit" src="https://img.shields.io/badge/%40wapikit-follow-blue"/></a>

<h4 align="center">
  <a href="https://join.slack.com/t/wapikit/shared_invite/zt-2kl7eg29s-4DfP9lFwojQg_yCcyW_w6Q">Slack</a> |
  <a href="https://docs.wapikit.com?ref=github">Docs</a> |
  <a href="https://wapikit.com?ref=github">Website</a>
</h4>
  
</p>

## The Open-source and AI Enabled WhatsApp Solution

![Dashboard view](https://res.cloudinary.com/dm4zlrwhs/image/upload/v1736856619/image_uospgn.png)

Watch the Demo Video [here](https://www.youtube.com/watch?v=wcUCGuGe2LY)

## 📖 About

WapiKit is an open-source, self-hosting enabled, single binary executable and performant WhatsApp campaign manager with team inbox & cross-platform integration availability.

## ✨ Major Features

- Contact List Management
- Campaign Broadcasting
- Multi Organization and Multi Agent Support
- Role Based Access Control
- Integration with software application via API
- Live Team Inbox
- Cross Platform Integration Suite

Check the [roadmap](#-roadmap) for upcoming features.

## Installation:

### Binary

- Download the [latest release](https://github.com/wapikit/wapikit/releases) and extract the binary.
- Make sure you have running instance of Postgres DB and Redis.
- `./wapikit --new-config` to generate config.toml with boilerplate configs. Add your configs by editing it.
- `./wapikit --install --idempotent` to setup the Postgres DB.
  -- You can use `--debug` flag to enable debug logs.
- Run `./wapikit` and visit `http://localhost:8000`

See [installation docs here](https://docs.wapikit.com/installation)

NOTE: WapiKit is right now available to self-hosting users only, our cloud version will be soon live [here](https://wapikit.com). You can join the wait-list, if want to get notified.

### Docker

COMING SOON...

## 📌 Status

Alpha Version - This application software is not stable right now. It is currently in public alpha. Report issues [here](https://github.com/wapikit/wapikit/issues).

## 📍 Roadmap:

- [x] Onboarding
- [x] Multi Organization Support with Member invite
- [x] Settings
- [x] Contact List Management with bulk contact import
- [x] RBAC
- [x] Campaign Manager
- [x] API Access Support
- [x] Analytics
- [ ] **In Progress: Live Team Inbox**
- [ ] Template Message Support Header Media, Copy Code button and other template message configuration while setting up a campaign
- [ ] Global AI Chat Assistant
- [ ] Feature flag System
- [ ] Cross Platform Integration Marketplace Infra
- [ ] Support OpenAI Integration (allowing users to connect AI with the application for auto reply configuration and chat summarization)
- [ ] No Code Chat Flow Configurator
- [ ] Configure rate limit response headers
- [ ] Support Embedded Sign-up as a Tech provider in the cloud SaaS keeping user owned API key in self-hosted version.
- [ ] Support Slack and Discord Integration (allowing users to connect slack workspace and discord server with application for notifications)
- [ ] Support E-commerce Product Catalog along with Order Management and Payments Support
- [ ] Support HubSpot Integration (Sync WhatsApp campaigns with HubSpot CRM to manage leads and automate follow-ups)
- [ ] Support Linear Integration (allowing users to create issues directly from the chat dashboard)
- [ ] Support Shopify Integration (Send order confirmations, shipping updates, and promotions via WhatsApp)
- [ ] Support WooCommerce Integration (Automate abandoned cart reminders and offer personalized discounts through WhatsApp)

We love to hear what do you want add in the list above. If you have got any idea / feature requests. Please reach out to us on our slack [here](https://join.slack.com/t/wapikit/shared_invite/zt-2kl7eg29s-4DfP9lFwojQg_yCcyW_w6Q)

### 🔗 Links:

- [Website](https://wapikit.com)
- [Documentation](https://docs.wapikit.com)
- [Wapi.go](https://go.wapikit.com): You can use this library if you want to build you own whatsapp Cloud API based chatbots.
- [Wapi.js](https://js.wapikit.com): You can use this javascript modules to build whatsapp chatbots in javascript.

## 🤝 Contribution Guidelines

Being an open-source project, we appreciate even the smallest contribution from your end. Please join our [slack channel](https://join.slack.com/t/wapikit/shared_invite/zt-2kl7eg29s-4DfP9lFwojQg_yCcyW_w6Q), to get involved.

For detailed guidelines, check [Contributing.md](./CONTRIBUTING.md).

## 📜 License

WapiKit is open-source and distributed under the AGPL 3.0 License. View [LICENSE](./LICENSE).

## 📞 Follow us

- [Slack Channel](https://join.slack.com/t/wapikit/shared_invite/zt-2kl7eg29s-4DfP9lFwojQg_yCcyW_w6Q)
- [Email](contact@wapikit.com)
- [Twitter](https://twitter.com/wapikit)
- [LinkedIn](https://www.linkedin.com/in/company/wapikit)
- [Github](https://github.com/wapikit)
