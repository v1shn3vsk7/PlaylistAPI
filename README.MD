## Запуск сервера: 
из папки проекта:
`docker-compose up` (zsh)

## Запуск тестов:
### Для тестирования плейлиста
`docker build -t test-pkg -f internal/test/pkg/Dockerfile .`

`docker run --rm test-pkg go test ./internal/test/pkg -v`
### Для тестирования сервера
`docker build -t test-internal -f internal/test/internal/Dockerfile .`

`docker run --rm test-internal go test ./internal/test/internal -v`

## Запросы (gRPC):
```
Порт: 5536
```

Каждый метод возвращает Response. Responce представляет собой ответ от сервера
```
message Response {
  string result = 1;
} 
```
### Начать воспроизведение:
gRPC метод `Play`
```
rpc Play(google.protobuf.Empty) returns(Response){}
```

### Остановить воспроизведение:
gRPC метод `Pause`

```
rpc Pause(google.protobuf.Empty) returns(Response){}
```

### Следующая песня:
gRPC метод `Next`

```
rpc Next(google.protobuf.Empty) returns(Response){}
```

### Предыдущая песня:
gRPC метод `Prev`

```
rpc Prev(google.protobuf.Empty) returns(Response){}
```

### Добавить песню:
gRPC: `Add`

```
rpc Add(AddRequest) returns(Response) {}
```

Длительность песни нужно указывать в секундах
```
message AddRequest {
  string name     = 1;
  string artist   = 2;
  int32  duration = 3;
}
```

### Редактировать песню:
gRPC: `Edit`

```
rpc Edit(EditRequest) returns(Response) {}
```

```
message EditRequest {
  string prev_name    = 1;
  string prev_artist  = 2;
  string new_name     = 3;
  string new_artist   = 4;
  int32  new_duration = 5;
}
```

### Удалить песню:
gRPC: `Delete`

```
rpc Delete(DeleteRequest) returns(Response) {}
```

```
message DeleteRequest {
  string name   = 1;
  string artist = 2;
}
```

### Текущее состояние плейлиста:
gRPC: `Status`

```
rpc Status(google.protobuf.Empty)  returns(Response)  {}
```


