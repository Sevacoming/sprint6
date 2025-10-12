#!/usr/bin/env bash
set -e
PORT="${PORT:-8083}"
MSYS2_ARG_CONV_EXCL='*' /c/Windows/System32/taskkill.exe /IM morse.exe /F 2>/dev/null || true
PID=$(/c/Windows/System32/netstat.exe -ano | awk -v p=":$PORT" '$0 ~ p && $6=="LISTENING"{pid=$5} END{print pid}')
[ -n "$PID" ] && MSYS2_ARG_CONV_EXCL='*' /c/Windows/System32/taskkill.exe /F /PID "$PID" || true
echo "Открывай: http://127.0.0.1:$PORT"
PORT="$PORT" ./morse.exe
