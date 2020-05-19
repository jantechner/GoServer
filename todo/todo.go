package todo

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

type Todo struct {
	Number  int    `json:"id"`
	Content string `json:"content"`
}

func (todo *Todo) save() error {
	filename := fmt.Sprintf("./todo/todos/%v.txt", todo.Number)
	return ioutil.WriteFile(filename, []byte(todo.Content), os.FileMode(0644))
}

func load(idStr string) (*Todo, error) {
	filename := idStr + ".txt"
	content, err := ioutil.ReadFile("./todo/todos/" + filename)
	if err != nil {
		return nil, errors.New("todo not found")
	}
	number, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, err
	}
	return &Todo{number, string(content)}, nil
}

func loadAll() (todos []Todo, err error) {
	files, err := ioutil.ReadDir("./todo/todos")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		filename := file.Name()
		numberStr := filename[:len(filename)-len(".txt")]
		todo, err := load(numberStr)
		if err != nil {
			return nil, err
		}
		todos = append(todos, *todo)
	}
	return
}
