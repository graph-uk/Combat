package main

import (
	"Combat/CLIParser"
	"os"
	"Combat/SerialRunner"
)

func main() {
	action := CLIParser.GetAction() //"run" action by default
	if action == "" {
		action = "run"
	}

	if action == "help" {
		println("Help")
		os.Exit(0)
	}

	var testManager TestManager

	testManager.Init("Tests_Examples/Tests")

	switch action {
	case "list":
		testManager.PrintListOrderedByNames()
	case "tags":
		testManager.PrintListOrderedByTag()
	case "params":
		testManager.PrintListOrderedByParameter()
	case "cases":
		testManager.PrintCases()
	case "run":
		testManager.PrintCases()
		SerialRunner.RunCasesSerial(testManager.AllCases(),"Tests_Examples/Tests")
	default:
		println("Incorrect action. Please run \"Combat help\" for find available actions.")
		os.Exit(1)
	}
	os.Exit(0)

	//Парсим все значения параметров из CLI.

	//Сначала фильтруем тесты по имени, потом по тегам.
	//combat run -name="lynx" -tags="xnd"
	//-name
	//-tags
	//-locale="sdf"
	//-nyx="sdf"

	//Потом собираем инфу о параметрах, которые нужно предоставить.
	//Читаем параметры.
	//Проверяем, что все необходимые параметры предоставлены. Если параметров предоставлено больше - ошибка.
	//Проверяем, что перечисления вписываются в ограничения.
	//Запускаем/выводим статистику.

	//combat <...> run  запуск
	//combat <...> list вывод тестов с параметрами и тегами, группировка по имени теста.
	//combat <...> params вывод тестов с параметрами и тегами, группировка по параметрам.
	//combat <...> tags вывод тестов с параметрами и тегами, группировка по тегам.
	//combat help вывод справки.

	//Print all tests with params

	return
}
