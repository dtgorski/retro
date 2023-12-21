#version 330 core

in  vec2 texCoord;
out vec4 color;

uniform sampler2D tex;

// Default fragment shader.
void main() {
    color = texture(tex, texCoord);

    // Scanline, kind of.
    if (int(texCoord.y * 192 * 3) % 3 == 0) {
        color = color * vec4(1, 1, 1, 0.7);
    }
}
