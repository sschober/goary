package main

import (
    //    "github.com/hoisie/web.go"
    //    "github.com/sschober/web.go"
    "web"
    "fmt"
    "strconv"
    "strings"
    "json"
    "time"
    "io/ioutil"
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
//
// Note:
//  It might be convenient to chose a Vector here, but, alas, Go
//  does not support generics/templates. So there is only IntVector
//  and StringVector.
// TODO:
//  - Investigate alternatives
type RoarList []*Roar

// Custom string representation of a RoarList
func (rl RoarList) String() string {
    var result string
    for i, r := range roarList {
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
    var i, err = strconv.Atoi(val)
    if err != nil {
        ctx.StartResponse(400)
        ctx.WriteString(err.String())
        return
    }
    if i >= len(roarList) {
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("No roar with id: %d", i))
        return
    }
    fmt.Printf("%v", ctx.Request.Headers)
    ctx.WriteString(roarList[i].toJson())
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
    fmt.Printf("\n%v\n", ctx.Request)

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
    fmt.Printf("Read %d bytes", len(buf))
    fmt.Printf("\n%s\n", string(buf))

    // try unmarshal bytes and extract a Roar
    var r *Roar = &Roar{}
    err = json.Unmarshal(buf, r)
    if nil != err {
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("Couldn't parse input: %s\n",
            err.String()))
        return
    }
    fmt.Printf("New Roar: %v", r)

    // append new roar
    // yay, this is efficient, baby *g*
    var roarListNew = make(RoarList, len(roarList)+1)
    for i, r := range roarList {
        roarListNew[i] = r
    }
    roarListNew[len(roarListNew)-1] = r
    roarList = roarListNew

    //tell client id of new roar
    ctx.WriteString(fmt.Sprintf("Created new Roar with id: %d\n", len(roarListNew)-1))
}

func getRoarAsString(ctx *web.Context, val string) {
    var i, err = strconv.Atoi(val)
    if err != nil {
        ctx.StartResponse(400)
        ctx.WriteString(err.String())
        return
    }
    if i >= len(roarList) {
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("No roar with id: %d", i))
        return
    }
    fmt.Printf("%v", ctx.Request.Headers)
    ctx.WriteString(roarList[i].String())
}

// TODO:
//  This whole (web.go)-client side dispatching is insane: It ought
//  to be done in web.go itself. Maybe via an addition to the api:
//  add the ability to specify a content type to a handler.

//
func getRoarDispatcher(ctx *web.Context, val string) {
    ct, _ := ctx.Request.Headers["Content-Type"]
    switch strings.Split(ct, ";", 2)[0] {
    case "text/plain", "":
	getRoarAsString(ctx, val);
    case "application/json":
        getRoarAsJson(ctx, val)
    default:
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("Unknown Content-Type: %s\n", ct))
    }
}
//
func getRoarsDispatcher(ctx *web.Context) {
    ct, _ := ctx.Request.Headers["Content-Type"]
    switch strings.Split(ct, ";", 2)[0] {
    case "text/plain", "":
	ctx.WriteString(getRoarsAsString())
    case "application/json":
        ctx.WriteString(getRoarsAsJson())
    default:
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("Unknown Content-Type: %s\n", ct))
    }
}
//
func postDispatcher(ctx *web.Context) {
    ct, _ := ctx.Request.Headers["Content-Type"]
    switch strings.Split(ct, ";", 2)[0] {
    case "text/plain", "application/x-www-form-urlencoded", "":
        // not yet handled
    case "application/json":
        postRoarAsJson(ctx)
    default:
        ctx.StartResponse(400)
        ctx.WriteString(fmt.Sprintf("Unknown Content-Type: %s\n", ct))
    }
}

func main() {
    roarList = make(RoarList, 3)
    for i := 0; i < 3; i++ {
        roarList[i] = NewRoar("Sven", fmt.Sprintf("Hello %d!", i))
    }
    /*
    var getRoarsDispatcher = ContentDispatcher
*/
    web.Get("/roars", getRoarsDispatcher)
    web.Post("/roars", postDispatcher)
    web.Get("/roars/(.*)", getRoarDispatcher)
    web.Run("0.0.0.0:9999")
}
