package tty

import (
	"fmt"
	"github.com/toby1984/go_flocking/list"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"unsafe"
)

const (
	VT100_CLEAR_SCREEN = "\033[2J"
)

type Color int

// see https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
const (
	COLOR_BLACK   Color = 30
	COLOR_RED           = 31
	COLOR_GREEN         = 32
	COLOR_YELLOW        = 33
	COLOR_BLUE          = 34
	COLOR_MAGENTA       = 35
	COLOR_CYAN          = 36
	COLOR_WHITE         = 37
	COLOR_DEFAULT       = 39
	COLOR_RESET         = 0
)

const (
	COLOR_COUNT = 8
)


type WindowSizeChangedListener struct {
	callback func(int, int)
}

type Terminal struct {
	windowSizeChanged              chan os.Signal
	windowSizeChangedCallbackMutex sync.Mutex
	windowSizeChangedCallbacks     list.List[WindowSizeChangedListener]
	windowSizeMutex                sync.Mutex
	ttyColumns                     int
	ttyRows                        int
}

func (t *Terminal) AddSizeChangeListener(f func(int, int)) *WindowSizeChangedListener {
	var listenerFunc WindowSizeChangedListener
	listenerFunc.callback = f

	t.windowSizeChangedCallbackMutex.Lock()
	defer t.windowSizeChangedCallbackMutex.Unlock()
	t.windowSizeChangedCallbacks.Add(listenerFunc)
	return &listenerFunc
}

func (t *Terminal) RemoveSizeChangeListener(l *WindowSizeChangedListener) bool {
	if l == nil {
		panic("Listener must not be null")
	}
	t.windowSizeChangedCallbackMutex.Lock()
	defer t.windowSizeChangedCallbackMutex.Unlock()
	return t.windowSizeChangedCallbacks.Remove(func(value WindowSizeChangedListener) bool {
		return &value == l
	})
}

func (t *Terminal) GetColor(colIdx int) Color
{
		idx := colIdx % COLOR_COUNT

}

func (t *Terminal) GetTerminalSize() (int, int) {
	t.windowSizeMutex.Lock()
	defer t.windowSizeMutex.Unlock()
	return t.ttyColumns, t.ttyRows
}

func (t *Terminal) updateTerminalWidth() (int, int) {
	t.windowSizeMutex.Lock()
	defer t.windowSizeMutex.Unlock()
	columns, rows, err := terminalWidth()
	if err == nil {
		t.ttyRows = rows
		t.ttyColumns = columns
	}
	return t.ttyColumns, t.ttyRows
}

func (t *Terminal) MoveCursorTo(x int, y int) {
	os.Stdout.WriteString(fmt.Sprintf("\033[%d;%dH", y, x))
}

func (t *Terminal) PrintAt(text string, x int, y int) {
	os.Stdout.WriteString(fmt.Sprintf("\033[%d;%dH%s", y, x, text))
}

func (t *Terminal) PrintAtWithColor(text string, x int, y int, color Color) {
	os.Stdout.WriteString(fmt.Sprintf("\033[%d;%dH\u001B[%dm%s", y, x, color, text))
}

func (t *Terminal) PrintWithColor(text string, color Color) {
	// \033[XXXm
	os.Stdout.WriteString(fmt.Sprintf("\u001B[%dm%s", color, text))
}

func (t *Terminal) ClearScreen() {
	os.Stdout.WriteString(VT100_CLEAR_SCREEN)
}

func NewTerminal() *Terminal {
	var t Terminal
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	t.ttyColumns = -1
	t.ttyRows = -1
	t.updateTerminalWidth()
	t.windowSizeChanged = make(chan os.Signal, 1)
	signal.Notify(t.windowSizeChanged, syscall.SIGWINCH)
	go func() {
		for {
			<-t.windowSizeChanged
			fmt.Println("DEBUG: Terminal size changed")
			var cols, rows = t.updateTerminalWidth()

			t.windowSizeChangedCallbackMutex.Lock()
			t.windowSizeChangedCallbacks.Visit(func(l WindowSizeChangedListener) {
				l.callback(cols, rows)
			})
			t.windowSizeChangedCallbackMutex.Unlock()
		}
	}()
	return &t
}

func (t *Terminal) SetForeground(color Color) {
	os.Stdout.WriteString(fmt.Sprintf("\033[%dm", color))
}

func (t *Terminal) SetBackground(color Color) {
	os.Stdout.WriteString(fmt.Sprintf("\033[48;5;%dm", color+10))
}

type window struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

func terminalWidth() (int, int, error) {
	w := new(window)
	tio := syscall.TIOCGWINSZ
	res, _, err := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(tio),
		uintptr(unsafe.Pointer(w)),
	)
	if int(res) == -1 {
		return 0, 0, err
	}
	return int(w.Col), int(w.Row), nil
}
