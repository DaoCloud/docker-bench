package main

import (
    "fmt"
    "os"
    "time"

    "github.com/codegangsta/cli"
    "github.com/montanaflynn/stats"
)

func bench(requests int, n uint64, image string, args []string) {
    start := time.Now()
    timings := make([]float64, requests)

    for i := 0; i < requests; i++ {
        ts := time.Now()
        CreateService(image, args, n)
        timings[i] = time.Since(ts).Seconds()
        fmt.Printf("[%3.f%%] %d/%d request done\n", float64(i+1)/float64(requests)*100, i+1, requests)
    }

    total := time.Since(start)
    mean, _ := stats.Mean(timings)
    p90th, _ := stats.Percentile(timings, 90)
    p99th, _ := stats.Percentile(timings, 99)

    fmt.Printf("\n")
    fmt.Printf("Time taken for tests: %.3fs\n", total.Seconds())
    fmt.Printf("Time per request: %.3fs [mean] | %.3fs [90th] | %.3fs [99th]\n", mean, p90th, p99th)
}

func main() {
    app := cli.NewApp()
    app.Name = "scale-out"
    app.Usage = "DaoCloud swarm-kit benchmarking tool"
    app.Version = "0.1"
    app.Author = "haipeng"
    app.Email = "haipeng.wu@daocloud.io"
    app.Flags = []cli.Flag{
        cli.IntFlag{
            Name:  "requests, r",
            Value: 1,
            Usage: "Number of requests to create service.",
        },
        cli.IntFlag{
            Name:  "instances, n",
            Value: 1,
            Usage: "Number of instances to be created in each request.",
        },
        cli.StringSliceFlag{
            Name:  "image, i",
            Value: &cli.StringSlice{},
            Usage: "Image(s) to use for benchmarking.",
        },
    }

    app.Action = func(c *cli.Context) {
        if !c.IsSet("image") && !c.IsSet("i") {
            cli.ShowAppHelp(c)
            os.Exit(1)
        }
        bench(c.Int("requests"), uint64(c.Int("instances")), c.StringSlice("image")[0], c.Args())
    }

    app.Run(os.Args)
}
