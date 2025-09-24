# Contyard

A tui application that monitors and manages Docker/Podman containers on a local  
or remote host. It displays real-time metrics (CPU, memory, network usage) for  
running containers, it allows you stopping/starting containers, and it shows logs  
interactively.

**Status**: Alpha (v0.0.1). Expect bugs. Report issues at https:://github.com/MikelGV/Contyard/issues.

## Features
- Real-time TUI with CPU/memory stats.
- Suport Docker(`-d`), Podman (`-p`), Kubernetes (`-k`, experimental).
- Commands: `contyard`, `contyard -d`, `contyard -p -a`.

## Prerequisites
- Docker: for `-d`.
- Podman: for `-p`.
- Kubectl: for `-k`.

## Installation
## Maual
1. Download from [GitHub Releases](https://github.com/MikelGV/Contyard/releases/tag/v0.0.1-beta):
    - Linux: `contyard-linux-amd64-0.0.1.tar.gz`
    - macOS: `contyard-darwin-amd64-0.0.1.tar.gz`
    - Windows: `contyard-windows-amd64-0.0.1.zip`
2. Extract:
    ```bash
    tar -xzf contyard-linux-amd64-0.0.1.tar.gz 
3. Move to PATH:
    mv contyard-linux-amd64-0.0.1 /usr/local/contyard
    chmod +x /usr/local/bin/contyard
4. Verify: contyard --version (outputs 0.0.1-beta

## One-Line Install
curl -L https://raw.githubusercontent.com/MikelGV/Contyard/v.0.0.1-beta/install.sh | bash

## Usage
contyard -d # Docker stats
contyard -p # Podman stats 
contyard -k # Kubernetes stats 
contyard -p -a # Podman + Kubernetes 
contyard -v # Show version 

## License
[MIT](https://github.com/MikelGV/Contyard/blob/main/LICENSE)
```
