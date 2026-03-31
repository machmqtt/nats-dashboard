package collector

import "testing"

func TestBuildTopologyRoutes(t *testing.T) {
	snap := &Snapshot{
		Varz: map[string]*Varz{
			"A": {ServerName: "nats-1", Connections: 10, Cluster: ClusterOptsVarz{Name: "dc1"}},
			"B": {ServerName: "nats-2", Connections: 5, Cluster: ClusterOptsVarz{Name: "dc1"}},
		},
		Routez: map[string]*Routez{
			"A": {Routes: []RouteInfo{{RemoteID: "B", RemoteName: "nats-2"}}},
			"B": {Routes: []RouteInfo{{RemoteID: "A", RemoteName: "nats-1"}}},
		},
		Health: map[string]*HealthStatus{
			"A": {Status: "ok"},
			"B": {Status: "ok"},
		},
		Rates: map[string]*ServerRates{},
	}

	g := buildTopology(snap, nil)
	if len(g.Nodes) != 2 {
		t.Errorf("nodes = %d, want 2", len(g.Nodes))
	}
	// Routes are bidirectional, should be deduplicated to 1.
	if len(g.Links) != 1 {
		t.Errorf("links = %d, want 1 (deduplicated)", len(g.Links))
	}
	if g.Links[0].Type != "route" {
		t.Errorf("link type = %q, want route", g.Links[0].Type)
	}
}

func TestBuildTopologyGateways(t *testing.T) {
	snap := &Snapshot{
		Varz: map[string]*Varz{
			"A": {ServerName: "nats-1"},
		},
		Gatewayz: map[string]*Gatewayz{
			"A": {
				OutboundGateways: map[string]*RemoteGatewayz{
					"dc2": {IsConfigured: true},
				},
			},
		},
		Health: map[string]*HealthStatus{"A": {Status: "ok"}},
		Rates:  map[string]*ServerRates{},
		Routez: map[string]*Routez{},
	}

	g := buildTopology(snap, nil)
	if len(g.Nodes) != 2 {
		t.Errorf("nodes = %d, want 2 (server + gateway)", len(g.Nodes))
	}
	if len(g.Links) != 1 {
		t.Errorf("links = %d, want 1", len(g.Links))
	}
}
