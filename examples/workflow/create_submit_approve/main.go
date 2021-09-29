package main

import (
	"fmt"
	"log"
	"suftsdk/pkg/api"
)

func main() {
	// инициализируем новый клиент с учётными данными Пака Юрия
	client1, err := api.NewClient("pakua", "147753", nil)
	if err != nil {
		log.Fatalln(err)
	}

	// добавляем новое расписание и выводим его на экран
	periodId := api.PeriodId(378)
	schedule, err := client1.AddSchedule(periodId)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Созданное расписание")
	fmt.Printf("%#v\n\n", *schedule)

	// добавляем новый logging-time к созданному расписанию
	loggingTime1 := api.AddLoggingTime{
		CommentEmployee: "test15",
		Day1Time:        1,
		Day2Time:        1,
		Day3Time:        0,
		Day4Time:        1,
		Day5Time:        0,
		Day6Time:        3,
		Day7Time:        2,
		ProjectId:       69753,
		Task:            "test15",
		WorkKindId:      21,
	}
	loggingTimeCreated1, err := client1.AddLoggingTime(api.ScheduleId(schedule.Id), &loggingTime1)
	if err != nil {
		log.Fatalln(err)
	}

	// добавляем ещё один logging-time к созданному расписанию
	loggingTime2 := api.AddLoggingTime{
		CommentEmployee: "test16",
		Day1Time:        5,
		Day2Time:        0,
		Day3Time:        4,
		Day4Time:        0,
		Day5Time:        3,
		Day6Time:        0,
		Day7Time:        2,
		ProjectId:       69753,
		Task:            "test16",
		WorkKindId:      21,
	}
	loggingTimeCreated2, err := client1.AddLoggingTime(api.ScheduleId(schedule.Id), &loggingTime2)
	if err != nil {
		log.Fatalln(err)
	}

	//получаем список объектов logging-time, относящихся к нашему расписанию и выводим их на экран
	loggingTimes, err := client1.LoggingTimeList(api.ScheduleId(schedule.Id), nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Список созданных Logging-Time'ов:\n")
	for _, lt := range loggingTimes {
		fmt.Println(*lt)
	}
	fmt.Println()
	fmt.Println()

	// отправляем расписание на утверждение руководителю/согласующему
	// и выводим его на экран, чтобы показать, что изменился его статус
	// здесь мы вызвали метод непосредственно у расписания
	scheduleForApprove, err := schedule.SubmitForApproveSchedule()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Расписание, отправленное на утверждение")
	fmt.Printf("%#v\n\n", *scheduleForApprove)

	//показываем, что у всех объектов logging-time, относящихся к нашему расписанию,
	//также изменился статус на НУ. Выводим их на экран
	loggingTimes, err = client1.LoggingTimeList(api.ScheduleId(schedule.Id), nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Список Logging-Time'ов, находящихся на утверждении:\n")
	for _, lt := range loggingTimes {
		fmt.Println(*lt)
	}
	fmt.Println()
	fmt.Println()

	// инициализируем новый экземпляр клиента с учётными данными Никонова Олега - руководителя Пака
	client2, err := api.NewClient("nikonovov", "147753", nil)
	if err != nil {
		log.Fatalln(err)
	}

	// получаем объекты logging-time и утверждаем их. Выводим на экран
	loggingTimeForApproveToAdmin1, err := client2.DetailLoggingTime(api.ScheduleId(scheduleForApprove.Id), api.LoggingTimeId(loggingTimeCreated1.Id))
	if err != nil {
		log.Fatalln(err)
	}
	loggingTimeApproved1, err := loggingTimeForApproveToAdmin1.ApproveLoggingTime("всё хорошо")
	if err != nil {
		log.Fatalln(err)
	}
	loggingTimeForApproveToAdmin2, err := client2.DetailLoggingTime(api.ScheduleId(scheduleForApprove.Id), api.LoggingTimeId(loggingTimeCreated2.Id))
	if err != nil {
		log.Fatalln(err)
	}
	loggingTimeApproved2, err := loggingTimeForApproveToAdmin2.ApproveLoggingTime("всё хорошо")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Список утверждённых объектов трудозатрат")
	fmt.Printf("%#v\n%#v\n\n", *loggingTimeApproved1, *loggingTimeApproved2)

	// После утверждения всех трудозатрат наш изначальный объект расписания
	// автоматически сам переводится в статус УТВ. Посмотрим его.
	scheduleApproved, err := client2.DetailSchedule(api.ScheduleId(scheduleForApprove.Id))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("Утверждённое расписание\n")
	fmt.Printf("%#v\n\n", *scheduleApproved)
}
