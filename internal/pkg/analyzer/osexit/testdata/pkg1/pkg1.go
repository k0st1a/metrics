package main

import "os"

func main() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	os.Exit(0) // want "direct call os.Exit"
}

func Exit() {
	os.Exit(0)
}
