# go-marathon-team-1

### Репозиторий команды 1

## СУФТ SDK

Go-библиотека, CLI и набор тестов для работы с СУФТ (Системой Учёта Фактических Трудозатрат).
Библиотека позволяет:
* Получать информацию по расписаниям
* Создавать новые расписания
* Добавлять трудозатраты к расписаниям
* Удалять трудозатраты
* Отправлять расписания на утверждение
* Утверждать трудозатраты по расписаниям
* Отклонять трудозатраты по расписаниям

### Пример использования библиотеки:
```
package main

import (
	"fmt"
	"log"
	"suftsdk/pkg/api"
)

func main() {
	client, err := api.NewClient("demo@example.com", "demo", nil)
	if err != nil {
		log.Fatalln(err)
	}

    // здесь вызываем метод получения списка расписаний без передачи опций вывода результата.
    // в этом случае берутся значения по умолчанию, указанные в библиотеке.
	schedules, err := client.Schedules(nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Список расписаний\n")
	for _, schedule := range schedules {
		fmt.Printf("%#v\n", schedule)
	}

    // здесь вызываем метод получения списка расписаний с передачей опций вывода результата.
    // нумерация страниц начинается с нуля.
	schedules, err = client.Schedules(&api.OptionsS{Page: 1, Size: 2})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("\nСписок расписаний при передаче опций\n")
	for _, schedule := range schedules {
		fmt.Printf("%#v\n", schedule)
	}
}
```
Другие примеры вы можете найти в папке examples.

## CLI
SUFT CLI - CLI предоставляет возможность взаимодействия с api СУФТ (системы учета фактических трудозатрат)

### Установка:
> go install cmd/suft/suft.go

### Пример использования CLI:
    suft [global options] command [command options] [arguments...]

### COMMANDS:
    help, h  Shows a list of commands or help for one command

#### Временные затраты:  
    logging-times, lts          Список временных затрат  
    logging-time, lt            Детализация временной затраты  
    add-logging-time, al        Добавление временной затраты  
    remove-logging-time, rmlt   Удаление временной затраты  
    approve-logging-time, aprv  Утверждение временной затраты
    decline-logging-time, dcl   Отклонение временной затраты

#### Клиент:
    login   Аутентификация клиента  
    logout  Выход из клиента  

#### Расписания:
    schedules, scs         Список расписаний  
    schedule, sc           Детализация расписания  
    add-schedule, as       Добавление расписания  
    submit-for-approve, s  Отправить расписание на утверждение

### GLOBAL OPTIONS:
    --help, -h  show help

*Авторы: Зинатуллин Дамир, Цокало Жан*