package main

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "regexp"
    "sync"

    "github.com/alecthomas/kingpin"
)

var (
    start = kingpin.Arg("dir", "Directory to start processing in").String()
    concur = kingpin.Flag("concurrency", "Number of concurrent threads").Default("1").Short('c').Int()
    wg sync.WaitGroup
    lg sync.WaitGroup
    file_list chan string
)

func walkfunc(path string, _ os.FileInfo, _ error) error {
    file_list <-path

    return nil
}

func file_processor(){
    defer wg.Done()
    for i := range file_list{
        err := scanfile(i)
        if err != nil {
            panic(err)
        }
    }
}

func activate_processor(){
    wg.Add(1)
    go file_processor()
}

func scanfile(path string) error {
    //fmt.Printf("Scanning %s...\n", path)
    fp, err := os.Open(path)

    if err != nil {
        fmt.Println(err)
        //os.Exit(1)
    }

    defer fp.Close()

    reader := bufio.NewReader(fp)
    scanner := bufio.NewScanner(reader)

    for scanner.Scan() {
        current := scanner.Text()
        match, _ := regexp.MatchString("-{5}BEGIN [RD]SA PRIVATE KEY-{5}", current)
        if match {
            fmt.Printf("%s\t%s\n", path, current)
        }
        match, _ = regexp.MatchString("[pP]assword", current)
        if match {
            fmt.Printf("%s\t%s\n", path, current)
        }
    }

    return nil
}

func main(){
    kingpin.Parse()

    var starting string

    if len(*start) > 0 {
        starting = *start
    } else {
        cwd, err := os.Getwd()
        if err != nil {
            panic(err)
        }
        starting = cwd
    }
    file_list = make( chan string, *concur)

    for i :=0; i < *concur; i++ {
        activate_processor()
    }


    filepath.Walk(starting, walkfunc)
    close(file_list)
}
