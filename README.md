# Draft-API-Framework

draft go api framework

- go serve http2
- middlewares
    - FrontMiddleware: apply base request requirement and log statistics
    - AuthMiddleware: token base authentication and rate limit and log statistics
    - CacheMiddleware: bolt db
    - IsAdminMiddleware
- signin/signup user handlers
- JSON validations and other Utils
- example use for response flatebuffers
- tests and benchmarks
- ...
