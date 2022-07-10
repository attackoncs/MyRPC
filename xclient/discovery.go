package xclient

import (
	"errors"
	"math"
	"math/rand"
	"sync"
	"time"
)

type SelectMode int

const (
	RandomSelect     SelectMode = iota //select randomly
	RoundRobinSelect                   //robin algorithm
)

type Discovery interface {
	Refresh() error //refresh from remote registry
	Update(servers []string) error
	Get(mode SelectMode) (string, error)
	GetAll() ([]string, error)
}

var _ Discovery = (*MultiServersDiscovery)(nil)

func (d MultiServersDiscovery) Refresh() error {
	//TODO implement me
	return nil
}

func (d MultiServersDiscovery) Update(servers []string) error {
	//TODO implement me
	d.mu.Lock()
	defer d.mu.Unlock()
	d.servers = servers
	return nil
}

func (d MultiServersDiscovery) Get(mode SelectMode) (string, error) {
	//TODO implement me
	d.mu.Lock()
	defer d.mu.Unlock()
	n := len(d.servers)
	if n == 0 {
		return "", errors.New("rpc discovery: no available servers")
	}

	switch mode {
	case RandomSelect:
		return d.servers[d.r.Intn(n)], nil
	case RoundRobinSelect:
		s := d.servers[d.index%n] //servers could be updated so mode n to ensure safety
		d.index = (d.index + 1) % n
		return s, nil

	default:
		return "", errors.New("rpc discovery: not supported select mode")
	}
}

func (d MultiServersDiscovery) GetAll() ([]string, error) {
	//TODO implement me
	d.mu.RLock()
	defer d.mu.RUnlock()

	servers := make([]string, len(d.servers), len(d.servers))
	copy(servers, d.servers)
	return servers, nil
}

//不需注册中心，服务列表手工维护的服务发现的结构体:MultiServersDiscovery
//MultiServersDiscovery is a discovery for multi servers without a registry center
// user provides the server addresses explicitly instead
type MultiServersDiscovery struct {
	r       *rand.Rand   //生成随机数
	mu      sync.RWMutex //保护
	servers []string
	index   int //记录robin算法选择的位置
}

func NewMultiServerDiscovery(servers []string) *MultiServersDiscovery {
	d := &MultiServersDiscovery{
		servers: servers,
		//r 是一个产生随机数的实例，初始化时使用时间戳设定随机数种子，避免每次产生相同的随机数序列
		r: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	d.index = d.r.Intn(math.MaxInt32 - 1) //为了避免每次从 0 开始，初始化时随机设定一个值。
	return d
}
