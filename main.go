package main

import (
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/sdl_ttf"
    "os"
)

const (
    RUNE_A string = "ᚠ"
    RUNE_B        = "ᚨ"
    RUNE_C        = "ᚰ"
    RUNE_D        = "ᚸ"
)

type RuneResource struct {
    Txtr  *sdl.Texture
    Rect  *sdl.Rect
    Value string
}

var (
    CL_WHITE = sdl.Color{255, 255, 255, 255}
    CL_BLACK = sdl.Color{0, 0, 0, 255}
)

func HandleEvents() bool {
    var event sdl.Event
    for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
        switch event.(type) {
        case *sdl.QuitEvent:
            return false
            break
        }
    }

    return true
}

func GetFullRect(s *sdl.Surface) *sdl.Rect {
    return &sdl.Rect{0, 0, s.W, s.H}
}

func BakeText(r *sdl.Renderer, font *ttf.Font, msg string, color sdl.Color) (*sdl.Texture, *sdl.Rect, error) {
    surf, err := font.RenderUTF8_Blended(msg, color)
    if err != nil {
        return nil, nil, err
    }

    rect := GetFullRect(surf)
    t, err := r.CreateTextureFromSurface(surf)

    return t, rect, err
}

func LoadRune(r *sdl.Renderer, font *ttf.Font, runa string, color sdl.Color) *RuneResource {
    fmt.Printf("Loading rune %q\n", runa)
    txtr, rect, _ := BakeText(r, font, runa, color)

    return &RuneResource{
        Txtr:  txtr,
        Rect:  rect,
        Value: runa,
    }
}

func LoadAllRunes(r *sdl.Renderer, font *ttf.Font, color sdl.Color, runes ...string) []*RuneResource {
    allRunes := make([]*RuneResource, len(runes))

    for i, runa := range runes {
        allRunes[i] = LoadRune(r, font, runa, color)
    }

    return allRunes
}

func run() int {
    var (
        running  bool
        window   *sdl.Window
        renderer *sdl.Renderer
        err      error
    )

    fmt.Printf("Stating RuneGrid - Day 001\n")

    window, renderer, err = sdl.CreateWindowAndRenderer(320, 540, 0)
    if err != nil {
        fmt.Printf("Error at sdl.CreateWindowAndRenderer: %v", err)
        return 1
    }

    ttf.Init()

    runeFont, err := ttf.OpenFont("assets/fonts/NotoSansRunic-Regular.ttf", 32)
    if err != nil {
        fmt.Printf("Error at ttf.OpenFont: %v", err)
        return 1
    }

    window.SetTitle("RuneGrid - Day 001")

    runesR := LoadAllRunes(renderer, runeFont, CL_WHITE, RUNE_A, RUNE_B, RUNE_C, RUNE_D)
    println(len(runesR))

    running = true
    for running {
        running = HandleEvents()
        renderer.Clear()

        for i, rr := range runesR {
            dr := *rr.Rect
            dr.Y = int32(i * 32)
            renderer.Copy(rr.Txtr, rr.Rect, &dr)
        }

        renderer.Present()
        sdl.Delay(16)
    }

    sdl.Quit()
    return 0
}

func main() {
    os.Exit(run())
}
