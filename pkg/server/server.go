package textserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"

	"word-histrogram/pkg/tasks"
	"word-histrogram/pkg/worker"
)

// TextServer listens and serves a tcp endpoint
type TextServer struct {
	addListener net.Listener
	getListener net.Listener
	pool        worker.Pool
	handler     *tasks.Handler
}

// New returns a new TextServer
func New(addPort int, getPort int, pool worker.Pool, handler *tasks.Handler) (*TextServer, error) {
	add, err := net.Listen("tcp", "localhost:"+strconv.Itoa(addPort))
	if err != nil {
		return nil, fmt.Errorf("could not create addListener: %w", err)
	}

	get, err := net.Listen("tcp", "localhost:"+strconv.Itoa(getPort))
	if err != nil {
		return nil, fmt.Errorf("could not create getListener: %w", err)
	}

	return &TextServer{
		addListener: add,
		getListener: get,
		pool:        pool,
		handler:     handler,
	}, nil
}

// ServeAdd handles new connections and delegates to a worker pool for persistence
func (t *TextServer) ServeAdd(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			t.close()

			return
		default:
			conn, err := t.addListener.Accept()
			if err != nil {
				log.Println("failed to acquire connection from getListener: %w", err)

				return
			}

			buffer := make([]byte, 1024)

			_, err = conn.Read(buffer)
			if err != nil {
				log.Println(err)
			}

			err = t.pool.Schedule(&tasks.WordsItem{
				Buffer:    buffer,
				WordStore: t.handler.WordStore,
			})
			if err != nil {
				log.Println("could not schedule task: %w", err)
			}

			err = conn.Close()
			if err != nil {
				log.Println("failed to close connection: %w", err)
			}
		}
	}
}

// ServeGet retrieves an histogram of saved data
func (t *TextServer) ServeGet(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			t.close()

			return
		default:
			conn, err := t.getListener.Accept()
			if err != nil {
				log.Println("failed to acquire connection from getListener: %w", err)

				return
			}

			wordsStore := t.handler.Get()

			type kv struct {
				Key   string
				Value int
			}

			var topRank []kv
			for k, v := range wordsStore {
				topRank = append(topRank, kv{k, v})
			}

			sort.Slice(topRank, func(i, j int) bool {
				return topRank[i].Value > topRank[j].Value
			})

			resultNum := 1

			for _, kv := range topRank {
				if resultNum == 5 {
					break
				}

				formatted := fmt.Sprintf("%s, %d\n", kv.Key, kv.Value)

				_, err = conn.Write([]byte(formatted))
				resultNum++
			}

			if err != nil {
				log.Println("failed to write response to connection: %w", err)
			}

			err = conn.Close()
			if err != nil {
				log.Println("failed to close connection: %w", err)
			}
		}
	}
}

func (t *TextServer) close() {
	if err := t.addListener.Close(); err != nil {
		log.Println("could not have close the getListener: %w", err)
	}
}
