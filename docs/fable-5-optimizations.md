# Model Profile Optimization Experiment

Reference notes from a tuning session (2026-07-05) that produced the current
`models.ini` profiles. Kept for future re-tuning when models, quantizations,
or llama.cpp versions change.

## Environment

- **GPU**: AMD Radeon RX 6800 XT, 16 GB VRAM (Vulkan/RADV) — single GPU,
  shared with the desktop session (~1.4 GiB baseline) and OBS while streaming
- **CPU/RAM**: Ryzen 9 5900X (12 cores), 32 GB RAM
- **llama.cpp**: build `b9876-4b2a0cdee` (vulkan-x64)
- **Models**: `gemma-4-26B-A4B` UD-Q4_K_XL (~17 GB, MoE) and
  `Qwen3.6-35B-A3B-MTP` UD-Q4_K_XL (~22.8 GB, MoE + Gated Delta Net)

## Goal

Two profiles per model:

- **Default (live/OBS-safe)**: ≤ ~11 GiB total VRAM peak, leaving ~5 GiB and
  4 CPU cores free so streaming and the desktop stay smooth.
- **`-performance`**: best speed while keeping the desktop usable,
  ≤ ~14.8 GiB total VRAM peak (~1.6 GiB headroom).

## Methodology

Measurements used `llama-ping` (via `llama-launcher` + `models.ini` presets)
while sampling `/sys/class/drm/card*/device/mem_info_vram_used` every 0.5 s
and recording the peak. Prompt processing (PP) was measured with a ~6.6k-token
prompt; generation speed comes from the server `timings` field.

> Note: `llama-bench` was **not** used — numbers are single-run measurements
> through the HTTP server, good enough for VRAM validation and coarse speed
> comparison, but not statistically rigorous benchmarks.

Each profile change requires a server restart (`models.ini` is read at launch).

## Results (final configuration)

Peak is total device VRAM (desktop included). GPU total: 16363 MiB.

| Profile | Offload strategy | Peak VRAM | PP (6.6k tok) | Generation |
|---|---|---|---|---|
| `gemma-4-26B` | `ngl=99` + `n-cpu-moe=16`, ctx 32k | 11.1 GiB | 262 t/s | 31.7 t/s |
| `gemma-4-26B-performance` | `fit-target=2048` (auto) | 14.8 GiB | 326 t/s | 40.5 t/s |
| `qwen-3.6-35B` | `ngl=99` + `n-cpu-moe=32`, ctx 32k | 10.6 GiB | 135 t/s | 36.5 t/s |
| `qwen-3.6-35B-performance` | `ngl=99` + `n-cpu-moe=24` | 14.7 GiB | 170 t/s | 42.4 t/s |

### Rejected configurations

| Configuration | Peak VRAM | Generation | Why rejected |
|---|---|---|---|
| gemma perf, `n-cpu-moe=6` | 16.26 GiB (99%) | 32.3 t/s | VRAM overflow spilled to host memory: *slower* than the default profile and no headroom for OBS |
| qwen perf, `n-cpu-moe=20` | 16.25 GiB (99%) | 46.1 t/s | Fastest result, but ~100 MiB free is unsafe with a desktop/OBS on the same GPU |
| qwen perf, `fit-target=2048` | 13.7 GiB | 38.1 t/s | Safe but barely faster than the default profile (see `--fit` finding below) |
| global `ubatch-size=64` | — | PP 47 t/s (qwen) | Starved prompt processing; 256 gave ~3× PP at no measurable VRAM cost |

## Key findings

1. **MoE offload beats layer offload.** Both models are MoE with few active
   parameters. Keeping every layer on GPU (`n-gpu-layers = 99`) and pushing
   expert weights to CPU (`n-cpu-moe = N`) preserves GPU-side attention/KV and
   is far faster than lowering `n-gpu-layers`.

2. **Never fill the VRAM.** At ~99% usage the Vulkan driver spills buffers to
   host memory over PCIe; for Gemma this made the "max offload" configuration
   slower than the conservative one. Peaks around 14.7 GiB were the safe
   ceiling on this 16 GiB card with a desktop running.

3. **New `--fit` system (b9876, on by default).** Auto-adjusts *unset*
   arguments to fit free device memory, honoring a `--fit-target MiB` margin
   (default 1024). Explicitly setting `n-gpu-layers` disables it
   (`common_fit_params: ... already set by user, abort` in the log).
   - Great for Gemma: `fit-target = 2048` alone beat manual tuning.
   - Poor for Qwen 3.6: fit offloads whole layers instead of experts — the
     log warns `layer 0 is assigned to device CPU but the fused Gated Delta
     Net tensor is assigned to device Vulkan0` — so manual `n-cpu-moe` wins.

4. **`ubatch-size` matters for agent workloads.** Raising the global
   `ubatch-size` from 64 to 256 (batch 512 → 1024) tripled prompt-processing
   speed with no measurable VRAM increase. Long-context coding agents are PP
   bound, so this is the single cheapest win found.

5. **CPU threads as a live/perf knob.** `threads = 8` on default profiles
   keeps 4 cores free for OBS/desktop; `threads = 12` on performance profiles
   uses all physical cores for the CPU-side experts.

6. **Other new options in this llama.cpp version** (not yet explored):
   speculative decoding types `draft-eagle3`, `draft-dflash` and the
   `ngram-*` family; `--fit-ctx`; per-draft-model `n-cpu-moe`. The `ngram-*`
   drafts may help repetitive coding workloads.

## Re-tuning checklist

1. Start the server with `./llama-launcher` (restart after any `models.ini` change).
2. Load a profile with `./llama-ping <profile> "<long prompt>"`.
3. Watch VRAM: `cat /sys/class/drm/card*/device/mem_info_vram_used`.
4. Adjust `n-cpu-moe` in steps of ±2: lower = faster until the peak
   approaches ~14.7 GiB; back off at the first sign of overflow (slower
   generation is the symptom).
5. Unload between tests with `./llama-free`.
