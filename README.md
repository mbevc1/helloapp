# helloapp
HelloApp is a demo example

## Application info
Application is demonstration typical microservices workloads and is serving
HTTP traffic on port 8080. It's connecting to MySQL database and Redis cache
cluster.

### Environment parameters

| variable | default |
|-|-|
| `REDIS_HOST` | localhost |
| `REDIS_PORT` | 6379 |
|-|-|
|`DB_USER` | root |
|`DB_PASS` | password1 |
|`DB_HOST` | localhost |
|`DB_PORT` | 3306 |
|`DB_NAME` | mysql |

### Enpoints

Application listens on port 8080 and has some endpoints:
* `/` - main Hello content
* `/version` - version
* `/health` - current health status
