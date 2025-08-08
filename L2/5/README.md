Вывод:
error
В данном случае test() возвращает nil типа *customError, но при присваивании err = test(), интерфейс err становится (type: *customError, value: nil).
Однако, поскольку err имеет тип *customError, но значение nil, проверка if err != nil возвращает true, и выполняется ветка println("error").