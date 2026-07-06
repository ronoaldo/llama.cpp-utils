# llama.cpp-utils

A lightweight wrapper, management, and automation toolkit for [`llama.cpp`](https://github.com/ggml-org/llama.cpp) designed for **Linux x86_64** systems.

This project is built around the intent to leverage **Vulkan inference** for hardware-accelerated LLM execution across a broad range of GPUs (AMD Radeon, NVIDIA GeForce/RTX, Intel Arc) without relying on vendor-specific SDKs (such as CUDA or ROCm). It is specifically configured and optimized to support GPUs with **12 or more GB of VRAM**, utilizing 8-bit KV caching, single-model VRAM management, and efficient layer offloading to run modern quantized models smoothly.

## Features

- **Broad GPU Compatibility via Vulkan**: Uses Vulkan compute shaders to provide hardware acceleration across AMD, NVIDIA, and Intel GPUs on Linux x86_64.
- **Optimized for 12GB+ VRAM**: Pre-configured defaults (such as 8-bit KV cache quantization `q8_0`, tuned batch sizes, and single-model concurrency limits) maximize model context and layer offloading within 12GB+ GPU VRAM budgets.
- **Automated Installation & Updates**: Automatically fetches and verifies the latest pre-compiled `vulkan-x64` release binaries of `llama.cpp`.
- **MCP & WebUI Integration**: Launches `llama serve` pre-configured with `--webui-mcp-proxy` for seamless integration with Model Context Protocol (MCP) clients and web interfaces.
- **VRAM & API Utilities**: Includes dedicated utilities to test model throughput (`llama-ping`) and unload inactive models from GPU VRAM (`llama-free`).

## Components

- [`llama-installer`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-installer): Automates the installation and updating of `llama.cpp` binaries using official pre-compiled `vulkan-x64` builds from GitHub.
- [`llama-launcher`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-launcher): Launches the `llama serve` daemon using model presets configured in `models.ini`, enforcing a single loaded model (`--models-max 1`) to preserve VRAM.
- [`models.ini`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/models.ini): The central configuration file for model presets, defining generation parameters, GPU layer offloading, Flash Attention, and memory parameters.
- [`llama-ping`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-ping): CLI utility to query the OpenAI-compatible `/v1/chat/completions` API and output response metrics (tokens/second).
- [`llama-free`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-free): Utility script to query active models and issue unload requests to free VRAM.

## Hardware & System Requirements

- **Operating System**: Linux x86_64 (AMD64)
- **GPU**: Any graphics card with Vulkan 1.2+ support (AMD Radeon, NVIDIA GeForce/RTX, Intel Arc)
- **VRAM**: **12 GB or more** of VRAM recommended
- **Software Dependencies**: `bash`, `curl`, `jq`, `tar`, and working Vulkan drivers (`vulkan-tools` / Mesa RADV / NVIDIA proprietary drivers)

## Quick Start

### 1. Installation

Install or update `llama.cpp` using the automated Vulkan installer:

```bash
./llama-installer
```

After installation, update your shell environment if `~/.local/bin` was newly added to your `PATH`:

```bash
source ~/.bashrc
```

### 2. Start the Server

Launch the `llama.cpp` server with Vulkan acceleration and MCP proxy support:

```bash
./llama-launcher
```

The server listens on `http://0.0.0.0:8084`.

### 3. Test a Model

Send a test prompt to an alias configured in `models.ini`:

```bash
./llama-ping gemma-4-26B-A4B "Explain quantum computing in simple terms"
```

### 4. Unload Models from VRAM

To free GPU memory without stopping the server daemon:

```bash
./llama-free
```

## Configuration

Model aliases, Hugging Face repositories, offload layers, and VRAM-saving options are configured in [`models.ini`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/models.ini).

Global defaults (`[*]`) set baseline VRAM-friendly values:
- `cache-type-k = q8_0` & `cache-type-v = q8_0` (quantized KV cache for memory savings)
- `ctx-size = 65536`
- `ubatch-size = 256` & `batch-size = 1024`

Each model has two presets: a default one that leaves VRAM and CPU headroom
for the desktop and streaming software (e.g. OBS), and a `-performance` one
that uses most of the GPU. See
[docs/fable-5-optimizations.md](docs/fable-5-optimizations.md) for the
measurements and tuning guide behind these values.

