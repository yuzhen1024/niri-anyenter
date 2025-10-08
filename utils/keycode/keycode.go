package keycode

// https://github.com/torvalds/linux/blob/master/include/uapi/linux/input-event-codes.h

var Letters = map[int64]string{
	16: "q", 17: "w", 18: "e", 19: "r", 20: "t", 21: "y", 22: "u", 23: "i", 24: "o", 25: "p",
	30: "a", 31: "s", 32: "d", 33: "f", 34: "g", 35: "h", 36: "j", 37: "k", 38: "l",
	44: "z", 45: "x", 46: "c", 47: "v", 48: "b", 49: "n", 50: "m",
}
