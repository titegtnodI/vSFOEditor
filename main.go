//vSFOEditor by titegtnodI

package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

const (
    VERSION = "V1.1.0"
    vitainput = "PSPBUILD"
    psinput = "PSBUILD"
    output = "PARAM.SFO"
    inplen = 4912
    details = 272
    title = 1296
    name = 4656
    title2 = 4784
    DEBUG = false
    defaultsCfg = "defaults.cfg"
)

func nullPad (what string) (string) {
    var out string
    out = what
    for i := len(out);i < 20;i++ {
        out = out + "\x00"
    }
    return out[0:20]
}

func readFile (where string, length int) (data []byte, err error) {
    fmt.Printf("Opening \"%s\" ...\n", where)
    file, err := os.Open(where)
    defer file.Close()
    if err != nil {
        return
    }

    fmt.Printf("Reading \"%s\" ...\n", where)
    data = make([]byte, length)
    count, err := file.Read(data)
    if err != nil {
        return
    }

    data = data[:count]

    return
}

func main () {
    stdin := bufio.NewReader(os.Stdin)

    defer func () {
            fmt.Printf("Press enter to exit.")
            var garbage string
            fmt.Fscanf(stdin, "%s", &garbage)
    } ()

    fmt.Printf("vSFOEditor %s by titegtnodI\n\n", VERSION)

    var sDefaults []string
    var dTitle, dName, dDetails string
    defaults, err := readFile(defaultsCfg, 62)
    if err == nil {
        sDefaults = strings.Split(string(defaults), "\n")
        if len(sDefaults) >= 3 {
            dTitle = strings.Trim(sDefaults[0], "\r")
            dName = strings.Trim(sDefaults[1], "\r")
            dDetails = strings.Trim(sDefaults[2], "\r")
        }
    }
    if err != nil || len(sDefaults) < 3 {
        fmt.Printf("Failed to read from %s.\n", defaultsCfg)
        dTitle = ""
        dName = "VHBL READY SAVE"
        dDetails = "Made with vSFOEditor"
    }
    fmt.Printf("\n")

    var choice string
    fmt.Printf("Is this for PSP(1) or Playstation(2)? [1]: ")
    choice, _ = stdin.ReadString('\n')
    choice = strings.Trim(choice, "\r\n")
    if choice != "1" && choice != "2" {
        choice = "1"
    }
    fmt.Printf("\n")

    var data []byte
    if choice == "1" {
        data, err = readFile(vitainput, inplen)
    } else {
        data, err = readFile(psinput, inplen)
    }
    if err != nil {
        fmt.Printf("An error has occurred.\n")
        panic(err)
    }
    fmt.Printf("\n")

    if DEBUG {
        if choice == "1" {
            fmt.Printf("Title: %s\nName: %s\nDetails: %s\n\n",
                       data[title:title+20], data[name:name+20], data[details:details+20])
        } else {
            fmt.Printf("Title: %s\nName: %s\n\n",
                       data[title:title+20], data[title2:title2+20])
        }
    }

    var nTitle, nName, nDetails string
    if choice == "1" {
        fmt.Printf("App Name (Same as folder) [%s]: ", dTitle)
    } else {
        fmt.Printf("App Name (Same as folder, in format \"SXXX#####\") [%s]: ", dTitle)
    }
    nTitle, _ = stdin.ReadString('\n')
    nTitle = strings.Trim(nTitle, "\r\n")
    nTitle = nullPad(nTitle)

    fmt.Printf("\"Name\" [%s]: ", dName)
    nName, _ = stdin.ReadString('\n')
    nName = strings.Trim(nName, "\r\n")
    if nName == "" {
        nName = dName
    }
    nName = nullPad(nName)

    if choice == "1" {
        fmt.Printf("Details [%s]: ", dDetails)
        nDetails, _ = stdin.ReadString('\n')
        nDetails = strings.Trim(nDetails, "\r\n")
        if nDetails == "" {
            nDetails = dDetails
        }
        nDetails = nullPad(nDetails)
    }

    if DEBUG {
        if choice == "1" {
            fmt.Printf("\nnTitle: %s\nnName: %s\nnDetails: %s\n", nTitle, nName,
                       nDetails)
        } else {
            fmt.Printf("\nnTitle: %s\nnName: %s\n", nTitle, nName)
        }
    }

    fmt.Printf("\nOpening \"%s\" ...\n", output)
    file, err := os.Create(output)
    defer file.Close()
    if err != nil {
        fmt.Printf("An error has occurred.\n")
        panic(err)
    }

    fmt.Printf("Writing \"%s\" ...\n\n", output)
    //TODO Don't use formatted strings here, it's stilly
    if choice == "1" {
        data = []byte(fmt.Sprintf("%s%s%s%s%s%s%s%s%s", data[:details], nDetails,
                              data[details+20:title], nTitle, data[title+20:name], nName,
                              data[name+20:title2], nTitle, data[title2+20:]))
    } else {
        data = []byte(fmt.Sprintf("%s%s%s%s%s", data[:title], nTitle,
                              data[title+20:title2], nName, data[title2+20:]))
    }
    _, err = file.Write(data)
    if err != nil {
        fmt.Printf("An error has occurred.\n")
        panic(err)
    }

    fmt.Printf("Done! No errors were detected.\n")
}
