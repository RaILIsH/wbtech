package main

import "fmt"

type Human struct {
	Name string
	Age  int
}

func (h *Human) SayHello() {
	fmt.Printf("Привет, меня зовут %s и мне %d лет.\n", h.Name, h.Age)
}

// Структура Action, которая встраивает Human
type Action struct {
	Human      // Встраивание структуры Human
	Profession string
}

// Метод структуры Action
func (a *Action) Work() {
	fmt.Printf("%s, %s.\n", a.Name, a.Profession)
}

func main() {

	// Создаем экземпляр Action
	person := Action{
		Human: Human{
			Name: "Иван",
			Age:  30,
		},
		Profession: "Программист",
	}

	// Вызываем методы Human через Action
	person.SayHello() // Метод Human доступен через Action

	// Вызываем метод Action
	person.Work()

	// Обращение к полям Human
	fmt.Printf("Имя: %s, Возраст: %d, Профессия: %s\n",
		person.Name, person.Age, person.Profession)

}
