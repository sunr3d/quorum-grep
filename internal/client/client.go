package client

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sunr3d/quorum-grep/internal/config"
	"github.com/sunr3d/quorum-grep/models"
	pbg "github.com/sunr3d/quorum-grep/proto/grepsvc"
)

const capacity = 1024

type Client struct {
	servers   []string
	quorum    int
	timeout   time.Duration
	chunkSize int
}

// New - конструктор Client.
func New(cfg *config.Config) *Client {
	timeout, _ := time.ParseDuration(cfg.Client.Timeout)
	quorum := len(cfg.Client.ServerList)/2 + 1

	return &Client{
		servers:   cfg.Client.ServerList,
		quorum:    quorum,
		timeout:   timeout,
		chunkSize: cfg.Client.ChunkSize,
	}
}

// ProcessFile - обрабатывает файл.
// Разбивает на чанки и отправляет на серверы.
// Ожидает результатов от серверов и собирает их в один результат.
// Выводит результат в консоль.
func (c *Client) ProcessFile(filename string, opts models.GrepOptions) error {
	lines, err := c.readInput(filename)
	if err != nil {
		return fmt.Errorf("readInput: %w", err)
	}

	tasks := c.splitData(lines, len(c.servers), opts)

	results, errors := c.sendToServers(tasks)

	out, err := c.waitForQuorum(results, errors)
	if err != nil {
		return fmt.Errorf("waitForQuorum: %w", err)
	}

	c.printResults(out, opts)

	return nil
}

// readInput - читает входные данные из файла или stdin.
func (c *Client) readInput(filename string) ([][]byte, error) {
	var scanner *bufio.Scanner

	if filename == "-" || filename == "" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return nil, fmt.Errorf("ошибка открытия файла %s: %w", filename, err)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	lines := make([][]byte, 0, capacity)
	for scanner.Scan() {
		line := make([]byte, len(scanner.Bytes()))
		copy(line, scanner.Bytes())
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения %s: %w", filename, err)
	}

	return lines, nil
}

// splitData - разбивает данные на чанки.
// Обрабатывает перекрытие контекста.
func (c *Client) splitData(lines [][]byte, numServers int, opts models.GrepOptions) []models.Task {
	lineLen := len(lines)
	chunkSize := c.chunkSize
	if chunkSize <= 0 || chunkSize > lineLen {
		chunkSize = lineLen / numServers
		if chunkSize == 0 {
			chunkSize = 1
		}
	}

	contextOverlap := opts.After
	if opts.Before > contextOverlap {
		contextOverlap = opts.Before
	}
	if opts.Around > contextOverlap {
		contextOverlap = opts.Around
	}

	tasks := make([]models.Task, numServers)

	for i := range numServers {
		start := i * chunkSize
		end := start + chunkSize

		if i == numServers-1 {
			end = lineLen
		}

		if i > 0 {
			start -= contextOverlap
			if start < 0 {
				start = 0
			}
		}
		if i < numServers-1 {
			end += contextOverlap
			if end > lineLen {
				end = lineLen
			}
		}

		if start < lineLen {
			tasks[i].Data = bytes.Join(lines[start:end], []byte("\n"))
			tasks[i].Index = i
			tasks[i].Options = opts

			tasks[i].LineNumbers = make([]int64, end-start)
			for j := range tasks[i].LineNumbers {
				tasks[i].LineNumbers[j] = int64(start + j + 1)
			}
		}
	}

	return tasks
}

// buildRequest - строит gRPC запрос для отправки на сервер.
func (c *Client) buildRequest(id int, task models.Task) *pbg.ChunkRequest {
	return &pbg.ChunkRequest{
		TaskId:      fmt.Sprintf("task-%d", id),
		Data:        task.Data,
		ChunkIndex:  int64(task.Index),
		LineNumbers: task.LineNumbers,
		Options: &pbg.GrepOptions{
			Pattern:    task.Options.Pattern,
			After:      int64(task.Options.After),
			Before:     int64(task.Options.Before),
			Around:     int64(task.Options.Around),
			Count:      task.Options.Count,
			IgnoreCase: task.Options.IgnoreCase,
			Invert:     task.Options.Invert,
			Fixed:      task.Options.Fixed,
			LineNum:    task.Options.LineNum,
		},
	}
}

// sendToServers - отправляет чанки на серверы в горутинах.
// Устанавливает соединение с сервером и отправляет запрос.
// Ожидает результатов от серверов и собирает их в один результат.
// Возвращает результаты и ошибки.
func (c *Client) sendToServers(tasks []models.Task) ([]models.Result, []error) {
	results := make([]models.Result, len(tasks))
	errors := make([]error, len(tasks))

	var wg sync.WaitGroup

	for i, task := range tasks {
		wg.Add(1)
		go func(i int, task models.Task) {
			defer wg.Done()

			server := c.servers[i%len(c.servers)]

			conn, err := grpc.NewClient(server, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				errors[i] = fmt.Errorf("не удалось подключиться к серверу %s: %w", server, err)
				return
			}
			defer conn.Close()

			client := pbg.NewGrepServiceClient(conn)

			req := c.buildRequest(i, task)

			ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
			defer cancel()

			resp, err := client.ProcessChunk(ctx, req)
			if err != nil {
				errors[i] = fmt.Errorf("ошибка при обработке куска %d на сервере %s: %w", i, server, err)
				return
			}

			matches := make([]models.Match, len(resp.Matches))
			for i, match := range resp.Matches {
				matches[i] = models.Match{
					Content:    match.Content,
					LineNumber: match.LineNumber,
				}
			}

			results[i] = models.Result{
				Matches:    matches,
				MatchCount: int(resp.MatchCount),
				Error:      resp.Error,
				TaskIndex:  i,
			}
		}(i, task)
	}

	wg.Wait()

	return results, errors
}

// waitForQuorum - ожидает результатов от серверов и собирает их в один результат.
// Возвращает результаты и ошибки.
func (c *Client) waitForQuorum(results []models.Result, errors []error) ([]models.Match, error) {
	success := 0
	seen := make(map[int64]bool)
	var out []models.Match

	for i, result := range results {
		if errors[i] == nil {
			for _, match := range result.Matches {
				if !seen[match.LineNumber] {
					seen[match.LineNumber] = true
					out = append(out, match)
				}
			}
			success++
		}
	}

	if success < c.quorum {
		return nil, fmt.Errorf("недостаточно успешных результатов")
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].LineNumber < out[j].LineNumber
	})

	return out, nil
}

// printResults - выводит результаты в консоль.
func (c *Client) printResults(matches []models.Match, opts models.GrepOptions) {
	if opts.Count {
		fmt.Println(len(matches))
	} else {
		for _, match := range matches {
			if opts.LineNum {
				fmt.Printf("%d:%s\n", match.LineNumber, string(match.Content))
			} else {
				fmt.Printf("%s\n", string(match.Content))
			}
		}
	}
}
