package main

import (
    //"github.com/hoisie/web.go"
    "github.com/sschober/web.go"
    "fmt"
    "strconv"
    //"strings"
    "json"
    "time"
    "io/ioutil"
    "container/vector"
    "log"
)

// Type to represent a roar.
//
// Note:
//    json.Unmarshal() can only unmarshal public fields, so we are
//    not able to make them private.
type Roar struct {
    Author       string
    Text         string
    CreationDate string
}

// Constructor to create a new Roar
//
// Sets CreationDate.
func NewRoar(author string, text string) *Roar {
    return &Roar{Author: author, Text: text,
        CreationDate: time.LocalTime().Format(time.RFC1123)}
}

// Custom string representation
func (r Roar) String() string {
    return fmt.Sprintf("%s - %s\n%s", r.Author, r.CreationDate, r.Text)
}

// Marshal Roar object to Json string representation
func (r Roar) toJson() string {
    var result []byte
    result, _ = json.Marshal(r)
    return string(result)
}

// Type to represent a list of Roars
type RoarList struct {
    // anonymous embedding is like inheriting
    vector.Vector
}

// Custom string representation of a RoarList
func (rl RoarList) String() string {
    var result string
    for i, r := range roarList.Data() {
        result += fmt.Sprintf("[%d]: %v\n", i, r)
    }
    return result
}

// Marshal RoarList to Json string representation
func (rl RoarList) toJson() string {
    var result []byte
    result, _ = json.Marshal(rl)
    return string(result)
}

// (Global) RoarList instance to track Roars
var roarList RoarList

/*

interface ContentHandler{
  func handle(ctx
}
type ContentDipatcher struct {
  contentHandlerMap map[string]:
}

*/

// Return list of Roars as string represenation
func getRoarsAsString() string {
    return fmt.Sprintf("%v\n", roarList)
}

// Return list of Roars as JSON represenation
func getRoarsAsJson() string {
    return roarList.toJson()
}

// Return a specific Roar as JSON represenation
// Param val needs to be Atoi parsable and point to a valid index.
func getRoarAsJson(ctx *web.Context, val string) {
    log.Stderrf("Request.Headers:\n%v\n", ctx.Request.Headers)
    var i, err = strconv.Atoi(val)
    if err != nil {
        ctx.StartResponse(400)
        ctx.WriteString(err.String())
        return
    }
    if i >= roarList.Len() {
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("No roar with id: %d", i))
        return
    }
    ctx.WriteString(roarList.At(i).(*Roar).toJson())
}

// Create a new Roar from user input.
//
// Reads body of request tries to unmarshal a Roar from it. A lot of
// things can go wrong here *g*. A curl command that works is:
//
// $ curl -i -X POST -d \
// '{"Author":"Sven","Text":"Hello!!!!","CreationDate":"Sat, 10 Jul 2010 19:46:48 CEST"}' \
// -H "Content-Type: application/json" \
// http://localhost:9999/roars
//
func postRoarAsJson(ctx *web.Context) {
    log.Stderrf("Request:\n%v\n", ctx.Request)

    // read request body to buf
    var buf, err = ioutil.ReadAll(ctx.Request.Body)
    if nil != err {
        ctx.StartResponse(400)
        ctx.WriteString(err.String())
        return
    } else if 0 == len(buf) {
        ctx.StartResponse(400)
        ctx.WriteString("Empty request body.\n")
        return
    }
    log.Stderrf("Read %d bytes from Request.Body:", len(buf))
    log.Stderrf("\n%s\n", string(buf))

    // try unmarshal bytes and extract a Roar
    var r *Roar = &Roar{}
    err = json.Unmarshal(buf, r)
    if nil != err {
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("Couldn't parse input: %s\n",
            err.String()))
        return
    }

    // append new roar
    roarList.Push(r)
    log.Stdoutf("Created new Roar with id %d:\n%v\n",
        r, roarList.Len())

    //tell client id of new roar
    ctx.WriteString(fmt.Sprintf("Created new Roar with id: %d\n",
        roarList.Len()))
}

func getRoarAsString(ctx *web.Context, val string) {
    var i, err = strconv.Atoi(val)
    if err != nil {
        ctx.StartResponse(400)
        ctx.WriteString(err.String())
        return
    }
    if i >= roarList.Len() {
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("No roar with id: %d", i))
        return
    }
    ctx.WriteString(roarList.At(i).(*Roar).String())
}

func deleteRoar(ctx *web.Context, val string) {
    var i, err = strconv.Atoi(val)
    if err != nil {
        ctx.StartResponse(400)
        ctx.WriteString(err.String())
        return
    }
    if i >= roarList.Len() {
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("No roar with id: %d", i))
        return
    }
    roarList.Delete(i)
    ctx.WriteString(fmt.Sprintf("Deleted %d", i))
}

func main() {
    // initialize our roarList with some Roars
    for i := 0; i < 3; i++ {
        roarList.Push(NewRoar("Sven", fmt.Sprintf("Hello %d!", i)))
    }
    // register the dispatchers with web.go
    web.Get("/roars", getRoarsAsString)
    web.GetAs("/roars", "application/json", getRoarsAsJson)
    web.PostAs("/roars", "application/json", postRoarAsJson)
    web.Get("/roars/(.*)", getRoarAsString)
    web.GetAs("/roars/(.*)", "application/json", getRoarAsJson)
    web.Delete("/roars/(.*)", deleteRoar)
    web.Run("0.0.0.0:9999")
}
