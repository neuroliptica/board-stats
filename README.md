# board-stats
Собирает статистику постинга по разным доскам сосача и строит прикольные графики в веб морде.

## server

```bash
$ cd server 
$ go build
$ ./stats
```

будет запущено на localhost:8080

запрос как `GET localhost:8080/api/stats?board={board}`

## client

```bash
$ cd client
$ npm start
```

