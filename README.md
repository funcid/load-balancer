# Балансировщики нагрузки

## Абстрактно

Данный каталог содержит различные алгоритмы балансировки нагрузки, реализованные на Go. Каждый алгоритм предназначен для эффективного распределения входящего трафика на основе различных стратегий с целью оптимизации производительности, использования ресурсов и отзывчивости пользователей.

## Сравнение

| Реализация                                                              | Описание                                                        | Когда использовать                                                                  |
|-------------------------------------------------------------------------|-----------------------------------------------------------------|-------------------------------------------------------------------------------------|
| <a href="./internal/algorithms/geographical.go">geographical.go</a>     | Направляет трафик на основе местоположения                      | Когда нужно минимизировать задержку, обслуживая пользователей с ближайших серверов  |
| <a href="./internal/algorithms/geographical.go">iphash.go</a>           | Использует хэш IP клиента для выбора сервера                    | Когда необходимо сохранить сессию на основе IP клиента                              |
| <a href="./internal/algorithms/geographical.go">leastactive.go</a>      | Выбирает сервер с наименьшим количеством активных соединений    | Когда нужно балансировать на основе текущей загрузки сервера                        |
| <a href="./internal/algorithms/geographical.go">leastconnections.go</a> | Направляет трафик на сервер с наименьшим количеством соединений | Когда нужно оптимизировать распределение соединений между серверами                 |
| <a href="./internal/algorithms/geographical.go">random.go</a>           | Случайно распределяет трафик                                    | Когда нужен простой и равномерный алгоритм без сложной логики                       |
| <a href="./internal/algorithms/geographical.go">roundrobin.go</a>       | Поочередно направляет трафик на сервера                         | Когда нужно обеспечить равное распределение трафика между серверами последовательно |