# natsuki - Database via HTTP
Designed for use with ROBLOX HTTP Requests via `NatsukiDriver.luau` Driver
Supported Databases:
- Postgres
- Redis

<br>

# Features
- 🏎 High Performance
  - Connection Pooling
  - Minimal HTTP Overhead
- 🔖 Logging
  - Store internal debug logs and http requests
  - Available in `kaho:entries` in redis database
  - Subscribe to log updates via channel `kaho:entries`
- 🛠 Easy Configuration
  - Use environment variables to quickly change settings
- 🧮 Always JSON Response
  - If response body contains an `error` field it has failed.
  - If response body contains a `results` field it has succeeded.

<br>

# Environment Variables
| Key                        |  Type  | Optional? |                 Default Value                  | Description                                                                                  |
| :------------------------- | :----: | :-------: | :--------------------------------------------: | :------------------------------------------------------------------------------------------- |
| NATSUKI_POSTGRES_URL       | string |    No     | `postgres://postgres:password@localhost:5432/` | Your Postgres Connection URL                                                                 |
| NATSUKI_POSTGRES_POOL_SIZE |  int   |    Yes    |                      `10`                      | Your Postgres Pool Size, the more connections the more concurrent queries can be run         |
| NATSUKI_REDIS_ADDR         | string |    No     |               `127.0.0.1:6379"`                | Your Redis Database Address                                                                  |
| NATSUKI_REDIS_PASS         | string |    Yes    |                     `nil`                      | Your Redis Database Password                                                                 |
| NATSUKI_REDIS_POOL_SIZE    |  int   |    Yes    |                      `10`                      | Your Redis Pool Size, the more connections the more concurrent commands can be run           |
| NATSUKI_PROXY              | string |    Yes    |                     `none`                     | Whether or not to enable proxy mode, allowed values are: cloudflare, none. Defaults to none. |
| NATSUKI_JWT                | string |    Yes    |             `your-256-bit-secret`              | Your JWT Secret, it has a default so ensure you change it in production                      |
| NATSUKI_KAHO_PRINT         |  any   |    Yes    |                     `true`                     | Whether or not to log to console. Disables if set to anything other than `nil`.               |
| NATSUKI_KAHO_ENABLE        |  any   |    Yes    |                     `true`                     | Whether or not to store log on Redis. Disables if set to anything other than `nil`.           |
| NATSUKI_ENABLE_SSL         |  any   |    Yes    |                     `nil`                      | Enables SSL if set to anything other than nil                                                |
| SSL_CERT_PATH              | string |    Yes    |                     `nil`                      | Path to your SSL Certificate File                                                            |
| SSL_KEY_PATH               | string |    Yes    |                     `nil`                      | Path to your SSL Key File                                                                    |

<br>

# Error Codes
| Code  | Description                                                |
| :---: | :--------------------------------------------------------- |
|   0   | Generic Error from Client (Bad Request, Bad Authorization) |
| 50011 | Generic Query/Command Error, see included `message` field  |
| 50012 | Generic Row Processing Error, see included `message` field |