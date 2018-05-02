package dispatcher

import "github.com/kobtea/iapetus/pkg/config"

func NewDispatcher(cluster config.Cluster) *Dispatcher {
	return &Dispatcher{
		Cluster: cluster,
	}
}

type Dispatcher struct {
	Cluster config.Cluster
}

func (d Dispatcher) FindNode() *config.Node {
	return d.defaultNode()
}

func (d Dispatcher) defaultNode() *config.Node {
	for _, r := range d.Cluster.Rules {
		if r.Default {
			for _, n := range d.Cluster.Nodes {
				if n.Name == r.Target {
					return &n
				}
			}
		}
	}
	if len(d.Cluster.Nodes) > 0 {
		return &d.Cluster.Nodes[0]
	}
	return nil
}
