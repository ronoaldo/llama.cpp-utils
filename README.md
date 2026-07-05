# llama.cpp-utils

A lightweight wrapper and automation tool for `llama.cpp` with Vulkan support on Ubuntu AMD64.

## Components

- `llama-installer`: Automates the installation and update of `llama.cpp` using pre-compiled Vulkan binaries.
- `llama-launcher`: A wrapper for `llama serve` with `--webui-mcp-proxy` enabled.
- `llama-ping`: A CLI utility to test models via chat completion requests.
- `llama-free`: Unloads all loaded models from the server.

## Installation

To install or update `llama.cpp` with Vulkan support, run:

```bash
./llama-installer
```

### Requirements

- Ubuntu AMD64 (x86_64)
- `curl`, `jq`, `tar`
- Vulkan drivers installed on your system

### Post-Installation

After running the installer, you may need to update your current shell session's `PATH`:

```bash
source ~/.bashrc
```

The `llama` command will be available in your terminal, pointing to the latest installed version of `llama-cli`.

## Usage

### Starting the Server

```bash
./llama-launcher
```

### Testing a Model

```bash
./llama-ping <model_alias> "[prompt]"
```

### Unloading Models

```bash
./llama-free
```

## Configuration

Model aliases and generation settings (e.g., `temp`, `top-p`) are managed in `models.ini`.
