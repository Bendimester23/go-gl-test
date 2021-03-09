package main

import (
	"fmt"
	"log"
	"runtime"

	"github.com/Bendimester23/gogl-test/shaders"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	window   *glfw.Window
	program  uint32
	triangle = []float32{
		0, 0.5, 0,
		-0.5, -0.5, 0,
		0.5, -0.5, 0,
	}
	vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
` + "\x00"

	fragmentShaderSource = `
    #version 410
    out vec4 frag_colour;
    void main() {
        frag_colour = vec4(1, 1, 1, 1);
    }
` + "\x00"
)

func main() {
	runtime.LockOSThread()
	log.Println(fmt.Sprintf("Started! %o", runtime.NumCPU()))
	glfw.Init()
	//window = CreateWindow(1280, 720, "First Go GLFW Window")
	if err := glfw.Init(); err != nil {
		log.Panicln(fmt.Sprintf("Error initialising glfw!\n Error: %e", err))
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, 0)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	w, erro := glfw.CreateWindow(1280, 720, "First GO-GLFW window", nil, nil)
	if erro != nil {
		log.Panicln(fmt.Sprintf("Error creating window!\n Error: %e", erro))
	}

	w.MakeContextCurrent()
	window = w
	program = InitGl()

	vao := makeVao(triangle)
	for !window.ShouldClose() {
		Draw(vao)
	}
}

func makeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func Draw(vao uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle)/3))

	window.SwapBuffers()
	glfw.PollEvents()
}

func InitGl() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)
	vertexShader, err := shaders.CompileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := shaders.CompileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)

	return prog
}

func CreateWindow(width int, height int, title string) *glfw.Window {
	if err := glfw.Init(); err != nil {
		log.Panicln(fmt.Sprintf("Error initialising glfw!\n Error: %e", err))
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, 0)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	w, erro := glfw.CreateWindow(width, height, title, nil, nil)
	if erro != nil {
		log.Panicln(fmt.Sprintf("Error creating window!\n Error: %e", erro))
	}

	w.MakeContextCurrent()
	return w
}
