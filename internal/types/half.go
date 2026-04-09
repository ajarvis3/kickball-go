package types

type Half string

const (
    HalfTop    Half = "top"
    HalfBottom Half = "bottom"
)

func IsHalf(s string) bool {
    return s == string(HalfTop) || s == string(HalfBottom)
}
