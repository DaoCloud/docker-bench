package main

import (
    "fmt"
    "time"

    "golang.org/x/net/context"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/swarm"
    "github.com/docker/docker/client"
)

func NewClient() *client.Client {
    c, err := client.NewClient("unix:///var/run/docker.sock", "v1.24", nil, map[string]string{"User-Agent": "engine-api-cli-1.0"})
    if err != nil {
        panic(err)
    }

    return c
}

func CreateService(image string, args []string, n uint64) string {
    var service swarm.ServiceSpec
    var options types.ServiceCreateOptions
    service.TaskTemplate.ContainerSpec.Image = image
    service.TaskTemplate.ContainerSpec.Args = args
    service.Mode.Replicated = &swarm.ReplicatedService{Replicas: &n}

    c := NewClient()
    res, err := c.ServiceCreate(context.Background(), service, options)
    if err != nil {
        panic(err)
    }

    WaitServiceInstanceN(res.ID, n, 500, 20*int(n))
    return res.ID
}

func WaitServiceInstanceN(serviceID string, instanceNum uint64, interval int, times int) {
    ok := false
    for i := 0; i < times; i++ {
        var ops types.TaskListOptions
        c := NewClient()
        tasks, err := c.TaskList(context.Background(), ops)
        if err != nil {
            panic(err)
        }

        var running uint64 = 0
        for _, task := range tasks {
            if task.ServiceID == serviceID && task.Status.State == "running" {
                running += 1
            }
        }

        if running == instanceNum {
            ok = true
            break
        }

        time.Sleep(time.Duration(interval)*time.Millisecond)
    }

    if !ok {
        panic(fmt.Errorf("service %s can not reach %d instance in %dms", serviceID, instanceNum, interval*times))
    }
}

func Scale(serviceID string, n uint64) {
    c := NewClient()
    service, _, err := c.ServiceInspectWithRaw(context.Background(), serviceID)
    if err != nil {
        panic(err)
    }

    serviceMode := &service.Spec.Mode
    if serviceMode.Replicated == nil {
        panic(fmt.Errorf("scale can only be used with replicated mode"))
    }

    var ori uint64 = 0
    if serviceMode.Replicated.Replicas != nil {
        ori = *serviceMode.Replicated.Replicas
    }

    if ori == n {
        return
    }

    var diff uint64 = 0
    if ori > n {
        diff = ori - n
    } else {
        diff = n - ori
    }

    serviceMode.Replicated.Replicas = &n
    c = NewClient()
    err = c.ServiceUpdate(context.Background(), service.ID, service.Version, service.Spec, types.ServiceUpdateOptions{})
    if err != nil {
        panic(err)
    }

    WaitServiceInstanceN(serviceID, n, 500, 20*int(diff))
}
