package todo

import (
	"io/ioutil"
	"strconv"
	"sync/atomic"
)

type TodosCounter struct {
	value int32
}

var Counter TodosCounter

func (counter *TodosCounter) Init() error {
	files, err := ioutil.ReadDir("./todo/todos")
	if err != nil {
		return err
	}
	var max int
	for _, file := range files {
		filename := file.Name()
		numberStr := filename[:len(filename)-len(".txt")]
		num, err := strconv.Atoi(numberStr)
		if err != nil {
			return err
		}
		if num > max {
			max = num
		}
	}
	counter.value = int32(max)
	return nil
}

func (counter *TodosCounter) GetNewId() int {
	return int(atomic.AddInt32(&counter.value, 1))
}