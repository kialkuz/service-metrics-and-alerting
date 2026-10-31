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
Build ID: bf6cc6ee5d3749e0aedced5768742605dd99009a
Type: inuse_space
Time: 2026-10-23 22:08:26 MSK
Showing nodes accounting for 1966.77kB, 14.63% of 13447.64kB total
Dropped 1 node (cum <= 67.24kB)
      flat  flat%   sum%        cum   cum%
    3078kB 22.89% 22.89%     3078kB 22.89%  runtime.allocm
-1097.69kB  8.16% 14.73% -1097.69kB  8.16%  compress/flate.(*compressor).initDeflate (inline)
-1024.28kB  7.62%  7.11% -1024.28kB  7.62%  github.com/jackc/pgx/v5.(*Conn).getRows
 -525.43kB  3.91%  3.20%  -525.43kB  3.91%  github.com/jackc/pgx/v5/internal/stmtcache.NewLRUCache (inline)
 -516.01kB  3.84%  0.64%  -516.01kB  3.84%  github.com/jackc/pgx/v5/internal/iobufpool.init.0.func1
  514.63kB  3.83%  3.19% -1021.71kB  7.60%  kialkuz/service-metrics-and-alerting/internal/service/server.(*MetricsService).SaveMetricList
  513.69kB  3.82%  7.01%   513.69kB  3.82%  syscall.init.func1
  513.31kB  3.82% 10.83%  1008.46kB  7.50%  github.com/jackc/pgx/v5/pgxpool.NewWithConfig.func3
 -513.12kB  3.82%  7.01%  -513.12kB  3.82%  github.com/jackc/pgx/v5/pgxpool.(*connResource).getConn (inline)
    -513kB  3.81%  3.20%     -513kB  3.81%  github.com/ugorji/go/codec.helperEncDriverSimpleBytes.fastpathEList
  512.56kB  3.81%  7.01%  1024.57kB  7.62%  github.com/jackc/pgx/v5/pgconn.(*PgConn).receiveMessage
  512.22kB  3.81% 10.82%   512.22kB  3.81%  runtime.malg
 -512.19kB  3.81%  7.01%  -512.19kB  3.81%  crypto/internal/impl.Register
  512.17kB  3.81% 10.82%   512.17kB  3.81%  net/textproto.readMIMEHeader
 -512.07kB  3.81%  7.01%  -512.07kB  3.81%  net/url.parse
 -512.05kB  3.81%  3.20%  -512.05kB  3.81%  github.com/jackc/pgx/v5/pgconn/ctxwatch.(*ContextWatcher).Watch
  512.02kB  3.81%  7.01%   512.02kB  3.81%  syscall.anyToSockaddr
  512.01kB  3.81% 10.82%   512.01kB  3.81%  github.com/jackc/pgx/v5/pgproto3.(*AuthenticationSASL).Decode
     512kB  3.81% 14.63%      512kB  3.81%  bufio.NewWriterSize (inline)
         0     0% 14.63% -1097.69kB  8.16%  compress/flate.(*compressor).init
         0     0% 14.63% -1097.69kB  8.16%  compress/flate.NewWriter (inline)
         0     0% 14.63% -1097.69kB  8.16%  compress/gzip.(*Writer).Close
         0     0% 14.63% -1097.69kB  8.16%  compress/gzip.(*Writer).Write
         0     0% 14.63%  -512.19kB  3.81%  crypto/internal/fips140/bigmod.init.0
         0     0% 14.63% -2632.52kB 19.58%  github.com/gin-gonic/gin.(*Context).Next
         0     0% 14.63% -2632.52kB 19.58%  github.com/gin-gonic/gin.(*Engine).ServeHTTP
         0     0% 14.63% -2632.52kB 19.58%  github.com/gin-gonic/gin.(*Engine).handleHTTPRequest
         0     0% 14.63%  -512.05kB  3.81%  github.com/jackc/pgx/v5.(*Conn).Exec
         0     0% 14.63% -1024.28kB  7.62%  github.com/jackc/pgx/v5.(*Conn).Query
         0     0% 14.63%  -512.05kB  3.81%  github.com/jackc/pgx/v5.(*Conn).exec
         0     0% 14.63%  -512.05kB  3.81%  github.com/jackc/pgx/v5.(*Conn).execPrepared
         0     0% 14.63%  -512.05kB  3.81%  github.com/jackc/pgx/v5.(*dbTx).Exec
         0     0% 14.63%   495.15kB  3.68%  github.com/jackc/pgx/v5.ConnectConfig
         0     0% 14.63%   495.15kB  3.68%  github.com/jackc/pgx/v5.connect
         0     0% 14.63%  -516.01kB  3.84%  github.com/jackc/pgx/v5/internal/iobufpool.Get
         0     0% 14.63%  -512.05kB  3.81%  github.com/jackc/pgx/v5/pgconn.(*PgConn).ExecStatement
         0     0% 14.63%  -512.05kB  3.81%  github.com/jackc/pgx/v5/pgconn.(*PgConn).execExtendedPrefix
         0     0% 14.63%   512.01kB  3.81%  github.com/jackc/pgx/v5/pgconn.(*PgConn).peekMessage
         0     0% 14.63%  1020.58kB  7.59%  github.com/jackc/pgx/v5/pgconn.ConnectConfig
         0     0% 14.63%  -516.01kB  3.84%  github.com/jackc/pgx/v5/pgconn.ParseConfigWithOptions.func1
         0     0% 14.63%  1020.58kB  7.59%  github.com/jackc/pgx/v5/pgconn.connectOne
         0     0% 14.63%  1020.58kB  7.59%  github.com/jackc/pgx/v5/pgconn.connectPreferred
         0     0% 14.63%   512.01kB  3.81%  github.com/jackc/pgx/v5/pgproto3.(*Frontend).Receive
         0     0% 14.63%  -516.01kB  3.84%  github.com/jackc/pgx/v5/pgproto3.NewFrontend
         0     0% 14.63%  -516.01kB  3.84%  github.com/jackc/pgx/v5/pgproto3.newChunkReader (inline)
         0     0% 14.63% -1024.28kB  7.62%  github.com/jackc/pgx/v5/pgxpool.(*Conn).Query
         0     0% 14.63%  -513.12kB  3.82%  github.com/jackc/pgx/v5/pgxpool.(*Pool).Acquire
         0     0% 14.63% -1537.41kB 11.43%  github.com/jackc/pgx/v5/pgxpool.(*Pool).Query
         0     0% 14.63%  -512.05kB  3.81%  github.com/jackc/pgx/v5/pgxpool.(*Tx).Exec
         0     0% 14.63%  1008.46kB  7.50%  github.com/jackc/puddle/v2.(*Pool[go.shape.*uint8]).initResourceValue.func1
         0     0% 14.63%     -513kB  3.81%  github.com/ugorji/go/codec.init
         0     0% 14.63%   513.69kB  3.82%  golang.org/x/sys/cpu.init.0
         0     0% 14.63%   513.69kB  3.82%  golang.org/x/sys/cpu.processOptions
         0     0% 14.63% -1537.41kB 11.43%  kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db.(*MemStorage).GetList
         0     0% 14.63%  -512.05kB  3.81%  kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db.(*MemStorage).SaveList
         0     0% 14.63%  -512.05kB  3.81%  kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db.(*MemStorage).updateListWithTransaction
         0     0% 14.63% -2632.52kB 19.58%  kialkuz/service-metrics-and-alerting/internal/middleware.WithComparing
         0     0% 14.63% -7786.75kB 57.90%  kialkuz/service-metrics-and-alerting/internal/middleware.WithLogging
         0     0% 14.63%  5154.23kB 38.33%  kialkuz/service-metrics-and-alerting/internal/middleware.WithLogging.func1
         0     0% 14.63%  1543.89kB 11.48%  kialkuz/service-metrics-and-alerting/internal/router.Init.WithResponseSign.func3
         0     0% 14.63%  1543.89kB 11.48%  kialkuz/service-metrics-and-alerting/internal/router.Init.WithSign.func2
         0     0% 14.63% -1534.83kB 11.41%  kialkuz/service-metrics-and-alerting/internal/server/handler.(*MetricsHandler).UpdatesHandler
         0     0% 14.63% -1097.69kB  8.16%  kialkuz/service-metrics-and-alerting/internal/service/compress.(*gzipWriter).Close
         0     0% 14.63% -1024.28kB  7.62%  kialkuz/service-metrics-and-alerting/internal/service/server.(*MetricsService).GetGroupedByTypeAndName
         0     0% 14.63%  -513.12kB  3.82%  kialkuz/service-metrics-and-alerting/internal/service/server.(*MetricsService).GetList
         0     0% 14.63%   512.02kB  3.81%  net.(*Dialer).DialContext
         0     0% 14.63%   512.02kB  3.81%  net.(*netFD).connect
         0     0% 14.63%   512.02kB  3.81%  net.(*netFD).dial
         0     0% 14.63%   512.02kB  3.81%  net.(*sysDialer).dialParallel
         0     0% 14.63%   512.02kB  3.81%  net.(*sysDialer).dialSerial
         0     0% 14.63%   512.02kB  3.81%  net.(*sysDialer).dialSingle
         0     0% 14.63%   512.02kB  3.81%  net.(*sysDialer).dialTCP
         0     0% 14.63%   512.02kB  3.81%  net.(*sysDialer).doDialTCP (inline)
         0     0% 14.63%   512.02kB  3.81%  net.(*sysDialer).doDialTCPProto
         0     0% 14.63%   512.02kB  3.81%  net.internetSocket
         0     0% 14.63%   512.02kB  3.81%  net.socket
         0     0% 14.63%  1026.10kB  7.63%  net/http.(*conn).readRequest
         0     0% 14.63% -2120.42kB 15.77%  net/http.(*conn).serve
         0     0% 14.63%      512kB  3.81%  net/http.newBufioWriterSize
         0     0% 14.63% -2632.52kB 19.58%  net/http.serverHandler.ServeHTTP
         0     0% 14.63%   512.17kB  3.81%  net/textproto.(*Reader).ReadMIMEHeader (inline)
         0     0% 14.63%  -512.07kB  3.81%  net/url.ParseRequestURI
         0     0% 14.63%   513.69kB  3.82%  os.Getenv
         0     0% 14.63%      513kB  3.81%  runtime.(*scavengerState).wake
         0     0% 14.63%  -511.50kB  3.80%  runtime.doInit (inline)
         0     0% 14.63%  -511.50kB  3.80%  runtime.doInit1
         0     0% 14.63%      513kB  3.81%  runtime.gopreempt_m (inline)
         0     0% 14.63%      513kB  3.81%  runtime.goschedImpl
         0     0% 14.63%      513kB  3.81%  runtime.injectglist
         0     0% 14.63%      513kB  3.81%  runtime.injectglist.func1
         0     0% 14.63%  -511.50kB  3.80%  runtime.main
         0     0% 14.63%      513kB  3.81%  runtime.morestack
         0     0% 14.63%     2565kB 19.07%  runtime.mstart
         0     0% 14.63%     2565kB 19.07%  runtime.mstart0
         0     0% 14.63%     2565kB 19.07%  runtime.mstart1
         0     0% 14.63%     3078kB 22.89%  runtime.newm
         0     0% 14.63%   512.22kB  3.81%  runtime.newproc.func1
         0     0% 14.63%   512.22kB  3.81%  runtime.newproc1
         0     0% 14.63%      513kB  3.81%  runtime.newstack
         0     0% 14.63%     2565kB 19.07%  runtime.resetspinning
         0     0% 14.63%     2565kB 19.07%  runtime.schedule
         0     0% 14.63%     3078kB 22.89%  runtime.startm
         0     0% 14.63%      513kB  3.81%  runtime.sysmon
         0     0% 14.63%   512.22kB  3.81%  runtime.systemstack
         0     0% 14.63%     2565kB 19.07%  runtime.wakep
         0     0% 14.63%   513.69kB  3.82%  sync.(*Once).Do
         0     0% 14.63%   513.69kB  3.82%  sync.(*Once).doSlow
         0     0% 14.63%  -516.01kB  3.84%  sync.(*Pool).Get
         0     0% 14.63%   513.69kB  3.82%  syscall.Getenv
         0     0% 14.63%   512.02kB  3.81%  syscall.Getpeername
         0     0% 14.63%   513.69kB  3.82%  syscall.init.OnceFunc.func3
         0     0% 14.63%   513.69kB  3.82%  syscall.init.OnceFunc.func3.1
