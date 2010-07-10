package main

import (
    "github.com/hoisie/web.go"
    "fmt"
    "strconv"
    "json"
    "time"
)

type Roar struct {
    author       string
    text         string
    creationDate string
}

func NewRoar(author string, text string) *Roar {
    return &Roar{author: author, text: text,
        creationDate: time.LocalTime().Format(time.RFC1123)}
}

func (r Roar) String() string {
    return fmt.Sprintf("%s - %s\n%s", r.author, r.creationDate, r.text)
}

func (r Roar) toJson() string {
    var result []byte
    result, _ = json.Marshal(r)
    return string(result)
}

type RoarList []*Roar

func (rl RoarList) String() string {
    var result string
    for i, r := range roarList {
        result += fmt.Sprintf("[%d]: %v\n", i, r)
    }
    return result
}

func (rl RoarList) toJson() string {
    var result []byte
    result, _ = json.Marshal(rl)
    return string(result)
}

var roarList RoarList

func getRoars() string {
    //return fmt.Sprintf("%v\n", roarList)
    return roarList.toJson()
}

func getRoar(ctx *web.Context, val string) {
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

func main() {
    roarList = make(RoarList, 3)
    for i := 0; i < 3; i++ {
        roarList[i] = NewRoar("Sven", fmt.Sprintf("Hello %d!", i))
    }
    web.Get("/roars", getRoars)
    web.Get("/roars/(.*)", getRoar)
    web.Run("0.0.0.0:9999")
}
