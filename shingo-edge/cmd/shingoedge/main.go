package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"shingoedge/config"
	"shingoedge/engine"
	"shingoedge/messaging"
	"shingo/protocol"
	"shingoedge/store"
	"shingoedge/www"
)

func main() {
	configPath := flag.String("config", "shingoedge.yaml", "path to config file")
	debug := flag.Bool("debug", false, "enable debug logging")
	port := flag.Int("port", 0, "HTTP port (overrides config)")
	flag.Parse()

	if *debug {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	}

	// Load config
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if *port > 0 {
		cfg.Web.Port = *port
	}

	// Open database
	db, err := store.Open(cfg.DatabasePath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	// Create and start engine
	eng := engine.New(engine.Config{
		AppConfig:  cfg,
		ConfigPath: *configPath,
		DB:         db,
		LogFunc:    log.Printf,
		Debug:      *debug,
	})
	eng.Start()
	defer eng.Stop()

	// Set up messaging
	msgClient := messaging.NewClient(&cfg.Messaging)
	if err := msgClient.Connect(); err != nil {
		log.Printf("messaging connect: %v (will retry via outbox)", err)
	} else {
		// Start outbox drainer
		drainer := messaging.NewOutboxDrainer(db, msgClient, &cfg.Messaging)
		drainer.Start()
		defer drainer.Stop()

		// Start inbound subscriber (old protocol)
		sub := messaging.NewSubscriber(msgClient, cfg, eng.OrderManager())
		if err := sub.Start(); err != nil {
			log.Printf("inbound subscriber: %v", err)
		}

		// Protocol ingestor (new unified protocol — runs alongside old subscriber during transition)
		nodeID := cfg.NodeID()
		edgeHandler := messaging.NewEdgeHandler(eng.OrderManager())
		ingestor := protocol.NewIngestor(edgeHandler, func(hdr *protocol.RawHeader) bool {
			return hdr.Dst.Node == nodeID || hdr.Dst.Node == "*"
		})
		if err := msgClient.Subscribe(cfg.Messaging.DispatchTopic, func(data []byte) {
			ingestor.HandleRaw(data)
		}); err != nil {
			log.Printf("protocol ingestor subscribe: %v", err)
		} else {
			log.Printf("protocol ingestor listening on %s (node=%s)", cfg.Messaging.DispatchTopic, nodeID)
		}

		// Heartbeater (registration + periodic heartbeat)
		hb := messaging.NewHeartbeater(msgClient, nodeID, cfg.Namespace, "dev", []string{cfg.LineID}, cfg.Messaging.OrdersTopic)
		hb.Start()
		defer hb.Stop()
	}
	defer msgClient.Close()

	// Set up HTTP server
	router, stopWeb := www.NewRouter(eng)
	defer stopWeb()

	addr := fmt.Sprintf("%s:%d", cfg.Web.Host, cfg.Web.Port)
	server := &http.Server{Addr: addr, Handler: router}

	// Start HTTP server
	go func() {
		log.Printf("ShinGo Edge listening on %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server: %v", err)
		}
	}()

	// Wait for shutdown signal
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutting down...")
}
