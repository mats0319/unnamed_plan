package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
)

type generator struct {
	reEngineIns *reEngines
}

var generatorIns = &generator{}

func main() {
	log.Println("> start  generate ts file(s).")
	defer log.Println("> finish generate ts file(s).")

	entry, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalln("read dir failed, error: ", err)
	}

	for i := range entry {
		if entry[i].IsDir() {
			continue
		}

		var fileInfo fs.FileInfo
		fileInfo, err = entry[i].Info()
		if err != nil {
			log.Println("get file info failed, error: ", err)
			continue
		}

		if !strings.HasSuffix(fileInfo.Name(), protoFileSuffix) {
			continue
		}

		pbFileIns := generatorIns.parse(fileInfo.Name())
		content := pbFileIns.generateTSContent()

		err = os.WriteFile(generateOutputFileName(pbFileIns.name), content, 0777)
		if err != nil {
			log.Fatalln("write file failed, error: ", err)
		}
	}
}

func generateOutputFileName(fileName string) string {
	return fmt.Sprintf("%s%s%s",
		outputDir,
		strings.TrimSuffix(fileName, protoFileSuffix),
		defaultOutputFileExtension)
}
