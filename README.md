# unmarshal
将 环境变量、default、flag、自定义数据源 等绑定到结构体.

绑定类型支持：
* string
* bool
* int、int8、int16、int32、int64
* uint、uint8、uint16、uint32、uint64
* float32、float64
* slice
* map
* time.Duration
* 自定义接口/类型等

# Example
```go
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/recallsong/unmarshal"
	unmarshalflag "github.com/recallsong/unmarshal/unmarshal-flag"
	"github.com/spf13/pflag"
)

type subConfig struct {
	Name string `flag:"sub_name" default:"dingling" desc:"sub config name"`
	Age  int    `flag:"age" env:"AGE" default:"18" desc:"age"`
}

type config struct {
	Name          string `flag:"name" env:"NAME" default:"recallsong" desc:"name"`
	SubConfig     *subConfig
	Duration      time.Duration          `flag:"duration" env:"DURATION" default:"1ns" desc:"duration"`
	Map           map[string]interface{} `env:"MAP" default:"{\"age\":123}"`
	URLs          []string               `flag:"urls" env:"URLS" default:"http://localhost:80,http://localhost:81,http://localhost:82" desc:"urls"`
	Numbers       []int                  `flag:"numbers" env:"NUMBERS" default:"123,456,789" desc:"numbers"`
	CustomData    customData             `flag:"custom" env:"CUSTOM" default:"custom value" desc:"custom"`
	privideStruct subConfig
	privideString string
}

type customData struct {
	Text string
}

// 实现 encoding.TextUnmarshaler 接口
func (c *customData) UnmarshalText(text []byte) error {
	c.Text = string(text)
	return nil
}

// 实现 pflag.Value 接口
func (c *customData) Set(text string) error {
	c.Text = text
	return nil
}

func (c *customData) String() string {
	return c.Text
}

func (c *customData) Type() string {
	return "customData"
}

func main() {
	var cfg config
	err := unmarshal.BindDefault(&cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	err = unmarshal.BindEnv(&cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	flags := pflag.NewFlagSet("example", pflag.ExitOnError)
	err = unmarshalflag.BindFlag(flags, &cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	flags.Parse(os.Args[1:])
	fmt.Println(cfg)
}
```