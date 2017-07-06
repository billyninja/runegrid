package main

import (
    "github.com/veandco/go-sdl2/sdl"
    "os"
)

const (
    KEY_ARROW_UP    = 1073741906
    KEY_ARROW_DOWN  = 1073741905
    KEY_ARROW_LEFT  = 1073741904
    KEY_ARROW_RIGHT = 1073741903
)

type Vector2d struct {
    X int32
    Y int32
}

type GridSlot struct {
    Active bool
    Rect   *sdl.Rect
}

type Grid struct {
    TileSize  int32
    Spacing   int32
    RootPoint Vector2d
    Slots     [4][4]*GridSlot
    CActive   [2]int
}

func BuildGrid(g *Grid) {
    for i := range g.Slots {
        for j := range g.Slots[i] {

            x := g.RootPoint.X + (int32(i) * g.Spacing) + (int32(i) * g.TileSize)
            y := g.RootPoint.Y + (int32(j) * g.Spacing) + (int32(j) * g.TileSize)

            g.Slots[i][j] = &GridSlot{
                Rect: &sdl.Rect{x, y, g.TileSize, g.TileSize},
            }
        }
    }
    g.Slots[0][0].Active = true
}

func handleKD(grid *Grid, kc sdl.Keycode) {
    grid.Slots[grid.CActive[0]][grid.CActive[1]].Active = false
    switch kc {
    case KEY_ARROW_UP:
        grid.CActive[1] -= 1
        break
    case KEY_ARROW_DOWN:
        grid.CActive[1] += 1
        break
    case KEY_ARROW_LEFT:
        grid.CActive[0] -= 1
        break
    case KEY_ARROW_RIGHT:
        grid.CActive[0] += 1
        break
    }

    for i := range grid.CActive {
        if grid.CActive[i] < 0 {
            grid.CActive[i] = 0
        }
        if grid.CActive[i] >= len(grid.Slots)-1 {
            grid.CActive[i] = len(grid.Slots) - 1
        }
    }

    grid.Slots[grid.CActive[0]][grid.CActive[1]].Active = true
}

func run() int {
    running := true
    w, r, _ := sdl.CreateWindowAndRenderer(320, 540, 0)
    w.SetTitle("Day 003 - Arrow Controls")

    ts := int32((320 - ((32 * 2) + (8 * 4))) / 4)

    Grid := &Grid{
        TileSize:  ts,
        Spacing:   8,
        RootPoint: Vector2d{32, 32},
    }

    BuildGrid(Grid)

    for running {
        for ev := sdl.PollEvent(); ev != nil; ev = sdl.PollEvent() {
            switch t := ev.(type) {
            case *sdl.QuitEvent:
                running = false
                break
            case *sdl.KeyDownEvent:
                handleKD(Grid, t.Keysym.Sym)
                break
            }
        }

        r.SetDrawColor(0, 0, 0, 0)

        r.Clear()
        r.SetDrawColor(255, 255, 255, 255)
        for i := range Grid.Slots {
            for j := range Grid.Slots[i] {
                s := Grid.Slots[i][j]
                if s.Active {
                    r.SetDrawColor(0, 255, 0, 255)
                }

                r.DrawRect(s.Rect)

                if s.Active {
                    r.SetDrawColor(255, 255, 255, 255)
                }
            }
        }
        r.Present()
        sdl.Delay(16)
    }

    return 0
}

func main() {
    os.Exit(run())
}
