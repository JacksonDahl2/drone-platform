### Overview
The goal is to build a data ingestion pipeline from drone geospatial data. I want to be able to analyze drone data coming in from many drones all at the same time, push that data to a kafka or redpanda message broker so that no data is lost. That data will then be saved to postgres with the timeseriesDB extension, and will have a websocket connection to a react frontend that's live showing different graphs, such as the location of all the drones.

The goal of this project will be to become more familiar with message brokers, as well as to use Go for the first time to actually build something.

### Steps
1. Use docker compose to spin up redpanda as an easy alternative to kafka. Have a Go consumer and a go producer that will just produce and consume a simple message.
	1. Producer - _Learning Focus:_ Go syntax, `main()` functions, and importing the `segmentio/kafka-go` library.
	2. Consumer - _Learning Focus:_ Go routines (concurrency) and `for` loops.
2. Build the data highway
	1.  Go struct that will define the go data
	2. a small `net/http` server built in go that will accept incoming data - _Learning Focus:_ JSON Marshaling/Unmarshaling and HTTP Handlers.
	3. Write a python mock script 'drone simulator' that will simulate all the data coming in, make sure it's coming into the consumers terminal
3. Persistant storage
	1. save the data so that it can be used later in a timeseriesDb
	2. Write a SQL schema to create a hypertable -> do I want to use a ORM in this instance?
	3. instead of having the consumer print to the terminal, have it store in the db using the `database/sql`, make sure to use goroutines and connection pools -> really good learning opportunity here
4. Visualizations
	1. Get data to the frontend without having the database act as a middleman
	2. Create a new websocket Go server, or add to the existing one, that maintains a websocket connection with the frontend
	3. Every time a message comes from kafka, the server should send it to the db and push it out to the websockets - _Learning Focus:_ Channels in Go (how different parts of your code talk to each other safely).
5. Productionization
	1. Use a local kubernetes tool like `kind` or `minikube`
	2. Write the kubernetes manifests to deploy cluster

drone-platform/
├── cmd/                         # The entry points for your applications
│   ├── ingestion-gateway/       # The "Front Door" (MQTT/HTTP to Kafka)
│   │   └── main.go
│   ├── consumer-worker/         # The "Brain" (Kafka to TimescaleDB)
│   │   └── main.go
│   └── mock-drone-swarm/        # Your Go or Python simulation script
│       └── main.go
├── internal/                    # Private code you don't want others to import
│   ├── platform/                # Shared logic for DB/Kafka connections
│   │   ├── kafka.go
│   │   └── database.go
│   └── telemetry/               # Domain logic (The "Drone" stuff)
│       ├── models.go            # Your DroneTelemetry struct
│       └── store.go             # SQL queries for the drone data
├── deployments/                 # Infrastructure & Orchestration
│   ├── docker-compose.yml       # Local development setup
│   ├── terraform/               # AWS/EKS infrastructure code
│   └── k8s/                     # Kubernetes manifests (YAML)
├── go.mod                       # Go project dependencies
└── go.sum                       # Checksums for dependencies