package strx_avatar

import (
    "bytes"
    "encoding/xml"
	"hash/fnv"
	"strings"
	"errors"
)

type ComponentFunc func(components ComponentPickCollection, colors ColorPickCollection) string
type ComponentGroup map[string]ComponentFunc

type ComponentPick struct {
    Name  string
    Value ComponentFunc
}

type ComponentPickCollection map[string]*ComponentPick
type ColorPickCollection map[string]string

type AvatarOptions struct {
    Base      string
    Hair      string
    Ears      string
    Earrings  string
    Eyebrows  string
    Eyes      string
    Glasses   string
    Nose      string
    Mouth     string
    Shirt     string
    FacialHair string
}

func escapeXML(s string) string {
    var b bytes.Buffer
    xml.EscapeText(&b, []byte(s))
    return b.String()
}

func color(colors ColorPickCollection, key string) string {
    return escapeXML(colors[key])
}

func renderOptional(components ComponentPickCollection, key string, colors ColorPickCollection) string {
    if c, ok := components[key]; ok && c != nil && c.Value != nil {
        return c.Value(components, colors)
    }
    return ""
}

func sprintf(format string, args ...string) string {
    size := len(format)
    for _, a := range args {
        size += len(a)
    }

    var b strings.Builder
    b.Grow(size)

    argIndex := 0
    for i := 0; i < len(format); i++ {
        if format[i] == '%' && i+1 < len(format) && format[i+1] == 's' {
            if argIndex < len(args) {
                b.WriteString(args[argIndex])
                argIndex++
            }
            i++
        } else {
            b.WriteByte(format[i])
        }
    }

    return b.String()
}

var baseGroup = ComponentGroup{
    "standard": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M154 319.5c-14.4-20-25.67-58.67-27-78L58.5 212 30 319.5h124Z" fill="%s" stroke="#000" stroke-width="4"/>`+
                `<path d="M130.37 263.69c-2.1.2-4.22.31-6.37.31-30.78 0-56.05-21.57-58.76-49.1L127 241.5c.38 5.48 1.55 13.32 3.37 22.19Z" fill="#000" style="mix-blend-mode:multiply"/>`+
                `<path d="M181.94 151.37v.01l.1.4.14.65A75.72 75.72 0 0 1 34.93 187.7l-.2-.74L18 117.13l-.06-.29A75.72 75.72 0 0 1 165.2 81.55l.05.21.02.08.05.2.05.2v.01l16.4 68.44.08.34.08.34Z" fill="%s" stroke="#000" stroke-width="4"/>`+
                `<g transform="translate(34 102.3)">%s</g>`,
            color(colors, "base"),
            color(colors, "base"),
            renderOptional(components, "facialHair", colors),
        )
    },
}

var earringsGroup = ComponentGroup{
    "hoop": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<path d="M24 0A24 24 0 1 1 0 24c0-6.4 3.5-11.5 6.57-16.5L7.5 6" stroke="%s" stroke-width="4"/>`,
            color(colors, "earring"))
    },
    "stud": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<circle cx="25" cy="2" r="4" fill="%s"/><circle cx="26" cy="1" r="1" fill="#fff"/>`,
            color(colors, "earring"))
    },
}

var earsGroup = ComponentGroup{
    "attached": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M30.5 6.18A23.78 23.78 0 0 0 23.08 5c-10.5 0-19 6.5-18 18.5 1.04 12.5 8.5 17 19 17A19.6 19.6 0 0 0 31 39.23" stroke="#000" stroke-width="8"/>`+
                `<path d="M31.5 39.04a19.38 19.38 0 0 1-7.42 1.46c-10.5 0-17.96-4.5-19-17-1-12 7.5-18.5 18-18.5 3.14 0 6.19.6 8.92 1.73l-.5 32.3Z" fill="%s"/>`+
                `<path d="M27.5 13.5c-4-1.83-12.8-2.8-16 8" stroke="#000" stroke-width="4"/>`+
                `<path d="M17 14c2.17 1.83 6.3 7.5 5.5 15.5" stroke="#000" stroke-width="4"/>`+
                `<g transform="translate(3 35)">%s</g>`,
            color(colors, "base"),
            renderOptional(components, "earrings", colors),
        )
    },
    "detached": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M37 8.25V7.13l-.95-.59A24.91 24.91 0 0 0 23.08 3C17.44 3 12.16 4.75 8.4 8.3c-3.8 3.58-5.86 8.83-5.31 15.37.52 6.37 2.66 11.06 6.2 14.17-.29 1-.37 2.08-.24 3.21a8.98 8.98 0 0 0 4.6 7.08C16.09 49.5 19.2 50 22.52 50c5.48 0 10.29-2.95 13.95-6.89l.53-.57V8.25Z" stroke="#000" stroke-width="4"/>`+
                `<path d="M42.97 23.98c.07-.65.1-1.3.1-1.98 0-10.22-9.5-17-20-17C12.6 5 4.09 11.5 5.09 23.5c.56 6.68 2.95 11.07 6.65 13.72a5.7 5.7 0 0 0-.68 3.6C11.68 46.1 16.19 48 22.52 48c11.1 0 19.9-14.05 20.45-24.02Z" fill="%s"/>`+
                `<path d="M27.5 13.5c-4-1.83-12.8-2.8-16 8" stroke="#000" stroke-width="4"/>`+
                `<path d="M17 14c2.17 1.83 6.3 7.5 5.5 15.5" stroke="#000" stroke-width="4"/>`+
                `<g transform="translate(3 42)">%s</g>`,
            color(colors, "base"),
            renderOptional(components, "earrings", colors),
        )
    },
}

var eyebrowsGroup = ComponentGroup{
    "down": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<path d="M27 26.5c6.17 2.5 21.1 3 31.5-15M94 4c5.17 5.33 18.1 12.8 28.5 0" stroke="%s" stroke-width="4" stroke-linecap="round"/>`,
            color(colors, "eyebrows"))
    },
    "eyelashesDown": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<path d="M27 26.5c6.17 2.5 21.1 3 31.5-15M94 4c5.17 5.33 18.1 12.8 28.5 0M37.15 26.46 31 21.03M116.22 9.44l1.78-8M45.6 22.81l-4.05-7.13M108.14 9.02l.94-8.15M52.67 17.2l-2.2-7.9M100 8.03l-.78-8.16" stroke="%s" stroke-width="4" stroke-linecap="round"/>`,
            color(colors, "eyebrows"))
    },
}

var glassesGroup = ComponentGroup{
    "round": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<g stroke="%s" stroke-width="4"><circle cx="122.5" cy="28" r="26"/><circle cx="55.5" cy="37" r="26"/><path d="M97.5 35a8 8 0 0 0-16 0M30 39 0 44.5"/></g>`,
            color(colors, "glasses"))
    },
    "square": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<g stroke="%s" stroke-width="4"><path d="M34.5 42.5 0 49.12" stroke-linecap="round"/><path d="M35.47 18.53 74.2 13.1a6 6 0 0 1 6.77 5.1l5.57 39.62a6 6 0 0 1-5.1 6.78l-34.48 4.84a6 6 0 0 1-6.65-4.48l-9.81-39.01a6 6 0 0 1 4.98-7.4ZM145.92 3.22 107.2 8.66a6 6 0 0 0-5.1 6.78l5.56 39.6a6 6 0 0 0 6.78 5.11l34.47-4.84a6 6 0 0 0 5.16-6.14l-1.32-40.2a6 6 0 0 0-6.83-5.75ZM83.5 37.12l22-3.5"/></g>`,
            color(colors, "glasses"))
    },
}

var eyesGroup = ComponentGroup{
    "eyes": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<g fill="%s"><ellipse cx="16.53" cy="29.4" rx="9" ry="13.5" transform="rotate(-6.78 16.53 29.4)"/><ellipse cx="80.53" cy="19.4" rx="9" ry="13.5" transform="rotate(-6.28 80.53 19.4)"/></g><g transform="translate(-40 -8)">%s</g>`,
            color(colors, "eyes"),
            renderOptional(components, "glasses", colors))
    },
    "smiling": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<path d="M5.29 34.07c.11.82 1.14 1 1.72.41 2.46-2.52 6.25-4.36 10.65-4.89 2.6-.3 5.1-.12 7.32.48.75.2 1.5-.44 1.23-1.17A10.84 10.84 0 0 0 5.3 34.07ZM69.38 24.07c.12.82 1.15 1 1.73.41 2.44-2.48 6.19-4.3 10.54-4.83 2.56-.3 5.03-.12 7.23.47.75.2 1.5-.44 1.23-1.17a10.74 10.74 0 0 0-20.73 5.12Z" fill="%s"/><g transform="translate(-40 -8)">%s</g>`,
            color(colors, "eyes"),
            renderOptional(components, "glasses", colors))
    },
}

var mouthGroup = ComponentGroup{
    "surprised": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M36.37 56.66c12.1-2.2 18.75-15.38 16.42-28.23C50.47 15.6 39.63 5.57 27.52 7.76 15.4 9.95 8.77 23.13 11.09 35.98c2.33 12.85 13.17 22.87 25.28 20.68Z" fill="%s"/><path d="M17.14 42.66c2.78-5.21 8.14-9.25 14.8-10.45 6.66-1.2 13.1.7 17.53 4.6-1.09 8.3-6.37 15-13.74 16.33-7.37 1.33-14.67-3.1-18.6-10.47Z" fill="#FC909F"/>`,
            color(colors, "mouth"),
        )
    },

    "laughing": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M64.73 25.7a36 36 0 0 0 1.18-12.54 4.98 4.98 0 0 0-6.8-4.2c-4.26 1.67-18.03 6.88-27.62 8.2-10.52 1.44-26.66-.32-31.44-.91a4.98 4.98 0 0 0-5.53 5.74 36 36 0 0 0 70.21 3.7Z" fill="%s"/><path d="M51.83 39.55a32 32 0 0 1-37.2 4.62 21.5 21.5 0 0 1 37.2-4.62Z" fill="#FC909F"/>`,
            color(colors, "mouth"),
        )
    },

    "nervous": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M68.42 31.57 67.4 17a8.06 8.06 0 0 0-9.74-7.3c-6.95 1.49-20.1 4.1-29.54 4.76-9.43.66-22.82-.1-29.9-.6a8.06 8.06 0 0 0-8.63 8.58L-9.4 37a8.06 8.06 0 0 0 9.73 7.3c6.95-1.48 20.1-4.1 29.54-4.76 9.44-.66 22.82.1 29.91.61a8.06 8.06 0 0 0 8.63-8.58Z" fill="%s"/><path d="m-.25 17.97 1.6 6.07a6 6 0 0 0 6.22 4.46 6 6 0 0 0-5.54 5.28l-.74 6.23c7.28-1.52 19.34-3.83 28.3-4.46 8.98-.63 21.24-.02 28.66.48l-1.6-6.07a6 6 0 0 0-6.21-4.46 6 6 0 0 0 5.54-5.28l.73-6.24c-7.27 1.53-19.33 3.84-28.3 4.47-8.97.62-21.23.01-28.65-.48Z" fill="#fff"/>`,
            color(colors, "mouth"),
        )
    },

    "smile": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M-.5 17.5c2.5 17 31 25 57 5.5" stroke="%s" stroke-width="4"/>`,
            color(colors, "mouth"),
        )
    },

    "sad": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M13 46c1.72-7.96 8.07-24.77 19.77-28.35 11.7-3.58 17.7 8.46 19.23 14.92" stroke="%s" stroke-width="4"/>`,
            color(colors, "mouth"),
        )
    },

    "pucker": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M26 16.7c4.17-2.34 21-5.3 21 1.5 0 8.5-11.5 8-11.5 8s13.04-3.16 10.5 6c-2.5 9-9.5 5.5-11.5 4.5" stroke="%s" stroke-width="4"/>`,
            color(colors, "mouth"),
        )
    },

    "frown": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M-5 41c3.21-7.96 15.1-24.77 37-28.35 21.9-3.58 33.13 8.46 36 14.92" stroke="%s" stroke-width="4"/>`,
            color(colors, "mouth"),
        )
    },

    "smirk": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M10 24.16c4.94 6.45 12.43 13.6 23.98 11.96 11.55-1.62 16.68-9.6 15.17-16.04" stroke="%s" stroke-width="4"/>`,
            color(colors, "mouth"),
        )
    },
}

var noseGroup = ComponentGroup{
    "pointed": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return `<path d="M16.5 3c0 14 7 25 7 25S20 34 10 32" stroke="#000" stroke-width="4"/>`
    },
    "curve": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return `<path d="M16.5 7c-.33 3.83 0 12.2 4 15 5 3.5-.5 12-10.5 10" stroke="#000" stroke-width="4"/>`
    },
}

var shirtGroup = ComponentGroup{
    "crew": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(`<g stroke="#000" stroke-width="4"><path d="M260.7 91H-12.64C3.67 61.66 26.86 42.98 64.44 34.4c16.02-3.65 34.67-5.47 56.56-5.47 9.46 0 16.81 1.44 23.8 3.35 2.58.7 5.18 1.5 7.84 2.3 4.4 1.34 8.97 2.72 13.91 3.86l.14.03.15.01c46.12 3.8 73.78 24.3 93.85 52.5Z" fill="%s"/><path d="m52.93 36.58 9.15-19.6a1 1 0 0 1 1.25-.51c37.93 13.42 72.43 12.48 104.4 3.57a1 1 0 0 1 1.09.38l13.93 19.05a.98.98 0 0 1-.42 1.5c-33.6 13.2-96.67 10.95-128.91-3.07a.98.98 0 0 1-.49-1.32Z" fill="%s"/><path opacity=".75" d="m52.93 36.58 9.15-19.6a1 1 0 0 1 1.25-.51c37.93 13.42 72.43 12.48 104.4 3.57a1 1 0 0 1 1.09.38l13.93 19.05a.98.98 0 0 1-.42 1.5c-33.6 13.2-96.67 10.95-128.91-3.07a.98.98 0 0 1-.49-1.32Z" fill="#fff"/></g>`,
            color(colors, "shirt"),
            color(colors, "shirt"))
    },
}

var hairGroup = ComponentGroup{
    "fonze": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M210.18 49.4c-1.27 6.05-4.6 11.32-9.43 15.9 9.4 34.06 9.6 53.87 4.38 57.65l-14.8-49.99C158.4 91.7 98.65 94.9 83.69 94.9c-1.55.17-3.02.4-4.43.67-14.65 9-2.6 52.12 11.75 70.43l-11 2c-5.14-24.97-17.41-22.92-26.61-21.38l-.32.05c2.2 13.63 6.72 27.74 10.45 39.32.95 2.99 1.86 5.8 2.66 8.4-.79.1-1.48.3-2.12.47-5.5 1.53-7.41 2.06-33.38-61.97-6.47-15.95-6.03-30.16-.97-42.62-4.78-4.8-14.37-7.14-19.71-7.78 10.44-6.12 20.58-4.87 25.54-3.1.5-.75 1.02-1.49 1.56-2.22-.97-4.41-7.96-9.46-12.11-11.82 8.55-4.3 18.6-2.03 22.98-.2C67.63 47.13 97.03 35.05 122 29 170.81 17.17 189.5.5 189.5.5c20.68 8.5 25.62 25.22 20.68 48.9Z" fill="%s" stroke="#000" stroke-width="4"/>`,
            color(colors, "hair"),
        )
    },

    "mrT": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<g fill="%s"><path opacity=".1" d="M187.99 77.18c-8-6.4-21.84-7-27.5-6.5l-8-26.5c13.6 3.2 32 24 35.5 33Z"/><path d="M85.8 11.76S91.52 7.8 115.74 1.7c24.21-6.1 33.04-3.72 33.04-3.72l11.8 72.84s-8.05-.18-28.04 4.19c-20 4.38-29.56 9.67-29.56 9.67l-17.2-72.9Z"/><path d="M48.99 86.68c-6.8-41.6 23.33-68.17 37-75.5l16.98 73.5c-19.2-39.6-45.33-15.17-54 2Z"/><path opacity=".1" d="M67.49 130.68c-7.2-27.2 22-41.84 35.5-46-7-16.34-23-31-42.5-13-18 30.5-11 54-5.5 72l12.5-13Z"/></g>`,
            color(colors, "hair"),
        )
    },

    "dougFunny": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M140 56c14.67-.67 40.4-8.8 26-36M114 54c14.67-.67 40.4-8.8 26-36M78 65c14.67-.67 40.4-8.8 26-36" stroke="%s" stroke-width="4"/>`,
            color(colors, "hair"),
        )
    },

    "mrClean": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return `<ellipse cx="147.85" cy="58.18" rx="6.86" ry="18.44" transform="rotate(117 147.85 58.18)" fill="#FCFDFF"/>`
    },

    "dannyPhantom": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="M123.79 17.49H123.94a96.78 96.78 0 0 1 62.07 24.36c14.06 12.4 22.45 26.87 25.19 36.73-4.06 2.32-11.01 4.31-19.88 5.95-9.68 1.78-21.3 3.08-33.15 4.01-23.7 1.86-48.2 2.2-59.63 1.96l-6.07-.13 4.8 3.71c2.5 1.93 5.83 3.28 9.34 4.22 3.55.95 7.42 1.54 11.14 1.87 3.82.34 7.55.42 10.64.34-10.59 8.16-24.06 14.44-37.35 19.09a225.88 225.88 0 0 1-39.83 9.92l-2.15.32.5 2.11c3.34 14.43 9.5 39.65 13.62 56.57 1.83 7.5 3.26 13.38 3.87 15.94 1.09 4.56 4.5 11.05 8.4 17.03 3.6 5.52 7.78 10.89 11.32 14.2l-7.84 31.81H49.37c8.34-12.71 10.1-27.4 8.4-42.98-1.84-16.87-7.76-35-14-53.17l-1.85-5.36c-5.69-16.46-11.36-32.88-14.43-48.6-3.4-17.44-3.56-33.75 2.83-48.09 10.34-23.21 28.66-36.7 47-44.12 18.37-7.45 36.61-8.76 46.46-7.71Z" fill="%s" stroke="#000" stroke-width="4"/>`,
            color(colors, "hair"),
        )
    },

    "full": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<path d="m-13.4 312.86-.62-1.64c-.56-1.55-1.35-3.84-2.23-6.74a163.21 163.21 0 0 1-5.28-23.98c-2.81-19.77-2.8-45.8 8.4-71.12 1.8-4.09 4.07-8.3 6.52-12.85 9.47-17.6 21.77-40.46 21.77-82.45 0-30.59 14.84-56.35 36.7-74.51 21.88-18.18 50.7-28.66 78.38-28.66 15.13 0 27.3 1.73 37.07 7.64 9.72 5.87 17.37 16.05 23.05 33.58a3.47 3.47 0 0 0 4.36 2.27c11.31-3.67 28.47.04 42.95 9.5 14.42 9.4 25.62 24.15 25.62 41.91 0 15.43-2.64 25.85-5.22 36-3.12 12.3-6.13 24.16-4 43.5.7 6.45 2.15 11.03 4.16 14.82 1.98 3.73 4.48 6.62 7.12 9.66l.05.07c6.28 7.25 9.13 13.22 10.06 18.47.92 5.23-.05 9.98-1.84 14.9-.9 2.48-2 4.97-3.15 7.59l-.1.22c-1.12 2.53-2.3 5.19-3.35 7.98-2.18 5.77-3.89 12.2-3.72 19.83.15 6.48 1.3 10.91 3.01 14.27 1.7 3.32 3.89 5.44 5.8 7.3l.05.05c1.74 1.68 3.2 3.1 4.27 5.1.96 1.78 1.67 4.13 1.79 7.66a172.14 172.14 0 0 1-87.4 23.9 110.86 110.86 0 0 1-7.28-.28c-6.15-9.4-11.75-24.88-16.1-40.8-4.21-15.46-7.18-31.08-8.3-41.4 37.08-10.72 60.32-48.98 54.73-88.46l-.01-.12a5.97 5.97 0 0 0-.08-.47 76.68 76.68 0 0 0-.43-2.25 486.97 486.97 0 0 0-6.53-28.2 276.64 276.64 0 0 0-7.45-24.2c-2.76-7.5-5.83-14.15-9.1-18.24l-.3-.37-.44-.2c-1.93-.83-3.94-1.77-6.03-2.74-9.93-4.62-21.84-10.17-37.26-10.78-18.78-.74-42.56 5.78-74.7 29.09l-2.02 1.46 1.95 1.57c15.92 12.83 19.37 29.86 18.63 44.3a89.2 89.2 0 0 1-5.24 25.1c-1.16-1.69-1.9-3.82-2.45-6.33a85.19 85.19 0 0 1-.96-5.45l-.05-.3c-.3-2-.63-4.08-1.1-6.09-.96-4.01-2.57-8.02-6.14-10.86-3.58-2.84-8.8-4.25-16.4-3.83l-2.22.13.46 2.18 11.36 53.31.02.08.03.09a79.8 79.8 0 0 0 19.91 32.81 65.49 65.49 0 0 1 1.54 2.68c1.01 1.82 2.37 4.38 3.76 7.33 2.82 5.94 5.66 13.24 6.2 19.2.57 6.05-.96 13.86-2.7 20.31a129.63 129.63 0 0 1-2.84 9.14c-5.03-2.4-9.53-2.23-13.38.01-4.16 2.43-7.21 7.06-9.48 12.22-4.15 9.42-6.14 21.64-7.06 29.22A601.65 601.65 0 0 1 6.2 320.1 353.85 353.85 0 0 1-9 314.64a190.8 190.8 0 0 1-4.4-1.77Z" fill="%s" stroke="#000" stroke-width="3.82"/>`,
            color(colors, "hair"),
        )
    },

    "turban": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<g stroke="#000" stroke-width="4"><path d="M222.73 100.8c0-66.1-36.46-110.8-80.87-110.8C84.96-10 27 11.34 27 112.25c0 24.97 10.66 43.58 25.56 57.29a42.88 42.88 0 0 1-3.5-4.92c-2.88-26.98 29.17-47.7 60.54-67.96 16.65-10.75 33.11-21.39 44.05-32.76 9.6 5.43 37.79 28.2 43.16 37.42 2.88 4.94 7.51 21.87 10.67 41.63 10.34-12.42 15.25-24.84 15.25-42.16Z" fill="%s"/><path d="M154.26 63.25c13.18-11.95 32.89-39.21 31.62-56.92"/></g>`,
            color(colors, "hair"),
        )
    },

    "pixie": func(components ComponentPickCollection, colors ColorPickCollection) string {
        return sprintf(
            `<g stroke="#000"><path d="m105.84 88.82 1.88.3v.08l-.04.16-.12.52c-.12.45-.3 1.09-.6 1.9a41.7 41.7 0 0 1-3.4 6.92c-3.17 5.32-8.7 12.66-18.31 21.6 29.97.74 55-8.92 72.82-19.04a156.35 156.35 0 0 0 21.71-14.87 118.43 118.43 0 0 0 7.5-6.7l.1-.1.01-.01 2.29-2.38.93 3.16-1.84.54 1.84-.54v.05l.05.13.15.52a817.06 817.06 0 0 1 2.69 9.28c1.75 6.14 4.14 14.58 6.66 23.77 5.03 18.35 10.6 39.81 12.7 51.97 3.49 20.32-1.91 35.74-5.1 44.87-.63 1.77-1.17 3.3-1.54 4.6.06.04.12.1.22.16.6.42 1.65.91 3.18 1.4 3.02.96 7.32 1.74 11.92 2.19 4.6.44 9.38.53 13.34.14 1.98-.2 3.7-.5 5.08-.92a6.85 6.85 0 0 0 2.58-1.27c.02-.08.03-.29-.02-.68-.1-.78-.4-1.87-.93-3.32-.8-2.15-1.97-4.8-3.35-7.88l-1.46-3.3c-3.93-8.94-8.65-20.56-9.68-32.4-1.15-13.16 1.45-24.42 3.79-34.54l.11-.48c2.39-10.34 4.38-19.32 2.34-28.42-3.1-13.8-7.32-25.3-13.8-33.57-6.42-8.17-15.13-13.27-27.5-14.21l-1.28-.1-.4-1.22c-5.7-17.57-13.38-29.05-23.18-36.17-9.8-7.12-21.96-10.05-36.94-10.05-27.7 0-57.5 10.52-79.88 28.72-22.35 18.18-37.14 43.9-35.2 74.33 2.13 33.6-.3 59.06-5.8 77.72-5 17-12.59 28.46-21.68 35.14 1.19.61 2.78 1.39 4.77 2.27 5.02 2.24 12.59 5.2 22.67 8.03 19.34 5.43 47.95 10.38 85.71 8.89-8.32-7.15-16.72-15.83-23.37-25-7.63-10.53-13.14-21.94-13.3-32.58-.19-12.62.67-45.67.93-52.5.62-15.9 10.82-28.15 20.54-36.3a96.6 96.6 0 0 1 18.96-12.34l.2-.1.13-.06.1-.05h.02v-.01l.8 1.74Zm0 0 1.88.3.56-3.5-3.23 1.46.79 1.74Zm100.23 126.57.03.04-.03-.04Z" fill="%s" stroke-width="3.82"/><path d="M191 58c.5 4.5-.3 13.5-1.5 19.5" stroke-width="4"/></g>`,
            color(colors, "hair"),
        )
    },
}

var registry = map[string]ComponentGroup{
    "base":       baseGroup,
    "hair":       hairGroup,
    "ears":       earsGroup,
    "earrings":   earringsGroup,
    "eyebrows":   eyebrowsGroup,
    "eyes":       eyesGroup,
    "glasses":    glassesGroup,
    "nose":       noseGroup,
    "mouth":      mouthGroup,
    "shirt":      shirtGroup,
}

func pick(groupName, variant string) (*ComponentPick, error) {
    g, ok := registry[groupName]
    if !ok {
        return nil, errors.New("group not found:"+groupName)
    }
    f, ok := g[variant]
    if !ok {
        return nil, errors.New("variant %s not found in group"+variant+ groupName)
    }
    return &ComponentPick{Name: variant, Value: f}, nil
}

func generate_SVG(opt AvatarOptions, colors ColorPickCollection) (string, error) {
    components := ComponentPickCollection{}

    var err error
    if components["base"], err = pick("base", opt.Base); err != nil {
        return "", err
    }
    if components["hair"], err = pick("hair", opt.Hair); err != nil {
        return "", err
    }
    if components["ears"], err = pick("ears", opt.Ears); err != nil {
        return "", err
    }
    if components["earrings"], err = pick("earrings", opt.Earrings); err != nil {
        return "", err
    }
    if components["eyebrows"], err = pick("eyebrows", opt.Eyebrows); err != nil {
        return "", err
    }
    if components["eyes"], err = pick("eyes", opt.Eyes); err != nil {
        return "", err
    }
    if opt.Glasses != "" {
        if components["glasses"], err = pick("glasses", opt.Glasses); err != nil {
            return "", err
        }
    }
    if components["nose"], err = pick("nose", opt.Nose); err != nil {
        return "", err
    }
    if components["mouth"], err = pick("mouth", opt.Mouth); err != nil {
        return "", err
    }
    if components["shirt"], err = pick("shirt", opt.Shirt); err != nil {
        return "", err
    }
    if opt.FacialHair != "" {
        if components["facialHair"], err = pick("facialHair", opt.FacialHair); err != nil {
            return "", err
        }
    }

    svg := sprintf(
        `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 360 360" fill="none" shape-rendering="auto">`+
            `<mask id="viewboxMask"><rect width="360" height="360" rx="0" ry="0" x="0" y="0" fill="#fff"/></mask>`+
            `<g mask="url(#viewboxMask)">`+
            `<g transform="translate(80 23)">%s</g>`+
            `<g transform="translate(170 183)">%s</g>`+
            `<g transform="translate(110 102)">%s</g>`+
            `<g transform="translate(49 11)">%s</g>`+
            `<g transform="translate(142 119)">%s</g>`+
            `<g transform="rotate(-8 1149.44 -1186.92)">%s</g>`+
            `<g transform="translate(84 154)">%s</g>`+
            `<g transform="translate(53 272)">%s</g>`+
            `</g></svg>`,
        components["base"].Value(components, colors),
        components["mouth"].Value(components, colors),
        components["eyebrows"].Value(components, colors),
        components["hair"].Value(components, colors),
        components["eyes"].Value(components, colors),
        components["nose"].Value(components, colors),
        components["ears"].Value(components, colors),
        components["shirt"].Value(components, colors),
    )

    return svg, nil
}

func NewAvatar(generate string) (string,error) {
	h := fnv.New64a()
	h.Write([]byte(generate))
	seed:=h.Sum64()
    pick := func(arr []string) string {
		return arr[seed % uint64(len(arr))]
    }
    basePalette := []string{"#f9c9b6", "#e69177", "#e1a57b", "#edb98a", "#ffdbb4"}
    earringPalette := []string{"#000000", "#C0C0C0", "#FFD700"}
    eyebrowPalette := []string{"#000000", "#3b2f2f", "#6b4f3f"}
    eyesPalette := []string{"#000000", "#2b2b2b", "#1f3a5f", "#2f5f1f"}
    eyeShadowPalette := []string{"#000000", "#6b5b95", "#a87ca0", "#8b5a2b"}
    glassesPalette := []string{"#000000", "#5a5a5a", "#8b4513", "#1f3a5f"}
    hairPalette := []string{"#ffeba4", "#2f1b0c", "#6d4c41", "#b5651d", "#d6b370", "#c0c0c0"}
    facialHairPalette := []string{"#5b4636", "#2f1b0c", "#6d4c41", "#8d6e63"}
    mouthPalette := []string{"#000000", "#7a1f1f", "#5a2a2a"}
    shirtPalette := []string{"#e0ddff", "#bde0fe", "#caffbf", "#ffd6a5", "#ffadad"}

    colors := ColorPickCollection{
        "base":       pick(basePalette),
        "earring":    pick(earringPalette),
        "eyebrows":   pick(eyebrowPalette),
        "eyes":       pick(eyesPalette),
        "eyeShadow":  pick(eyeShadowPalette),
        "glasses":    pick(glassesPalette),
        "hair":       pick(hairPalette),
        "facialHair": pick(facialHairPalette),
        "mouth":      pick(mouthPalette),
        "shirt":      pick(shirtPalette),
    }
    opt := AvatarOptions{
        Base:     pick([]string{"standard"}),
        Hair:     pick([]string{"fonze","mrT","dougFunny","mrClean","dannyPhantom","full","turban","pixie"}),
        Ears:     pick([]string{"attached", "detached"}),
        Earrings: pick([]string{"hoop", "stud"}),
        Eyebrows: pick([]string{"down", "eyelashesDown"}),
        Eyes:     pick([]string{"eyes", "smiling"}),
        Nose:     pick([]string{"pointed", "curve"}),
        Mouth:    pick([]string{"frown", "smile","surprised","laughing","nervous","sad","pucker","smirk"}),
        Shirt:    pick([]string{"crew"}),
    }
    if seed&1 == 1 {
        opt.Glasses = pick([]string{"round", "square"})
    }
    s, err := generate_SVG(opt, colors)
    if err != nil {
        return "",err
    }
    return s,nil
}