# gRPC-Cache
## Сервис кеширования данных 

### Сервис Cache:

    Cчитывает конфиг из файла config.yml (см ниже)
    Реализует GRPC метод GetRandomDataStream() без параметров, возвращающий поток из string
    Получив запрос  он NumberOfRequests(из конфига) раз в параллельных горутинах вызывает случайный URL из 
    набора URLs(из конфига)
    Полученные ответы по мере их поступления отдаются через поток и кэшируются в БД Redis 
    с временем жизни = случайным числом между MinTimeout(из конфига) и MaxTimeout(из конфига)
    При каждом подзапросе к URL необходимо сначала проверять наличие данных в БД Redis, 
    и если они там есть, то отдавать их оттуда, а не через запрос по URL - т.е. выполнять фунцию “кэша”
    При этом если в БД Redis записей нет или они просрочены, то сервис обращаеться напрямую по URL.
        
        ВАЖНО: Другие параллельно выполняющиеся горутины в данном процессе а также в других процессах / на других 
        серверах также не должны лезть в URL - они должны дожидаться окончания выполнения запроса первой 
        горутиной и получить данные от нее или через Redis. 
        Т.е. одновременно никогда не должно быть повторных обращений по одинаковым URL


### Сервис Consumer:

    Создает 1000 горутин и в каждой из них делает запросы ко второму сервису.