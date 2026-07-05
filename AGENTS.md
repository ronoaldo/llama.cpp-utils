# Agent Instructions: llama.cpp-utils

## Core Workflow
- **Start Server**: Run `./llama-launcher` to start the `llama serve` instance on port `8080`.
- **Interact**: Use `./llama-ping <model_alias> [prompt]` to test a model.
- **Unload Models**: Run `./llama-free` to unload all loaded models from the server.

## Key Components
- `models.ini`: The source of truth for model aliases and generation settings (e.g., `temp`, `top-p`, `spec-type`).
- `llama-launcher`: Wrapper for `llama serve` with `--webui-mcp-proxy` enabled.
- `llama-ping`: A CLI utility to send chat completion requests to the server.

## Configuration & API
- **API Base URL**: `http://localhost:8080/v1`
- **Model Aliases**: Defined in `models.ini` under `[alias]`.
- **Dependencies**: Requires `llama` server, `curl`, and `jq`.
