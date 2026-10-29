// Package zonecfg đọc cấu hình khu vực từ Parameter Store, có cache để tránh gọi SSM mỗi event.
package zonecfg

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Loader struct {
	client   *ssm.Client
	path     string
	ttl      time.Duration
	mu       sync.Mutex
	cached   Values
	loadedAt time.Time
}

type Values map[string]string

func New(client *ssm.Client, path string) *Loader {
	return &Loader{client: client, path: path, ttl: 5 * time.Minute}
}

func (l *Loader) Get(ctx context.Context) (Values, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.cached != nil && time.Since(l.loadedAt) < l.ttl {
		return l.cached, nil
	}

	vals := Values{}
	p := ssm.NewGetParametersByPathPaginator(l.client, &ssm.GetParametersByPathInput{
		Path:           aws.String(l.path),
		WithDecryption: aws.Bool(true),
	})
	for p.HasMorePages() {
		page, err := p.NextPage(ctx)
		if err != nil {
			return nil, fmt.Errorf("load %s: %w", l.path, err)
		}
		for _, prm := range page.Parameters {
			vals[strings.TrimPrefix(aws.ToString(prm.Name), l.path+"/")] = aws.ToString(prm.Value)
		}
	}
	l.cached, l.loadedAt = vals, time.Now()
	return vals, nil
}

func (v Values) Float(key string) (float64, error) {
	s, ok := v[key]
	if !ok {
		return 0, fmt.Errorf("missing parameter %q", key)
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parameter %q: %w", key, err)
	}
	return f, nil
}
