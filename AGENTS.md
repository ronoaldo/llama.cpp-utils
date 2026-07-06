# Agent Instructions: llama.cpp-utils

## Overview & Architecture Intent
`llama.cpp-utils` is a lightweight management and automation toolkit for `llama.cpp` on **Linux x86_64**.
- **Inference Engine**: Uses Vulkan compute shaders for hardware acceleration across a broad range of GPUs (AMD Radeon, NVIDIA GeForce/RTX, Intel Arc).
- **Target Specs**: Optimized for GPUs with **12 GB or more VRAM**, leveraging 8-bit KV caching (`q8_0`), single-model concurrency (`--models-max 1`), and model offloading presets.

## Core Workflow
- **Install/Update**: Run [`./llama-installer`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-installer) to fetch and verify the latest pre-compiled `vulkan-x64` binaries of `llama.cpp`.
- **Start Server**: Run [`./llama-launcher`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-launcher) to start the `llama serve` instance on port `8084`.
- **Interact**: Use [`./llama-ping <model_alias> [prompt]`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-ping) to test a model.
- **Unload Models**: Run [`./llama-free`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-free) to unload all loaded models from VRAM.

## Key Components
- [`llama-installer`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-installer): Automates installation and updating using official pre-compiled `vulkan-x64` builds.
- [`llama-launcher`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-launcher): Wrapper for `llama serve` with `--webui-mcp-proxy` enabled and `--models-max 1` VRAM protection.
- [`models.ini`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/models.ini): Source of truth for model aliases, generation parameters, and VRAM settings (`q8_0` KV cache).
- [`llama-ping`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-ping): CLI utility to send OpenAI-compatible chat completion requests to the server.
- [`llama-free`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/llama-free): CLI utility to clear active models from GPU VRAM via `/models/unload`.

## Configuration & API
- **API Base URL**: `http://localhost:8084/v1`
- **Model Presets**: Defined in [`models.ini`](file:///media/ronoaldo/data/workspace/ai/llama-launcher/models.ini).
- **Dependencies**: Requires Linux x86_64, Vulkan drivers, `llama` server, `curl`, `jq`, and `tar`.

