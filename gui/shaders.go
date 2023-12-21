// MIT License · Daniel T. Gorski · dtg [at] lengo [dot] org · 12/2023

package gui

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"retro/gui/files"
	"runtime"
)

func createDefaultShader() (err error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	// From gui/files/*
	defaultV := files.MustLoad(files.SHADER_DEFAULT_VERT)
	defaultF := files.MustLoad(files.SHADER_DEFAULT_FRAG)

	var vert, frag uint32
	if vert, err = compileShader(defaultV, gl.VERTEX_SHADER); err != nil {
		return err
	}
	if frag, err = compileShader(defaultF, gl.FRAGMENT_SHADER); err != nil {
		return err
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vert)
	gl.AttachShader(prog, frag)
	gl.LinkProgram(prog)
	gl.UseProgram(prog)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(verCoords)*4, gl.Ptr(verCoords), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(0)

	var tbo uint32
	gl.GenBuffers(1, &tbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, tbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(texCoords)*4, gl.Ptr(texCoords), gl.STATIC_DRAW)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 0, nil)
	gl.EnableVertexAttribArray(1)

	return nil
}

func compileShader(source []byte, shaderType uint32) (uint32, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	src, free := gl.Strs(string(source) + "\x00")
	defer free()

	sh := gl.CreateShader(shaderType)
	gl.ShaderSource(sh, 1, src, nil)
	gl.CompileShader(sh)

	var ok int32
	if gl.GetShaderiv(sh, gl.COMPILE_STATUS, &ok); ok == gl.FALSE {
		msg := "failed to compile %s shader %08X: %s"

		var l int32
		gl.GetShaderiv(sh, gl.INFO_LOG_LENGTH, &l)

		buf := make([]byte, l)
		gl.GetShaderInfoLog(sh, l, &l, &buf[0])

		var typ = "vertex"
		if shaderType == gl.FRAGMENT_SHADER {
			typ = "fragment"
		}
		return 0, fmt.Errorf(msg, typ, shaderType, string(buf))
	}

	return sh, nil
}

var (
	// Vertex:
	//          y
	//  -1,+1   |    +1,+1
	//         ─┼─> x
	//  -1,-1   |    +1,-1
	verCoords = []float32{
		-1, -1, 0,
		+1, +1, 0,
		-1, +1, 0,

		-1, -1, 0,
		+1, +1, 0,
		+1, -1, 0,
	}

	// Texture:
	//         y
	//   0,1   |    1,1
	//        ─┼─> x
	//   0,0   |    1,0
	texCoords = []float32{
		0, 1,
		1, 0,
		0, 0,

		0, 1,
		1, 0,
		1, 1,
	}
)
