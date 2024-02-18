>Это - проект, сделанный за 14 дней для курса по Golang от Яндекса

### Запуск 🤔:
###### Front-end ✨:
Я запускаю используя расширение для VS Code - Live ( https://marketplace.visualstudio.com/items?itemName=ritwickdey.LiveServer ) и не парюсь. По умолчанию сервер запускается по адресу http://127.0.0.1:5500/  (дальше путь к файлу HTML), в моем случае это http://127.0.0.1:5500/frontend/index.html
###### Back-end 🤓:
Оркестратор:
Агент:


------------------------------------------------------------------------------------------

# Адреса для отправки запросов:

- http://localhost:5500/expression/add

		  Что делает:
			  Хэндлит ввод выражения

		  Пример запроса curl:
			  curl -X POST -H "Content-Type: application/json" -d '{"2+2"}' http://localhost:5500/expression/add

			  После -d json может содержать любую строку с выражением, примеры:
				  {"2+2"} - правильно введенное выражение
				  {"1+2+3="} - правильно введенное выражени
				  {"2+2=50"} - неправильно введенное выражени
				  {"2*2-5"} - правильно введенное выражени

- http://localhost:5500/expressions/get-all

			Что делает:
				Хэндлит запросы на получение списка всех выражений

			Пример запроса curl:
				curl -X GET -H "Content-Type: application/json" http://localhost:5500/expressions/get-all

- http://localhost:5500/expressions/get-id

			Что делает:
				Хэндлит запросы на получение выражения по его ID

			Пример запроса curl:
				curl -X GET -H "Content-Type: application/json" http://localhost:5500/expressions/get-id?id=IDHere

- http://localhost:5500/operations/get

		Что делает:
			Хэндлит запрос на получение значения операций

		Пример запроса curl:
			curl -X GET -H "Content-Type: application/json" http://localhost:5500/operations/get?operation=all
		Значения operation: all, plus, minus, multiply, divide, agent

- http://localhost:5500/operations/add

		Что делает:
			Хэндлит запрос с значениями операций

		Пример запроса Curl:
			curl -X POST -H "Content-Type: application/json" -d "[{\"OperationName\":\"plus\",\"OperationDuration\":10},{\"OperationName\":\"minus\",\"OperationDuration\":15},{\"OperationName\":\"multiply\",\"OperationDuration\":20},{\"OperationName\":\"divide\",\"OperationDuration\":25},{\"OperationName\":\"agent\",\"OperationDuration\":30}]" http://localhost:5500/operations/add
		Обязательно должны быть только все приведенные в примере Имена операторов(OperationName)!

- http://localhost:5500/task/get
- http://localhost:5500/result/get