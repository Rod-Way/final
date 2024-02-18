function showContent(contentId) {
    // Скрываем все контейнеры
    var allContents = document.querySelectorAll('main > div')
    allContents.forEach(function(content) {
        content.style.display = 'none';
    });

    // Показываем выбранный контейнер
    var selectedContent = document.getElementById(contentId);
    if (selectedContent) {
        selectedContent.style.display = 'block';
    }
}

function preLoad(){
    // Вызываем функцию для отправки GET-запроса
    getAllExpressions()
        .then(expressionsData => {
            generateExpressionCards("expressions-cards", expressionsData);
        })
        .catch(error => {
            console.error('Error getting expressions:', error);
        });

    getAllOperations()
        .then(operations => {
            console.log(operation)
            for (var i = 0; i < operations.length; i++) {
                // Получаем текущий элемент массива
                var operation = operations[i];
                document.getElementById(operation.OperationName).value = operation.OperationDuration;
            }
        })
        .catch(error => {
            console.error('Error getting expressions:', error);
        });

    getAllAgents()
        .then(agentsData => {
            generateAgentCards("agents-cards", agentsData);
        })
        .catch(error => {
            console.error('Error getting expressions:', error);
        });
}

// Хэндлим обновление страницы
document.addEventListener('DOMContentLoaded', function() {

    // Отображаем один из контейнеров
    showContent('expressions');

    preLoad();
});

const statusColors = {
    ok: '#4CAF50',       // Зеленый
    bad: '#D32F2F',      // Красный
    lost: '#757575',     // Серый
    inactive: '#7B1FA2', // Фиолетовый
    error: '#212121'     // Очень темный цвет
};

function setBackgroundColor(element, color) {
    element.style.backgroundColor = color;
}

function generateExpressionCards(containerId, expressions) {
    const container = document.getElementById(containerId);
    container.innerHTML = '';

    expressions.forEach((expression, index) => {
        const card = document.createElement('div');
        card.classList.add('card');

        // Создаем заголовок для карточки с ExpressionData
        const title = document.createElement('h2');
        title.textContent = `${expression.Result} \n\n${expression.Status}`;
        card.appendChild(title);

        // Добавляем информацию о выражении
        const expressionInfo = document.createElement('p');
        expressionInfo.textContent = `Status: ${expression.Status}\n\nExpression: ${expression.ExpressionData}\n\nCreated Time: ${expression.CreatedTime}\n\nLast Update Time: ${expression.LastUpdateTime}`;
        card.appendChild(expressionInfo);

        if (expression.Status === 'Active') {
            setBackgroundColor(card, statusColors["ok"]);
        } else if (expression.Status === 'Done') {
            setBackgroundColor(card, statusColors["inactive"]);
        } else if (expression.Status === 'Bad') {
            setBackgroundColor(card, statusColors["bad"]);
        } else if (expression.Status === 'ServerError'){
            setBackgroundColor(card, statusColors["lost"]);
        }

        //Меняем background color
        setBackgroundColor(card)

        // Добавляем карточку в контейнер
        container.appendChild(card);
    });
}

function generateAgentCards(containerId, agents) {
    const container = document.getElementById(containerId);
    container.innerHTML = '';

    agents.forEach(agent => {
        const card = document.createElement('div');
        card.classList.add('card');

        // Создаем заголовок для карточки с ID агента
        const title = document.createElement('h2');
        title.textContent = `Agent ID: ${agent.ID}`;
        card.appendChild(title);

        // Добавляем информацию об агенте
        const agentInfo = document.createElement('p');
        agentInfo.textContent = `Status: ${agent.Status}\n\nAddress: ${agent.Addres}\n\nLast Ping: ${agent.LastPing}`;
        card.appendChild(agentInfo);

        // Устанавливаем картинку в зависимости от статуса агента
        const statusImage = document.createElement('img');
        if (agent.Status === 'Active') {
            statusImage.src = 'src/goodDaemon.png';
            statusImage.alt = 'Active';
            setBackgroundColor(card, statusColors["ok"]);
        } else if (agent.Status === 'Inactive') {
            statusImage.src = 'src/leftedDaemon.png';
            statusImage.alt = 'Inactive';
            setBackgroundColor(card, statusColors["inactive"]);
        } else if (agent.Status === 'Lost') {
            statusImage.src = 'src/lostDaemon.png';
            statusImage.alt = 'Unknown';
            setBackgroundColor(card, statusColors["bad"]);
        } else if (agent.Status === 'Disconnected'){
            statusImage.src = 'src/lostDaemon.png';
            statusImage.alt = 'Lost';
            setBackgroundColor(card, statusColors["lost"]);
        }
        card.appendChild(statusImage);

        // Добавляем карточку в контейнер
        container.appendChild(card);
    });
}

// Отслеживаем ввод выражения
var expressionInput = document.getElementById('expressionInput');
expressionInput.addEventListener('keydown', function(event) {
    if (event.key === 'Enter') {
        var enteredValue = expressionInput.value;
        document.getElementById('expressionInput').value = "";
        console.log('Введенное значение:', enteredValue);

        postExpression(enteredValue);
    }
});

// Отслеживаем ввод времени действия
var formNum = document.getElementById('operationsTimeForm');
formNum.addEventListener('submit', function(event) {
    event.preventDefault(); // Предотвращаем отправку формы

    postOperations();
});

function getAllExpressions() {
    return new Promise((resolve, reject) => {
        var xhr = new XMLHttpRequest();
        var url = "http://localhost:5500/expressions/get-all";
        xhr.open("GET", url, true);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    var responseData = JSON.parse(xhr.responseText);
                    resolve(responseData);
                } else {
                    reject(new Error("Failed to fetch expressions. Status code: " + xhr.status));
                }
            }
        };
        xhr.send();
    });
};

function getAllAgents() {
    return new Promise((resolve, reject) => {
        var xhr = new XMLHttpRequest();
        var url = "http://localhost:5500/agents/get";
        xhr.open("GET", url, true);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    var responseData = JSON.parse(xhr.responseText);
                    resolve(responseData);
                } else {
                    reject(new Error("Failed to fetch expressions. Status code: " + xhr.status));
                }
            }
        };
        xhr.send();
    });
};

function getAllOperations() {
    return new Promise((resolve, reject) => {
        var xhr = new XMLHttpRequest();
        var url = "http://localhost:5500/operations/get?operation=all";
        xhr.open("GET", url, true);
        xhr.onreadystatechange = function () {
            if (xhr.readyState === 4) {
                if (xhr.status === 200) {
                    var responseData = JSON.parse(xhr.responseText);
                    resolve(responseData);
                } else {
                    reject(new Error("Failed to fetch expressions. Status code: " + xhr.status));
                }
            }
        };
        xhr.send();
    });
};

function postOperations() {
    // Создаем массив объектов данных
    var data = [
        {OperationName: 'plus',OperationDuration: parseInt(document.getElementById("plus").value)},
        {OperationName: 'minus', OperationDuration: parseInt(document.getElementById("minus").value)},
        {OperationName: 'multiply', OperationDuration: parseInt(document.getElementById("multiply").value)},
        {OperationName: 'divide', OperationDuration: parseInt(document.getElementById("divide").value)},
        {OperationName: 'agent', OperationDuration: parseInt(document.getElementById("agent").value)}
    ];

    // Настраиваем параметры запроса
    var requestOptions = {
        method: 'POST', // Указываем метод запроса
        headers: {
            'Content-Type': 'application/json' // Указываем тип контента
        },
        body: JSON.stringify(data) // Преобразуем данные в JSON и устанавливаем их в тело запроса
    };

    // Выполняем запрос на сервер
    fetch('http://localhost:5500/operations/add', requestOptions)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            console.log(response)
            return response.json();
        })
        .then(data => {
            console.log(data); // Выводим полученные данные в консоль
        })
        .catch(error => {
            console.error('There was an error!', error); // Выводим сообщение об ошибке в консоль
        });
}

function postExpression(expression = '') {
    // Настраиваем параметры запроса
    var requestOptions = {
        method: 'POST', // Указываем метод запроса
        headers: {
            'Content-Type': 'application/json' // Указываем тип контента
        },
        body: expression // Преобразуем данные в JSON и устанавливаем их в тело запроса
    };

    // Выполняем запрос на сервер
    fetch('http://localhost:5500/expression/add', requestOptions)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            console.log(response)
            return response.json();
        })
        .then(data => {
            console.log(data); // Выводим полученные данные в консоль
        })
        .catch(error => {
            console.error('There was an error!', error); // Выводим сообщение об ошибке в консоль
        });
}
