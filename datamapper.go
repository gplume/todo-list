package main

type dataMapper interface {
	close()
	saveTodo(*todo) error
	listTodos(string) ([]*todo, error)
	getTodo(int) (*todo, error)
	updateTodo(*todo) error
	deleteTodo(int) error
}
