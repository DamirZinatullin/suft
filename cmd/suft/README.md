###NAME:  
SUFT CLI - CLI предоставляет возможность взаимодействия с api СУФТ (системы учета фактических трудозатрат)

###USAGE:
    suft [global options] command [command options] [arguments...]

###COMMANDS:
    help, h  Shows a list of commands or help for one command

####Временные затраты:  
    logging-times, lts          Список временных затрат  
    logging-time, lt            Детализация временной затраты  
    add-logging-time, al        Добавление временной затраты  
    remove-logging-time, rmlt   Удаление временной затраты  
    approve-logging-time, aprv  Утверждение временной затраты

####Клиент:
    login   Аутентификация клиента  
    logout  Выход из клиента  

####Расписания:
    schedules, scs         Список расписаний  
    schedule, sc           Детализация расписания  
    add-schedule, as       Добавление расписания  
    submit-for-approve, s  Отправить расписание на утверждение  

###GLOBAL OPTIONS:
    --help, -h  show help
