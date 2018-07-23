package main

type dataMapper interface {
	close()
	todosMapper
}

type todosMapper interface {
	saveTodo(*todo) error
	listTodos() ([]*todo, error)
	getTodo(int) (*todo, error)
	updateTodo(*todo) error
	deleteTodo(int) error
}
