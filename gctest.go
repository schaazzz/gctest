package main

import (
    "fmt"
    "time"
    _ "net/http/pprof"
    "net/http"
    "os"
    "log"
    "io"
)

type Pipe struct {
    Data      chan * DataChunk
}

type DataChunk struct {
    Length int
    Value string
}

func main() {
    runFile, _ := os.OpenFile("run.log", os.O_TRUNC | os.O_CREATE | os.O_WRONLY, 0666)

    defer func() {
        runFile.Close()
    } ()

    tee := io.MultiWriter(runFile, os.Stdout)
    runLogger := log.New(tee, "", 0)

    pipe0 := &Pipe {
        Data: make(chan * DataChunk),
    }

    pipe1 := &Pipe {
        Data: make(chan * DataChunk),
    }

    time.Sleep(5 * time.Second)

    go func() {
        fmt.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    go func() {
        for {
            time.Sleep(100 * time.Millisecond)
            data := "Hello from Pipe #1"
            chunk := &DataChunk{Length: len(data), Value: data}
            runLogger.Println("\r\n----------------------------")
            runLogger.Println("Init Pipe #1", chunk, &chunk)
            pipe1.Data <- chunk        
        }
    } ()

    go func() {
        for {
            select {
            case data := <- pipe1.Data:
                runLogger.Println("Pipe #1 sending", data, &data)
                pipe0.Data <- data
            }
        }
    } ()

    for {
        select {
        case data := <- pipe0.Data:
            runLogger.Println("Pipe #0 received", data, &data)
            runLogger.Println("----------------------------")
        }
    }
}
