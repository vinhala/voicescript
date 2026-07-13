#!/bin/sh
set -e

export PORT="${PORT:-80}"

exec /app/backend
