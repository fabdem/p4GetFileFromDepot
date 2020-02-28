package main

//	F.Demurger 2019-04
//  	args:
//			p4GetFileFromDepot <version> <file path/name in depot> <local path>
//
//      Option -v version
//      Option -u <user>
//
//      Get the file/version from the depot and store it to the local path.
//
//     	Returns <local  path>filename#<version>
//
//			P4 cli needs to be installed and in the path.
//
//
//	cross compilation AMD64:  env GOOS=windows GOARCH=amd64 go build p4getfilefromdepot.go

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {

	var versionFlg bool
	var user string
	var p4Cmd string // p4 command path
	var err error

	const usageVersion   = "Display Version"
	const usageUser      = "specify a username"

  // Have to create a specific set, the default one is poluted by some test stuff from another lib (?!)
  checkFlags := flag.NewFlagSet("check", flag.ExitOnError)

	checkFlags.BoolVar(&versionFlg, "version", false, usageVersion)
	checkFlags.BoolVar(&versionFlg, "v", false, usageVersion + " (shorthand)")
	checkFlags.StringVar(&user, "user", "", usageUser)
	checkFlags.StringVar(&user, "u", "", usageUser + " (shorthand)")
	checkFlags.Usage = func() {
        fmt.Printf("Usage: %s <opt> <version> <file path/name in depot> <localpath>\n",os.Args[0])
				fmt.Print("Get the file/version from the depot and store it to the local path.")
        fmt.Print("Returns <local  path>filename#<version> (P4 file naming convention).")
        checkFlags.PrintDefaults()
    }

    // Check parameters
	checkFlags.Parse(os.Args[1:])

	if versionFlg {
		fmt.Printf("Version %s\n", "2020-02  v1.0.5")
		os.Exit(0)
	}

	// Check presence of p4 cli
	if p4Cmd, err = exec.LookPath("p4"); err != nil {
		fmt.Printf("P4 command line is not installed - %s\n", err)
		os.Exit(1)
	}

  // Parse the command parameters
  index     := len(os.Args)
	version   := os.Args[index - 3]
	depotFile := os.Args[index - 2] + "#" + version
	localPath := os.Args[index - 1]

	// extract filename
	fileName  := filepath.Base(os.Args[index - 2])
	// fmt.Printf("fileName: %s\n", fileName)
	ext := filepath.Ext(fileName)
	localFileName := localPath + fileName[0:len(fileName)-len(ext)] + "#" + version + ext

	/*
	fmt.Printf("version: %s\n", version)
	fmt.Printf("depotFile: %s\n", depotFile)
	fmt.Printf("localPath: %s\n", localPath)
	fmt.Printf("localFileName: %s\n", localFileName)
	*/
	/*
	P4 print options :
	-k suppress RCS keyword expansion (p4 variables)
	-q suppress header printing
	-o <local file>  redirect to a file -> permissions are maintained so warning. REQUIRED! Output to std forces it to utf8,

	e.g. p4 -u myusername print -k -q -o ./test.cmd  //ap4rootproject/dev/folder/alocfile_bulgarian.txt#7
	*/
	var out []byte

	if len(user) > 0 {
		out, err = exec.Command(p4Cmd, "-u", user, "print","-k", "-q", "-o",localFileName, depotFile).CombinedOutput()
		// fmt.Printf("P4 command line result - %s\n %s\n", err, out)
	} else {
		out, err = exec.Command(p4Cmd, "print","-k", "-q", "-o",localFileName, depotFile).CombinedOutput()
	}
	if err != nil {
		fmt.Printf("P4 command line error\n%s\n%s\n", err, out)
		os.Exit(2)
	}

	// Unfortunately p4 print status in linux is not reliable.
	// err != nil when syntax err but not if file doesn't exist.
	// So manually check if a file was created:
	if _, err := os.Stat(localFileName); err != nil {
	   if os.IsNotExist(err) {
	       // file does not exist
				fmt.Printf("Error - No file produced\n%s\n%s\n", err, out)
				os.Exit(3)
	   } else {
	       // Can't get file stat
				fmt.Printf("Error - can't access the status of file produced\n%s\n%s\n", err, out)
				os.Exit(4)
		}
	}

	// Build return string:
	fmt.Printf("%s\n",localFileName)
}
