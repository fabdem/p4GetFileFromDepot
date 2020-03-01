package main

//	F.Demurger 2019-04
//  	args:
//			p4GetFileFromDepot <version> <file path/name in depot> <local path only>
//
//      Option -v version
//      Option -u <user>
//			Option -r <revision number>
//
//      Get the file/version from the depot and store it to the local path.
//			If no revision specified, returns the head.
//
//     	Returns <local  path>filename#<version>
//
//			P4 cli needs to be installed and in the path.
//
//
//	cross compilation AMD64:  env GOOS=windows GOARCH=amd64 go build p4getfilefromdepot.go

import (
	"errors"
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
	const usageUser      = "Specify a username"
	const usageRev       = "Specify a revision number"

  // Have to create a specific set, the default one is poluted by some test stuff from another lib (?!)
  checkFlags := flag.NewFlagSet("check", flag.ExitOnError)

	checkFlags.BoolVar(&versionFlg, "version", false, usageVersion)
	checkFlags.BoolVar(&versionFlg, "v", false, usageVersion + " (shorthand)")
	checkFlags.StringVar(&user, "user", "", usageUser)
	checkFlags.StringVar(&user, "u", "", usageUser + " (shorthand)")
	checkFlags.IntVar(&rev, "revision", 0, usageRev)
	checkFlags.IntVar(&rev, "r", 0, usagerev + " (shorthand)")
	checkFlags.Usage = func() {
        fmt.Printf("Usage: %s [opt] <file path/name in depot> <localpath>\n",os.Args[0])
				fmt.Print("Get the head file from the depot and store it to the local path.")
        fmt.Print("Returns <local  path>filename#<revision> (P4 file naming convention).")
				fmt.Print("Use option -r to specify a revision number, if not the head rev is downloaded.")
        checkFlags.PrintDefaults()
    }

    // Check parameters
	checkFlags.Parse(os.Args[1:])

	if versionFlg {
		fmt.Printf("Version %s\n", "2020-02  v1.1.0")
		os.Exit(0)
	}

	// Check presence of p4 cli
	if p4Cmd, err = exec.LookPath("p4"); err != nil {
		fmt.Printf("P4 command line is not installed - %s\n", err)
		os.Exit(1)
	}

  // Parse the command parameters
  index     := len(os.Args)
//	version   := os.Args[index - 3]
//	depotFile := os.Args[index - 2] + "#" + version
	depotFile := os.Args[index - 2]
	localPath := os.Args[index - 1]

	fileName  := filepath.Base(depotFile) // extract filename
	ext := filepath.Ext(depotFile) 				// Read extension

	if rev > 0 {  // If a specific version has been requested through -r
		localFileName := localPath + fileName[0:len(fileName)-len(ext)] + "#" + rev + ext

	} else { // Get head rev
		rev, err := p4GetHeadRev(depotFile,user)
		if err != nil {
			fmt.Printf("P4 command line error - %s\n", err)
			os.Exit(1)
		}
		localFileName := localPath + fileName[0:len(fileName)-len(ext)] + "#" + rev + ext
	}


	// fmt.Printf("fileName: %s\n", fileName)

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
		os.Exit(1)
	}

	// Unfortunately p4 print status in linux is not reliable.
	// err != nil when syntax err but not if file doesn't exist.
	// So manually check if a file was created:
	if _, err := os.Stat(localFileName); err != nil {
	   if os.IsNotExist(err) {
	       // file does not exist
				fmt.Printf("Error - No file produced\n%s\n%s\n", err, out)
				os.Exit(1)
	   } else {
	       // Can't get file stat
				fmt.Printf("Error - can't access the status of file produced\n%s\n%s\n", err, out)
				os.Exit(1)
		}
	}

	// Build return string:
	fmt.Printf("%s\n",localFileName)
}

// Get a file from P4 depot and store it under the path/name provided
func p4GetFile(depotFile string, localFileName string, user string) error {

	var out []byte

	if len(user) > 0 {
		out, err = exec.Command(p4Cmd, "-u", user, "print","-k", "-q", "-o",localFileName, depotFile).CombinedOutput()
		// fmt.Printf("P4 command line result - %s\n %s\n", err, out)
	} else {
		out, err = exec.Command(p4Cmd, "print","-k", "-q", "-o",localFileName, depotFile).CombinedOutput()
	}
	if err != nil {
		fmt.Printf("P4 command line error\n%s\n%s\n", err, out)
		return err
	}

	// Unfortunately p4 print status in linux is not reliable.
	// err != nil when syntax err but not if file doesn't exist.
	// So manually check if a file was created:
	if _, err := os.Stat(localFileName); err != nil {
	 	if os.IsNotExist(err) {
		 	// file does not exist
			fmt.Printf("Error - No file produced\n%s\n%s\n", err, out)
			return err
	 	} else {
		 	// Can't get file stat
			fmt.Printf("Error - can't access the status of file produced\n%s\n%s\n", err, out)
			return err
		}
	}
	return nil
}

// Get from P4 the head revision number of a file
func p4GetHeadRev(depotFileName string, user string) (rev int, error)

	var out []byte
	if len(user) > 0 {
		// fmt.Printf(p4Cmd + " -u " + user + " files " + " " + depotFile + "\n")
		out, err = exec.Command(p4Cmd, "-u", user, "files", depotFile).Output()
	} else {
		// fmt.Printf(p4Cmd + " files " + depotFile + "\n")
		out, err = exec.Command(p4Cmd, "files", depotFile).Output()
	}
	if err != nil {
		fmt.Printf("P4 command line error - %s\n", err)
		return 0,err
	}

	// Read version
	// e.g. //Project/dev/localization/afile_bulgarian.txt#8 - edit change 4924099 (utf16)
	idxBeg := strings.LastIndex(string(out),"#") + len("#")
	idxEnd := strings.LastIndex(string(out)," - ")
	// Check response to prevent out of bound index
	if idxBeg == -1 || idxEnd == -1 || idxBeg >= idxEnd {
		fmt.Printf("Format error in P4 response: %s\n", string(out))
		return 0,err
	}
	// sRev := string(out[strings.LastIndex(string(out),"#") + len("#"):strings.LastIndex(string(out)," - ")])
	sRev := string(out[idxBeg:idxEnd])

	rev, err := strconv.Atoi(sRev) // Check format
	if err != nil {
		fmt.Printf("sRev=%s\n", sRev)
		fmt.Printf("Format err=%s\n", err)
		return 0,err
	}

	return rev, nil
}
