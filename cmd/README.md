# tetrio-gateway

## Why?
Because I have multiple processes on the same machine using the tetrio api. I want to share the ratelimits.

## Goals
- Respect the rate limit of a single request every second.
- Inject some http cache headers.
- Fully compatible with the upstream API.