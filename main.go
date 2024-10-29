package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Estrutura que representa uma tarefa para ser processada
type Task struct {
	id      int
	payload int
}

// Função para processar uma tarefa simulando tempo de processamento
func (t *Task) Process() {
	fmt.Printf("Processando tarefa #%d com carga %d...\n", t.id, t.payload)
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond) // Simula um trabalho que leva um tempo aleatório
	fmt.Printf("Tarefa #%d concluída\n", t.id)
}

// Worker que processa tarefas do canal de trabalho
func worker(id int, tasks <-chan *Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		fmt.Printf("Worker %d começou a processar tarefa #%d\n", id, task.id)
		task.Process()
	}
}

// Pool de workers que distribui as tarefas entre workers disponíveis
func workerPool(numWorkers int, tasks []*Task) {
	// Cria um canal para enviar as tarefas aos workers
	taskChannel := make(chan *Task, len(tasks))
	// WaitGroup para aguardar todos os workers terminarem
	var wg sync.WaitGroup

	// Inicializa os workers
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, taskChannel, &wg)
	}

	// Envia as tarefas para o canal de trabalho
	for _, task := range tasks {
		taskChannel <- task
	}
	close(taskChannel) // Fecha o canal para indicar que não há mais tarefas

	// Espera todos os workers concluírem
	wg.Wait()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Cria uma lista de tarefas
	tasks := []*Task{}
	for i := 1; i <= 10; i++ {
		tasks = append(tasks, &Task{
			id:      i,
			payload: rand.Intn(100),
		})
	}

	// Inicia o pool de workers com número de workers definido
	fmt.Println("Iniciando processamento de tarefas com pool de workers...")
	workerPool(3, tasks) // Exemplo com 3 workers

	fmt.Println("Todas as tarefas foram processadas.")
}
