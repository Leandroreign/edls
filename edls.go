package main

import (
	"time"
)

// file types
const (
	fileRegular int = iota
	fileDirectory
	fileExecutable
	fileCompress
	fileImage
	fileLink
)

// file extension
const (
	exe = ".exe"
	deb = ".deb"
	zip = ".zip"
	gz  = ".gz"
	tar = ".tar"
	rar = ".rar"
	png = ".png"
	jpg = ".jpg"
	gif = ".gof"
)

type file struct {
	name            string
	fileType        int
	isDir           bool
	isHidden        bool
	userName        string
	groupName       string
	size            uint64
	modificatioTime time.Time
	mode            string
}

type styleFileType struct {
	icon   string
	color  string
	symbol string
}

var mapStyleByFileType = map[int]styleFileType{
	fileRegular:    {icon: "📃"},
	fileDirectory:  {icon: "📂", color: "Orange", symbol: "/"},
	fileExecutable: {icon: "🎽", color: "Red", symbol: "*"},
	fileCompress:   {icon: "🗃️", color: "Yellow"},
	fileImage:      {icon: "📸", color: "Purple"},
	fileLink:       {icon: "🔗", color: "Blue"},
}
