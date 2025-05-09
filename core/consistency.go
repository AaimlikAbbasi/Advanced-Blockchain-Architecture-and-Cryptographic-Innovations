package core

import "fmt"

// NetworkTelemetry represents network performance metrics
type NetworkTelemetry struct {
	Latency    int     // in milliseconds
	PacketLoss float64 // percentage (0-1)
	Throughput float64 // in Mbps
}

// ConsistencyOrchestrator manages blockchain consistency based on network conditions
type ConsistencyOrchestrator struct {
	NetworkMetrics   NetworkTelemetry
	ConsistencyLevel int // 1-5, where 5 is highest consistency
}

// NewConsistencyOrchestrator creates a new ConsistencyOrchestrator instance
func NewConsistencyOrchestrator() *ConsistencyOrchestrator {
	return &ConsistencyOrchestrator{
		ConsistencyLevel: 3, // Default to medium consistency
	}
}

// SetNetworkTelemetry updates the network metrics
func (co *ConsistencyOrchestrator) SetNetworkTelemetry(telemetry NetworkTelemetry) {
	co.NetworkMetrics = telemetry
	fmt.Printf("Updated network telemetry: Latency=%dms, PacketLoss=%.2f%%, Throughput=%.2fMbps\n",
		telemetry.Latency, telemetry.PacketLoss*100, telemetry.Throughput)
}

// AdjustConsistency adjusts the consistency level based on network conditions
func (co *ConsistencyOrchestrator) AdjustConsistency() {
	// Adjust consistency based on network conditions
	if co.NetworkMetrics.Latency > 200 || co.NetworkMetrics.PacketLoss > 0.2 {
		// Poor network conditions - reduce consistency
		co.ConsistencyLevel = 2
		fmt.Println("Reduced consistency level due to poor network conditions")
	} else if co.NetworkMetrics.Latency < 50 && co.NetworkMetrics.PacketLoss < 0.05 {
		// Good network conditions - increase consistency
		co.ConsistencyLevel = 4
		fmt.Println("Increased consistency level due to good network conditions")
	} else {
		// Normal network conditions - maintain medium consistency
		co.ConsistencyLevel = 3
		fmt.Println("Maintained medium consistency level")
	}
}
