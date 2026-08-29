# Agent Instructions: llama.cpp-utils

## Overview & Architecture Intent
`llama.cpp-utils` is a lightweight management and automation toolkit for `llama.cpp` on **Linux x86_64**.
- **Inference Engine**: Uses Vulkan compute shaders for hardware acceleration across a broad range of GPUs (AMD Radeon, NVIDIA GeForce/RTX, Intel Arc).
- **Target Specs**: Optimized for GPUs with **16 GB or more VRAM**, leveraging 8-bit KV caching (`q8_0`), single-model concurrency (`--models-max 1`), and model offloading presets.

## Core Workflow
- **Install/Update**: Run [`./llama-installer`](llama-installer) to fetch and verify the latest pre-compiled `vulkan-x64` binaries of `llama.cpp`.
- **Start Server**: Run [`./llama-launcher`](llama-launcher) to start the `llama serve` instance on port `1234`.
- **Interact**: Use [`./llama-ping <model_alias> [prompt]`](llama-ping) to test a model.
- **Unload Models**: Run [`./llama-free`](llama-free) to unload all loaded models from VRAM.

## Key Components
- [`llama-installer`](llama-installer): Automates installation and updating using official pre-compiled `vulkan-x64` builds.
- [`llama-launcher`](llama-launcher): Wrapper for `llama serve` with `--webui-mcp-proxy` enabled and `--models-max 1` VRAM protection.
- [`models.ini`](models.ini): Source of truth for model aliases, generation parameters, and VRAM settings (`q8_0` KV cache).
- [`llama-ping`](llama-ping): CLI utility to send OpenAI-compatible chat completion requests to the server.
- [`llama-free`](llama-free): CLI utility to clear active models from GPU VRAM via `/models/unload`.

## Configuration & API
- **API Base URL**: `http://localhost:1234/v1`
- **Model Presets**: Defined in [`models.ini`](models.ini).
- **Dependencies**: Requires Linux x86_64, Vulkan drivers, `llama` server, `curl`, `jq`, and `tar`.

