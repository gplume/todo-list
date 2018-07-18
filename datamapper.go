package main

type dataMapper interface {
	close()
	todosMapper
}

type todosMapper interface {
	createTodo(*todo) error
	listTodos() ([]*todo, error)
	updateTodo(*todo) error
	deleteTodo(*todo) error
}
