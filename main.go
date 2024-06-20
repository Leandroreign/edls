package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"

	"golang.org/x/exp/constraints"
)

const Windows = "windows"

func main() {
	// el flag h ya esta siendo utilizado como help para ayuda
	// filter flags
	flagPattern := flag.String("p", "", "Filter files by pattern \n pattern might be '.png, .img'")
	//flagAll := flag.Bool("a", false, "Return all the files including hidden ones")
	flagNumeberRecord := flag.Int("n", 0, "Number of records to be returned")

	flagOrderByTime := flag.Bool("t", false, "Return records sorted by modification time")
	flagOrderBySize := flag.Bool("s", false, "Return records sorted by size")
	flagReverse := flag.Bool("r", false, "Reverse the order of the records")

	// mapea cada uno de los flags en variables para que sean accesibles
	// flags de salida
	flag.Parse()
	// para ver el puntero ya que flag devuelve un puntero
	//fmt.Println(flagPattern)
	// para ver el valor del puntero
	/*
		fmt.Println(*flagPattern)
		fmt.Println(*flagAll)
		fmt.Println(*flagNumeberRecord)
		fmt.Println(*flagOrderByTime)
		fmt.Println(*flagOrderBySize)
		fmt.Println(*flagReverse)*/

	// la funcion args devuelve los argumentos cuando son mas de uno
	// la funcion arg(0) devuelve un unico argumento
	path := flag.Arg(0)

	if path == "" {
		path = "."
	}

	dirs, err := os.ReadDir(path)

	if err != nil {
		panic(err)
	}

	fs := []file{}

	for _, dir := range dirs {
		if *flagPattern != "" {
			isMatched, err := regexp.MatchString("(?i)"+*flagPattern, dir.Name())
			if err != nil {
				panic(err)
			}

			if !isMatched {
				continue
			}
		}

		f, err := getFile(dir, false)
		if err != nil {
			panic(err)
		}
		// esto "(?I)" es para decirle que no sea case Sensitive

		fs = append(fs, f)
	}
	if !*flagOrderBySize && !*flagOrderByTime {
		orderByName(fs, *flagReverse)
	} else if *flagOrderBySize && *flagOrderByTime {
		orderBySize(fs, *flagReverse)
	} else if *flagOrderByTime {
		orderByTime(fs, *flagReverse)
	} else if *flagOrderBySize {
		orderBySize(fs, *flagReverse)
	}

	if *flagNumeberRecord == 0 || *flagNumeberRecord > len(fs) {
		*flagNumeberRecord = len(fs)
	}

	printFileList(fs, *flagNumeberRecord)

}

func sortList[T constraints.Ordered](i, j T, isReversed bool) bool {
	if isReversed {
		return i > j
	}
	return i < j
}

func orderByName(files []file, isReversed bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return sortList[string](strings.ToLower(files[i].name), strings.ToLower(files[j].name), isReversed)
	})
}

func orderBySize(files []file, isReversed bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return sortList[uint64](files[i].size, files[j].size, isReversed)
	})
}

func orderByTime(files []file, isReversed bool) {
	sort.SliceStable(files, func(i, j int) bool {
		return sortList[uint64](uint64(files[i].modificatioTime.Unix()), uint64(files[j].modificatioTime.Unix()), isReversed)
	})
}

func printFileList(fs []file, nRecords int) {
	fmt.Println("")
	for _, file := range fs[:nRecords] {
		style := mapStyleByFileType[file.fileType]

		fmt.Printf("%s %s %s %10d %s %s %s %s\n",
			file.mode,
			file.userName,
			file.groupName,
			file.size,
			file.modificatioTime.Format(time.DateTime),
			style.icon,
			file.name,
			style.symbol,
		)
	}
}

func getFile(dir fs.DirEntry, isHidden bool) (file, error) {
	// fmt.Println("dir", dir)
	info, err := dir.Info()
	if err != nil {
		return file{}, fmt.Errorf("dir.Info():  %v", err)
	}

	f := file{
		name:            dir.Name(),
		isDir:           dir.IsDir(),
		isHidden:        isHidden,
		userName:        "",
		groupName:       "",
		size:            uint64(info.Size()),
		modificatioTime: info.ModTime(),
		mode:            info.Mode().String(),
	}
	setFile(&f)

	return f, nil
}

func setFile(f *file) {
	switch {
	case isLink(*f):
		f.fileType = fileLink
	case isCompress(*f):
		f.fileType = fileCompress
	case isImage(*f):
		f.fileType = fileImage
	case isExecute(*f):
		f.fileType = fileExecutable
	case f.isDir:
		f.fileType = fileDirectory
	default:
		f.fileType = fileRegular
	}

}

func isLink(f file) bool {
	return strings.HasPrefix(strings.ToUpper(f.mode), "L")
}

func isExecute(f file) bool {
	if runtime.GOOS == Windows {
		return strings.HasSuffix(f.name, exe)
	}
	// en unix para saber si es un ejecutable el mode contiene permisos de ejecucion

	return strings.Contains(f.mode, "x")
}

func isCompress(f file) bool {
	return strings.HasSuffix(f.name, zip) ||
		strings.HasSuffix(f.name, gz) ||
		strings.HasSuffix(f.name, rar) ||
		strings.HasSuffix(f.name, tar) ||
		strings.HasSuffix(f.name, deb)
}

func isImage(f file) bool {
	return strings.HasSuffix(f.name, png) ||
		strings.HasSuffix(f.name, jpg) ||
		strings.HasSuffix(f.name, gif)
}

func print[T any](msgs ...T) {
	for _, v := range msgs {
		fmt.Printf("%v \n", v)
	}
}
