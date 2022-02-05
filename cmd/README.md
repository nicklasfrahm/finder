# Desktop App and CLI

This folder contains the directly consumable applications of `finder`, a desktop app and a CLI.

## Desktop App

A desktop application using WebViews and Svelte to provide a graphical user interface for file management.

### Development

During development, you can start the development server as shown below.

```bash
cd /web && npm run dev
```

Then, open a second terminal and run the following command.

```bash
make watch URL=http://localhost:3000 TARGET=bin/app-linux-amd64
```
