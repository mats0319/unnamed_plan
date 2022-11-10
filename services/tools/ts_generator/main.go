package main

import (
    "github.com/pkg/errors"
    "log"
    "os"
    "strings"
)

func main() {
    path, err := os.Executable() // path, absolute path include exec protocolBuffer name
    if err != nil {
        log.Fatalln(errors.Wrap(err, "get path failed"))
    }

    index := strings.LastIndex(path, "/")
    dir := path[:index] + "/proto"

    entry, err := os.ReadDir(dir)
    if err != nil {
        log.Fatalln(errors.Wrap(err, "read path failed"))
    }

    for i := range entry {
        if entry[i].IsDir() {
            continue
        }

        fileInfo, err := entry[i].Info()
        if err != nil {
            log.Println("get protocolBuffer info failed", err)
            continue
        }

        if strings.HasSuffix(fileInfo.Name(), ".proto") {
            pbFileIns := parsePBFile(dir+"/"+fileInfo.Name(), strings.TrimSuffix(fileInfo.Name(), ".proto"))
            generateTS(pbFileIns)
        }
    }
}
