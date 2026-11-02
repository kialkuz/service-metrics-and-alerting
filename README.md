# Service Metrics & Alerting
Service Metrics & Alerting — это система сбора и обработки метрик, реализованная на Go по архитектуре клиент–сервер.
## О проекте
Система предназначена для сбора, хранения и обработки метрик. Клиент собирает данные и отправляет их на сервер, где они проходят валидацию и сохраняются в базе данных PostgreSQL. Пользователи могут получать как отдельные метрики, так и их полный список.
## Основные возможности
### Клиент
    • Сбор системных метрик. 
    • Работа по протоколу HTTP. 
    • Гибкая настройка параметров отправки. 
### Сервер
    • Прием и обработка метрик. 
    • Валидация входных данных. 
    • Получение отдельной метрики. 
    • Получение списка всех метрик. 
    • Хранение данных в PostgreSQL. 
    • Масштабируемая и модульная архитектура.
## Архитектура
### Проект реализован с применением лучших практик backend-разработки:  
    • Client–Server Architecture — четкое разделение ответственности. 
    • Repository Pattern — абстракция слоя доступа к данным. 
    • Context-Based Operations — корректное управление временем выполнения и отменой запросов. 
    • REST API — удобное взаимодействие между компонентами. 
    • PostgreSQL — надежное хранение метрик. 
    • Валидация данных — обеспечение целостности информации.
## Установка и запуск
1. Клонирование репозитория
git clone https://github.com/kialkuz/service-metrics-and-alerting.git  
cd service-metrics-and-alerting
2. Настройка PostgreSQL
Создайте базу данных и укажите параметры подключения через переменные окружения:
export DATABASE_DSN="postgres://user:password@localhost:5432/metrics?sslmode=disable"
3. Запуск сервера
go run ./cmd/server
4. Запуск клиента
go run ./cmd/client

## API Примеры
Добавление метрики  
POST /update/counter/RandomValue/123.45  
Типы метрик: counter, gauge
Значение: дробное число  
Получение метрики  
GET /value/{type}/{name}  
Получение списка метрик  
GET /metrics

File: main
Build ID: 97ec911143d7ef02fdfd96832de54298216bfe43
Type: inuse_space
Time: 2026-10-23 22:08:26 MSK
Showing nodes accounting for 8736.48kB, 64.97% of 13447.64kB total
Dropped 1 node (cum <= 67.24kB)
      flat  flat%   sum%        cum   cum%
 9120.60kB 67.82% 67.82%  9120.60kB 67.82%  golang.org/x/net/webdav.(*memFile).Write
    1028kB  7.64% 75.47%     1028kB  7.64%  bufio.NewReaderSize (inline)
    1026kB  7.63% 83.10%     1026kB  7.63%  runtime.allocm
-1024.28kB  7.62% 75.48% -1024.28kB  7.62%  github.com/jackc/pgx/v5.(*Conn).getRows
 1024.09kB  7.62% 83.10%  1024.09kB  7.62%  context.(*cancelCtx).propagateCancel.func2
 -902.59kB  6.71% 76.38%  -902.59kB  6.71%  compress/flate.NewWriter (inline)
  525.43kB  3.91% 80.29%   525.43kB  3.91%  github.com/jackc/pgx/v5/internal/stmtcache.NewLRUCache (inline)
 -516.01kB  3.84% 76.45%  -516.01kB  3.84%  github.com/jackc/pgx/v5/internal/iobufpool.init.0.func1
 -515.19kB  3.83% 72.62%  -515.19kB  3.83%  github.com/go-playground/validator/v10.map.init.8
 -514.63kB  3.83% 68.80% -2053.59kB 15.27%  kialkuz/service-metrics-and-alerting/internal/server/handler.(*MetricsHandler).UpdatesHandler
 -514.63kB  3.83% 64.97% -1537.84kB 11.44%  kialkuz/service-metrics-and-alerting/internal/service/server.(*MetricsService).SaveMetricList
  514.38kB  3.83% 68.79%   514.38kB  3.83%  github.com/go-playground/validator/v10.map.init.9
    -514kB  3.82% 64.97%     -514kB  3.82%  bufio.NewWriterSize (inline)
    -513kB  3.81% 61.16%     -513kB  3.81%  github.com/ugorji/go/codec.helperEncDriverSimpleBytes.fastpathEList
  512.56kB  3.81% 64.97%  1024.57kB  7.62%  github.com/jackc/pgx/v5/pgconn.(*PgConn).receiveMessage
 -512.19kB  3.81% 61.16%  -512.19kB  3.81%  crypto/internal/impl.Register
 -512.07kB  3.81% 57.35%  -512.07kB  3.81%  net/url.parse
 -512.05kB  3.81% 53.54%  -512.05kB  3.81%  github.com/jackc/pgx/v5/pgconn/ctxwatch.(*ContextWatcher).Watch
  512.04kB  3.81% 57.35%   512.04kB  3.81%  golang.org/x/net/webdav.(*memFS).OpenFile
  512.01kB  3.81% 61.16%   512.01kB  3.81%  encoding/json.(*decodeState).literalStore
  512.01kB  3.81% 64.97%   512.01kB  3.81%  github.com/jackc/pgx/v5/pgproto3.(*AuthenticationSASL).Decode
         0     0% 64.97%     1028kB  7.64%  bufio.NewReader (inline)
         0     0% 64.97%  -902.59kB  6.71%  compress/gzip.(*Writer).Close
         0     0% 64.97%  -902.59kB  6.71%  compress/gzip.(*Writer).Write
         0     0% 64.97%  -512.19kB  3.81%  crypto/internal/fips140/bigmod.init.0
         0     0% 64.97%   512.01kB  3.81%  encoding/json.(*Decoder).Decode
         0     0% 64.97%   512.01kB  3.81%  encoding/json.(*decodeState).array
         0     0% 64.97%   512.01kB  3.81%  encoding/json.(*decodeState).object
         0     0% 64.97%   512.01kB  3.81%  encoding/json.(*decodeState).unmarshal
         0     0% 64.97%   512.01kB  3.81%  encoding/json.(*decodeState).value
         0     0% 64.97% -2956.17kB 21.98%  github.com/gin-gonic/gin.(*Context).Next
         0     0% 64.97%   512.01kB  3.81%  github.com/gin-gonic/gin.(*Context).ShouldBindJSON (inline)
         0     0% 64.97%   512.01kB  3.81%  github.com/gin-gonic/gin.(*Context).ShouldBindWith (inline)
         0     0% 64.97% -2956.17kB 21.98%  github.com/gin-gonic/gin.(*Engine).ServeHTTP
         0     0% 64.97% -2956.17kB 21.98%  github.com/gin-gonic/gin.(*Engine).handleHTTPRequest
         0     0% 64.97%   512.01kB  3.81%  github.com/gin-gonic/gin/binding.decodeJSON
         0     0% 64.97%   512.01kB  3.81%  github.com/gin-gonic/gin/binding.jsonBinding.Bind
         0     0% 64.97%  -512.05kB  3.81%  github.com/jackc/pgx/v5.(*Conn).Exec
         0     0% 64.97% -1024.28kB  7.62%  github.com/jackc/pgx/v5.(*Conn).Query
         0     0% 64.97%  -512.05kB  3.81%  github.com/jackc/pgx/v5.(*Conn).exec
         0     0% 64.97%  -512.05kB  3.81%  github.com/jackc/pgx/v5.(*Conn).execPrepared
         0     0% 64.97%  -512.05kB  3.81%  github.com/jackc/pgx/v5.(*dbTx).Exec
         0     0% 64.97%  1033.99kB  7.69%  github.com/jackc/pgx/v5.ConnectConfig
         0     0% 64.97%  1033.99kB  7.69%  github.com/jackc/pgx/v5.connect
         0     0% 64.97%  -516.01kB  3.84%  github.com/jackc/pgx/v5/internal/iobufpool.Get
         0     0% 64.97%  -512.05kB  3.81%  github.com/jackc/pgx/v5/pgconn.(*PgConn).ExecStatement
         0     0% 64.97%  -512.05kB  3.81%  github.com/jackc/pgx/v5/pgconn.(*PgConn).execExtendedPrefix
         0     0% 64.97%   512.01kB  3.81%  github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage
         0     0% 64.97%   508.56kB  3.78%  github.com/jackc/pgx/v5/pgconn.ConnectConfig
         0     0% 64.97%  -516.01kB  3.84%  github.com/jackc/pgx/v5/pgconn.ParseConfigWithOptions.func1
         0     0% 64.97%   508.56kB  3.78%  github.com/jackc/pgx/v5/pgconn.connectOne
         0     0% 64.97%   508.56kB  3.78%  github.com/jackc/pgx/v5/pgconn.connectPreferred
         0     0% 64.97%   512.01kB  3.81%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive
         0     0% 64.97%  -516.01kB  3.84%  github.com/jackc/pgx/v5/pgproto3.NewFrontend
         0     0% 64.97%  -516.01kB  3.84%  github.com/jackc/pgx/v5/pgproto3.newChunkReader (inline)
         0     0% 64.97% -1024.28kB  7.62%  github.com/jackc/pgx/v5/pgxpool.(*Conn).Query
         0     0% 64.97% -1024.28kB  7.62%  github.com/jackc/pgx/v5/pgxpool.(*Pool).Query
         0     0% 64.97%  -512.05kB  3.81%  github.com/jackc/pgx/v5/pgxpool.(*Tx).Exec
         0     0% 64.97%  1033.99kB  7.69%  github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func3
         0     0% 64.97%  1033.99kB  7.69%  github.com/jackc/puddle/v2.(*Pool[go.shape.*uint8]).initResourceValue.func1
         0     0% 64.97%  1859.42kB 13.83%  github.com/swaggo/files.init.10
         0     0% 64.97%  1197.01kB  8.90%  github.com/swaggo/files.init.11
         0     0% 64.97%  1588.17kB 11.81%  github.com/swaggo/files.init.12
         0     0% 64.97%   683.75kB  5.08%  github.com/swaggo/files.init.13
         0     0% 64.97%   794.19kB  5.91%  github.com/swaggo/files.init.14
         0     0% 64.97%  1197.01kB  8.90%  github.com/swaggo/files.init.7
         0     0% 64.97%  1595.23kB 11.86%  github.com/swaggo/files.init.8
         0     0% 64.97%   717.85kB  5.34%  github.com/swaggo/files.init.9
         0     0% 64.97%     -513kB  3.81%  github.com/ugorji/go/codec.init
         0     0% 64.97% -1024.28kB  7.62%  kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db.(*MemStorage).GetList
         0     0% 64.97%  -512.05kB  3.81%  kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db.(*MemStorage).SaveList
         0     0% 64.97%  -512.05kB  3.81%  kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db.(*MemStorage).updateListWithTransaction
         0     0% 64.97% -2956.17kB 21.98%  kialkuz/service-metrics-and-alerting/internal/middleware.WithComparing
         0     0% 64.97% -7786.75kB 57.90%  kialkuz/service-metrics-and-alerting/internal/middleware.WithLogging
         0     0% 64.97%  4830.58kB 35.92%  kialkuz/service-metrics-and-alerting/internal/router.Init.WithLogging.func2
         0     0% 64.97%  1025.13kB  7.62%  kialkuz/service-metrics-and-alerting/internal/router.Init.WithResponseSign.func4
         0     0% 64.97%  1025.13kB  7.62%  kialkuz/service-metrics-and-alerting/internal/router.Init.WithSign.func3
         0     0% 64.97%  -902.59kB  6.71%  kialkuz/service-metrics-and-alerting/internal/service/compress.(*gzipWriter).Close
         0     0% 64.97%  -511.16kB  3.80%  kialkuz/service-metrics-and-alerting/internal/service/server.(*MetricsService).GetGroupedByTypeAndName
         0     0% 64.97%  -513.12kB  3.82%  kialkuz/service-metrics-and-alerting/internal/service/server.(*MetricsService).GetList
         0     0% 64.97%  -512.07kB  3.81%  net/http.(*conn).readRequest
         0     0% 64.97% -2954.24kB 21.97%  net/http.(*conn).serve
         0     0% 64.97%     1028kB  7.64%  net/http.newBufioReader
         0     0% 64.97%     -514kB  3.82%  net/http.newBufioWriterSize
         0     0% 64.97%  -512.07kB  3.81%  net/http.readRequest
         0     0% 64.97% -2956.17kB 21.98%  net/http.serverHandler.ServeHTTP
         0     0% 64.97%  -512.07kB  3.81%  net/url.ParseRequestURI
         0     0% 64.97%  8606.63kB 64.00%  runtime.doInit (inline)
         0     0% 64.97%  8606.63kB 64.00%  runtime.doInit1
         0     0% 64.97%  8606.63kB 64.00%  runtime.main
         0     0% 64.97%      513kB  3.81%  runtime.mcall
         0     0% 64.97%      513kB  3.81%  runtime.mstart
         0     0% 64.97%      513kB  3.81%  runtime.mstart0
         0     0% 64.97%      513kB  3.81%  runtime.mstart1
         0     0% 64.97%     1026kB  7.63%  runtime.newm
         0     0% 64.97%      513kB  3.81%  runtime.park_m
         0     0% 64.97%     1026kB  7.63%  runtime.resetspinning
         0     0% 64.97%     1026kB  7.63%  runtime.schedule
         0     0% 64.97%     1026kB  7.63%  runtime.startm
         0     0% 64.97%     1026kB  7.63%  runtime.wakep
         0     0% 64.97%  -516.01kB  3.84%  sync.(*Pool).Get

Удалось снизить потребление памяти в:
1. middleware Logging -7786.75kB через инициализацию логгера в main.go
2. SaveMetricList через очистку времменных данных
