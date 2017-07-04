package main

import (
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/sdl_ttf"
    "os"
)

var (
    CL_WHITE = sdl.Color{255, 255, 255, 255}
    CL_GREEN = sdl.Color{0, 255, 0, 255}
    CL_BLACK = sdl.Color{0, 0, 0, 255}
    RUNES    = [4]*RuneResource{}
)

const (
    RUNA_A string = "ᚠ"
    RUNA_B        = "ᚨ"
    RUNA_C        = "ᚰ"
    RUNA_D        = "ᚸ"
)

type Vector2d struct {
    X int32
    Y int32
}

type RuneResource struct {
    Font      *ttf.Font
    FontSize  int32
    Value     string
    Texture   *sdl.Texture
    BoundRect *sdl.Rect
}

type RuneInstance struct {
    Resource *RuneResource
}

type GridSlot struct {
    Rect    *sdl.Rect
    Active  bool
    Holding *RuneInstance
}

type Grid struct {
    RootPos Vector2d
    Size    int32
    Spacing int32
    Slots   [4][4]*GridSlot
}

func SetColor(r *sdl.Renderer, c sdl.Color) {
    r.SetDrawColor(c.R, c.G, c.B, c.A)
}

func (g *Grid) Render(renderer *sdl.Renderer, base_color sdl.Color) {
    SetColor(renderer, base_color)
    for i := range g.Slots {
        for j := range g.Slots[i] {
            slot := g.Slots[i][j]
            if slot.Active {
                SetColor(renderer, CL_GREEN)
            }

            renderer.DrawRect(slot.Rect)
            rn := RUNES[j]
            dest := &sdl.Rect{
                slot.Rect.X + (slot.Rect.W / 2) - (rn.BoundRect.W / 2),
                slot.Rect.Y + (slot.Rect.H / 2) - (rn.BoundRect.H / 2),
                rn.BoundRect.W,
                rn.BoundRect.H,
            }
            renderer.Copy(rn.Texture, rn.BoundRect, dest)

            if slot.Active {
                SetColor(renderer, base_color)
            }
        }
    }
}

func (g *Grid) BuildRects() {
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            x := g.RootPos.X + ((g.Size * int32(i)) + (g.Spacing * int32(i)))
            y := g.RootPos.Y + ((g.Size * int32(j)) + (g.Spacing * int32(j)))
            rect := &sdl.Rect{x, y, g.Size, g.Size}
            g.Slots[i][j] = &GridSlot{
                Rect:   rect,
                Active: ((i+j)%2 == 0),
            }
        }
    }
}

func HandleEvents() bool {
    var ev sdl.Event
    for ev = sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
        switch ev.(type) {
        case *sdl.QuitEvent:
            return false
        }
    }
    return true
}

func LoadRune(renderer *sdl.Renderer, font *ttf.Font, color sdl.Color, runa string) *RuneResource {
    surf, _ := font.RenderUTF8_Blended(runa, color)
    txtr, _ := renderer.CreateTextureFromSurface(surf)

    return &RuneResource{
        Font:      font,
        Value:     runa,
        Texture:   txtr,
        BoundRect: &sdl.Rect{0, 0, surf.W, surf.H},
    }
}

func LoadAllRunes(renderer *sdl.Renderer, font *ttf.Font, color sdl.Color, runes ...string) {
    for i, runa := range runes {
        RUNES[i] = LoadRune(renderer, font, color, runa)
    }
}

func run() int {
    var (
        running bool = true
        grid    *Grid
    )

    window, renderer, err := sdl.CreateWindowAndRenderer(320, 480, 0)
    if err != nil {
        fmt.Printf("Error sdl.CreateWindowAndRenderer %v", err)
        return 1
    }
    ttf.Init()
    f32, err := ttf.OpenFont("assets/fonts/NotoSansRunic-Regular.ttf", 24)

    window.SetTitle("RuneGrid - Day 002")
    grid = &Grid{
        RootPos: Vector2d{32, 32},
        Size:    60,
        Spacing: 4,
    }
    grid.BuildRects()

    LoadAllRunes(renderer, f32, CL_WHITE, RUNA_A, RUNA_B, RUNA_C, RUNA_D)
    println(len(RUNES))

    for running {
        running = HandleEvents()
        SetColor(renderer, CL_BLACK)
        renderer.Clear()
        grid.Render(renderer, CL_WHITE)
        renderer.Present()
    }

    return 0
}

func main() {
    os.Exit(run())
}
