package main

import (
    "fmt"
    "os"
    "sync"
    "time"

    "github.com/codegangsta/cli"
    "github.com/montanaflynn/stats"
)

func worker(requests int, image string, args []string, completeCh chan time.Duration) {
    for i := 0; i < requests; i++ {
        start := time.Now()
        RunContainer(image, args)
        completeCh <- time.Since(start)
    }
}

func session(requests, concurrency int, images []string, args []string, completeCh chan time.Duration) {
    var wg sync.WaitGroup
    var size = len(images)

    for i := 0; i < concurrency; i++ {
        wg.Add(1)
        image := images[i%size]
        go func() {
            worker(requests, image, args, completeCh)
            wg.Done()
        }()
    }
    wg.Wait()
}

func bench(requests, concurrency int, images []string, args []string) {
    start := time.Now()

    timings := make([]float64, 0)
    completeCh := make(chan time.Duration, requests*concurrency)
    doneCh := make(chan struct{})
    current := 0
    go func() {
        for timing := range completeCh {
            timings = append(timings, timing.Seconds())
            current++
            percent := float64(current) / float64(requests*concurrency) * 100
            fmt.Printf("[%3.f%%] %d/%d containers started\n", percent, current, requests*concurrency)
        }
        doneCh <- struct{}{}
    }()
    session(requests, concurrency, images, args, completeCh)
    close(completeCh)
    <-doneCh

    total := time.Since(start)
    mean, _ := stats.Mean(timings)
    p90th, _ := stats.Percentile(timings, 90)
    p99th, _ := stats.Percentile(timings, 99)

    fmt.Printf("\n")
    fmt.Printf("Time taken for tests: %.3fs\n", total.Seconds())
    fmt.Printf("Time per container: %.3fs [mean] | %.3fs [90th] | %.3fs [99th]\n", mean, p90th, p99th)
}

func main() {
    app := cli.NewApp()
    app.Name = "create-container"
    app.Usage = "DaoCloud swarm-kit benchmarking tool"
    app.Version = "0.1"
    app.Author = "haipeng"
    app.Email = "haipeng.wu@daocloud.io"
    app.Flags = []cli.Flag{
        cli.IntFlag{
            Name:  "concurrency, c",
            Value: 1,
            Usage: "Number of multiple requests to perform at a time.",
        },
        cli.IntFlag{
            Name:  "requests, r",
            Value: 1,
            Usage: "Number of containers to start in each thred.",
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
        bench(c.Int("requests"), c.Int("concurrency"), c.StringSlice("image"), c.Args())
    }

    app.Run(os.Args)
}
