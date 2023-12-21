#version 330 core

layout (location = 0) in vec3 vertex;
layout (location = 1) in vec2 texCoords;

out vec2 texCoord;

// Default vertex shader.
void main() {
    gl_Position = vec4(vertex.x, vertex.y, 0, 1);
    texCoord = texCoords;
}
