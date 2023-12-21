// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package gui

import (
	"errors"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/png"
	"log"
	"retro/emu/device/render"
	"retro/emu/input"
	"retro/emu/virtual"
	"retro/gui/files"
	"runtime"
	"time"
	"unsafe"
)

type (
	// Window represents the GUI (X11) output window.
	Window struct {
		properties Properties
		renderer   *render.Driver
		channels   *virtual.Channels
		window     *glfw.Window
	}

	// Properties provides resource parameters for the GUI.
	Properties struct {
		Width  int
		Height int
		Title  string
	}
)

// NewWindow creates a new X11 Window.
func NewWindow(properties Properties, renderer *render.Driver, channels *virtual.Channels) *Window {
	return &Window{properties: properties, renderer: renderer, channels: channels}
}

// Clipboard returns the content of the clipboard.
func (win *Window) Clipboard() string {
	return win.window.GetClipboardString()
}

// Properties returns the window properties.
func (win *Window) Properties() Properties {
	return win.properties
}

// Open opens an X11 Window.
func (win *Window) Open() (err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	defer func() {
		if r := recover(); r != nil {
			err = errors.New(r.(*glfw.Error).Error())
		}
	}()

	if err = glfw.Init(); err != nil {
		return err
	}
	if err = gl.Init(); err != nil {
		return err
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	w := win.properties.Width
	h := win.properties.Height
	t := win.properties.Title

	win.window, err = glfw.CreateWindow(w, h, t, nil, nil)
	if err != nil {
		return err
	}

	// From gui/files/*
	loadIcon := func(name string) image.Image {
		icon, _ := png.Decode(files.MustOpen(name))
		return icon
	}
	win.window.SetIcon([]image.Image{
		loadIcon(files.WINDOW_ICON_16x),
		loadIcon(files.WINDOW_ICON_24x),
		loadIcon(files.WINDOW_ICON_32x),
		loadIcon(files.WINDOW_ICON_48x),
	})

	win.window.MakeContextCurrent()

	// After context is made current.
	glfw.SwapInterval(1) // VSYNC on = 1, off = 0

	// ---

	keyFunc := func(w *glfw.Window, k glfw.Key, c int, a glfw.Action, m glfw.ModifierKey) {
		win.channels.KeyInput() <- input.NewKeyInput(int(k), int(a), int(m))
	}
	butFunc := func(w *glfw.Window, b glfw.MouseButton, a glfw.Action, m glfw.ModifierKey) {
		win.channels.MouseButton() <- input.NewMouseButton(int(b), int(a), int(m))
	}
	posFunc := func(w *glfw.Window, x float64, y float64) {
		win.channels.CursorPos() <- input.NewCursorPos(x, y)
	}

	// Install event listeners.
	win.window.SetKeyCallback(keyFunc)
	win.window.SetMouseButtonCallback(butFunc)
	win.window.SetCursorPosCallback(posFunc)

	return nil
}

// Close closes the window.
func (*Window) Close() error {
	return nil
}

// RenderAndListen starts the main rendering loop.
func (win *Window) RenderAndListen() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	if false {
		gl.Enable(gl.DEBUG_OUTPUT)
		gl.DebugMessageCallback(message, nil)
	}

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.DST_ALPHA)

	var tex uint32
	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_BASE_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.MIRROR_CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.MIRROR_CLAMP_TO_EDGE)

	w := int32(render.Width)
	h := int32(render.Height)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	if err := createDefaultShader(); err != nil {
		return err
	}

	var frame []byte
	for !win.window.ShouldClose() {
		flash := time.Now().UnixMilli()%1000 > 450
		frame = win.renderer.Render(flash)

		gl.Clear(gl.COLOR_BUFFER_BIT)
		gl.TexSubImage2D(gl.TEXTURE_2D, 0, 0, 0, w, h, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(frame))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(verCoords)/3))

		win.window.SwapBuffers()
		glfw.PollEvents()

		time.Sleep(time.Second / 36) // ~ 30 f/s here
	}
	return nil
}

func message(_ uint32, typ uint32, _ uint32, s uint32, _ int32, m string, _ unsafe.Pointer) {
	log.Printf("[type: 0x%X, severity: 0x%X] %s\n", typ, s, m)
}
