# Go InfraKit

**Go InfraKit** is a modular SDK written in Go designed to unify and abstract access to infrastructure and artificial
intelligence services in distributed systems.

It provides standardized interfaces and interchangeable drivers for different providers and technologies, allowing
developers to configure and use services such as cache, queues, databases, OCR, TTS, STT, and LLM models in a pluggable,
decoupled, and extensible way.

---

## 🔍 Purpose

Go InfraKit was created to eliminate direct coupling with specific libraries and simplify the consumption of
infrastructure and intelligence services by centralizing initialization, fallback, connection, and reconnection logic in
a single SDK.

This way, it is possible to change the implementation of a service (e.g., switch Redis for in-memory cache, RabbitMQ for
Kafka, or use Ollama instead of OpenAI) without altering the application's business logic.

---

## 🧩 Features

- **Unified cache**: with support for Redis and in-memory, using a single interface.
- **Decoupled messaging**: drivers for RabbitMQ, Kafka, NATS, and others through a standard abstraction.
- **Pluggable database**: support for PostgreSQL, MySQL, MongoDB, and others with configuration based on options.
- **Intelligent services**:
    - **OCR**: integration with tools like Tesseract.
    - **Text-to-Speech (TTS)**: abstractions for voice services (AWS, Azure, etc.).
    - **Speech-to-Text (STT)**: interchangeable speech recognition.
    - **LLM**: integration with language models such as OpenAI, Ollama, or local LLMs.

---

## ✅ Advantages

- Reduction of duplicated code in projects using similar services.
- Technology swapping without impact on the main application.
- Interfaces oriented towards software engineering best practices (SOLID).
- Facilitates testing and simulations with mocks and in-memory implementations.
- Modularity for future growth and extension (new drivers, middlewares, etc.).

---

## 📦 Basic example

### Initializing a Redis cache

```go
package main

import (
	"context"
	"fmt"
	"log"

	infrakit "github.com/carlosealves2/go-infrakit"
	"github.com/carlosealves2/go-infrakit/cache"
	"github.com/carlosealves2/go-infrakit/eventbus"
)

func main() {
	ctx := context.Background()

	redisCache, err := infrakit.NewCache(cache.Options{
		Driver: cache.RedisDrive,
		Addr:   "localhost:6379",
		DB:     0,
	})
	if err != nil {
		log.Fatalf("error initializing cache: %v", err)
	}

	err = redisCache.Set(ctx, "key", "value")
	if err != nil {
		log.Fatalf("error setting value in cache: %v", err)
	}

	val, err := redisCache.Get(ctx, "key")
	if err != nil {
		log.Fatalf("error getting value from cache: %v", err)
	}

	fmt.Println("Value retrieved from cache:", val)

	mq, err := infrakit.NewEventBus(eventbus.Options{
		Driver: eventbus.RabbitMQDrive,
		URL:    "amqp://guest:guest@localhost:5672/",
		Queue:  "events",
	})
	if err != nil {
		log.Fatalf("error initializing messaging: %v", err)
	}

	err = mq.Publish(ctx, []byte("Hello world!"))
	if err != nil {
		log.Fatalf("error publishing message: %v", err)
	}

	fmt.Println("Message published successfully!")
}
```

## 📝 License

Distributed under the MIT license.
