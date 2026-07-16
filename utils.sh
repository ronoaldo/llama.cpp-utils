#!/usr/bin/env bash
# utils.sh - Funções auxiliares compartilhadas para os scripts do llama.cpp-utils

# ==========================================
# Definição de Cores (ANSI)
# ==========================================
readonly COLOR_RESET='\033[0m'
readonly COLOR_RED='\033[0;31m'
readonly COLOR_GREEN='\033[0;32m'
readonly COLOR_YELLOW='\033[1;33m'
readonly COLOR_BLUE='\033[0;34m'
readonly COLOR_CYAN='\033[0;36m'
readonly COLOR_BOLD='\033[1m'

# ==========================================
# Funções de Logging
# ==========================================

# log_info: Mensagens informativas normais (Verde)
log_info() {
    echo -e "${COLOR_GREEN}[INFO]${COLOR_RESET} $1"
}

# log_warn: Alertas (Amarelo)
log_warn() {
    echo -e "${COLOR_YELLOW}[WARN]${COLOR_RESET} $1"
}

# log_err: Erros críticos, impresso na saída de erro padrão stderr (Vermelho)
log_err() {
    echo -e "${COLOR_RED}[ERRO]${COLOR_RESET} $1" >&2
}

# log_step: Destaca etapas importantes da execução (Ciano Negrito)
log_step() {
    echo -e "${COLOR_CYAN}${COLOR_BOLD}* $1${COLOR_RESET}"
}

# log_success: Destaca um sucesso conclusivo
log_success() {
    echo -e "${COLOR_GREEN}${COLOR_BOLD}✔ $1${COLOR_RESET}"
}

# ==========================================
# Outros Helpers Úteis
# ==========================================

# die: Imprime erro e encerra o script com código de falha
die() {
    log_err "$1"
    exit 1
}

# check_dep: Garante que um comando necessário está instalado
check_dep() {
    local cmd=$1
    if ! command -v "$cmd" &> /dev/null; then
        die "A dependência obrigatória '$cmd' não foi encontrada. Instale-a e tente novamente."
    fi
}
