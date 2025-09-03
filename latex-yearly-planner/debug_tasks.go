package mainpackage main

import (

import (	"fmt"

	"fmt"	"github.com/kudrykv/latex-yearly-planner/internal/data"

	"github.com/kudrykv/latex-yearly-planner/internal/data")

)

func main() {

func main() {	reader := data.NewReader("../input/test_triple.csv")

	reader := data.NewReader("../input/test_triple.csv")	tasks, err := reader.ReadTasks()

	tasks, err := reader.ReadTasks()	if err != nil {

	if err != nil {		fmt.Printf("Error: %v

		fmt.Printf("Error: %v\n", err)", err)

		return		return

	}	}

	fmt.Printf("Read %d tasks:\n", len(tasks))	fmt.Printf("Read %d tasks:

	for i, task := range tasks {", len(tasks))

		fmt.Printf("Task %d: ID=%s, Name=%s, Category=%s, Start=%s, End=%s\n", 	for i, task := range tasks {

			i+1, task.ID, task.Name, task.Priority, task.StartDate.Format("2006-01-02"), task.EndDate.Format("2006-01-02"))		fmt.Printf("Task %d: ID=%s, Name=%s, Category=%s, Start=%s, End=%s

	}", 

}			i+1, task.ID, task.Name, task.Priority, task.StartDate.Format("2006-01-02"), task.EndDate.Format("2006-01-02"))
	}
}
