package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-redis/redis/v7"
)

const joinedServerKey = "servers"

// Database type acts as the control interface for the Redis datastore.
type Database struct {
	uri    string
	client *redis.Client
}

// New creates a database connection.
func New(uri string) Database {
	if strings.HasPrefix(uri, "redis://") {
		uri = uri[8:]
	}

	client := redis.NewClient(&redis.Options{
		Addr:     uri,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatalf("Could not connect to database: %s due to err: %s", uri, err)
	}

	log.Printf("Connected to database at %s", uri)

	return Database{
		uri:    uri,
		client: client,
	}
}

// GetBestCandidateWord returns the best replacement synonym for supplied word.
// The return order should be in the defined lexeme order in Lexeme.go.
func (d *Database) GetBestCandidateWord(word string) string {
	results := make([]*redis.StringCmd, 4)

	_, err := d.client.TxPipelined(func(pipe redis.Pipeliner) error {
		pipe.Expire(fmt.Sprintf("best_word_single_%s", word), 10*time.Second)

		for idx, l := range ordering {
			res := pipe.SRandMember(fmt.Sprintf("%s:%s", l, word))
			results[idx] = res
		}

		return nil
	})

	if err != redis.Nil && err != nil {
		log.Printf("Could not access datastore for word: %s, %s", word, err)
		return word
	}

	for _, c := range results {
		w, err := c.Result()
		if err == nil {
			return w
		}
	}

	// Fallback. Only return if nothing was found.
	return word
}

// WaitForReady waits for `ready` status message in `status` pubsub channel.
func (d *Database) WaitForReady(timeout int) error {
	if timeout == 0 {
		log.Println("Skipping database check...")
		return nil
	}

	pubsub := d.client.Subscribe("status")
	defer pubsub.Close()

	log.Printf("Waiting for ready status on channel `status` for %ds", timeout)

	ch := pubsub.Channel()
	for {
		select {
		case msg := <-ch:
			if msg.Payload == "ready" {
				return nil
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			return fmt.Errorf("Channel `status` timed out after %ds", timeout)
		}
	}
}

// AddJoinedServer adds a server to the `servers` set.
func (d *Database) AddJoinedServer(id string) error {
	result := d.client.SAdd(joinedServerKey, id)
	return result.Err()
}

// IsServerJoined checks if a specific guild is in the `servers` set.
func (d *Database) IsServerJoined(id string) (bool, error) {
	result := d.client.SIsMember(joinedServerKey, id)

	if result.Err() != nil {
		return false, result.Err()
	}

	return true, nil
}
