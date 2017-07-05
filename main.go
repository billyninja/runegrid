package main

import (
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/sdl_ttf"
    "math"
    "os"
    "time"
)

var (
    CL_WHITE         = sdl.Color{255, 255, 255, 255}
    CL_GREEN         = sdl.Color{0, 255, 0, 255}
    CL_BLACK         = sdl.Color{0, 0, 0, 255}
    RUNES            = [4]*RuneResource{}
    NOW              time.Time
    PREVIOUS_NOW     time.Time
    FRAME_LATENCY    time.Duration
    FRAME_LATENCY_MS float64
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

type Fx struct {
    Min      float64
    Max      float64
    Current  float64
    Duration time.Duration
    Elapsed  time.Duration
}

func (f *Fx) Tick() {
    f.Elapsed += FRAME_LATENCY

    RElp := float64(f.Elapsed) / float64(f.Duration)
    RPg := f.Current / f.Max

    stride := math.Cos(RElp)
    dl := (stride * (f.Max - f.Min)) / (f.Duration.Seconds() * FRAME_LATENCY_MS)
    f.Current += dl

    fmt.Printf("s:%.2f - d:%.2f - el:%.2f - %.2f \n", stride, dl, RElp, RPg)

    if f.Current >= f.Max || f.Elapsed >= f.Duration {
        f.Elapsed = 0
        f.Current = f.Min
    }
}

type RuneResource struct {
    Font      *ttf.Font
    FontSize  int32
    Value     string
    Texture   *sdl.Texture
    BoundRect *sdl.Rect
}

type RuneInstance struct {
    Resource  *RuneResource
    CurrScale *sdl.Rect
    Rotation  *Fx
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
            rn := slot.Holding
            rnR := slot.Holding.Resource

            dest := &sdl.Rect{
                slot.Rect.X + (slot.Rect.W / 2) - (rn.CurrScale.W / 2),
                slot.Rect.Y + (slot.Rect.H / 2) - (rn.CurrScale.H / 2),
                rn.CurrScale.W,
                rn.CurrScale.H,
            }
            //renderer.Copy(rnR.Texture, rnR.BoundRect, dest)
            renderer.CopyEx(rnR.Texture, rnR.BoundRect, dest, float64(rn.Rotation.Current), nil, 0)
            if slot.Active {
                SetColor(renderer, base_color)
            }
        }
    }
}

func SinWave(v int64) float64 {
    return math.Sin(float64(v))
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
                Holding: &RuneInstance{
                    CurrScale: &sdl.Rect{x, y, g.Size, g.Size},
                    Resource:  RUNES[i],
                    Rotation:  InstFx(0, 360, time.Second*10, true),
                },
            }
        }
    }
}

func InstFx(min, max float64, d time.Duration, loop bool) *Fx {
    return &Fx{
        Min:      min,
        Max:      max,
        Current:  min,
        Duration: d,
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

func update(g *Grid) {
    for i := 0; i < len(g.Slots); i++ {
        for j := 0; j < len(g.Slots[i]); j++ {
            slot := g.Slots[i][j]
            slot.Holding.Rotation.Tick()
        }
    }
}

func run() int {
    var (
        running bool = true
        grid    *Grid
    )

    window, renderer, err := sdl.CreateWindowAndRenderer(320, 480, 0)
    window.SetTitle("RuneGrid - Day 002")
    if err != nil {
        fmt.Printf("Error sdl.CreateWindowAndRenderer %v", err)
        return 1
    }

    ttf.Init()
    f32, err := ttf.OpenFont("assets/fonts/NotoSansRunic-Regular.ttf", 24)
    LoadAllRunes(renderer, f32, CL_WHITE, RUNA_A, RUNA_B, RUNA_C, RUNA_D)

    grid = &Grid{
        RootPos: Vector2d{32, 32},
        Size:    60,
        Spacing: 4,
    }
    grid.BuildRects()

    NOW = time.Now()
    for running {
        // Updating global timers
        PREVIOUS_NOW = NOW
        NOW = time.Now()
        FRAME_LATENCY = NOW.Sub(PREVIOUS_NOW)
        FRAME_LATENCY_MS = float64(FRAME_LATENCY.Nanoseconds() / 1000000)
        running = HandleEvents()
        update(grid)

        SetColor(renderer, CL_BLACK)
        renderer.Clear()
        grid.Render(renderer, CL_WHITE)
        renderer.Present()

        fmt.Printf("%v - %dhz\n", FRAME_LATENCY, time.Second/FRAME_LATENCY)
        sdl.Delay(16)
    }

    return 0
}

func main() {
    os.Exit(run())
}
