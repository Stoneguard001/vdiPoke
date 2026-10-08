package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	inputKeyboard = 1
	keyEventKeyUp = 0x0002
	vkTab         = 0x09
	vkMenu        = 0x12
	vkF15         = 0x7E
)

var procSendInput = syscall.NewLazyDLL("user32.dll").NewProc("SendInput")

type keyboardInput struct {
	vk        uint16
	scan      uint16
	flags     uint32
	time      uint32
	extraInfo uintptr
}

// Padded to match the size of the Win32 INPUT union, whose largest member is MOUSEINPUT.
type input struct {
	inputType uint32
	ki        keyboardInput
	_         [8]byte
}

func key(vk uint16, flags uint32) input {
	return input{inputType: inputKeyboard, ki: keyboardInput{vk: vk, flags: flags}}
}

func sendF15() error {
	return send(
		key(vkF15, 0),
		key(vkF15, keyEventKeyUp),
	)
}

func sendAltTab() error {
	return send(
		key(vkMenu, 0),
		key(vkTab, 0),
		key(vkTab, keyEventKeyUp),
		key(vkMenu, keyEventKeyUp),
	)
}

func send(inputs ...input) error {
	n, _, err := procSendInput.Call(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		unsafe.Sizeof(inputs[0]),
	)
	if int(n) != len(inputs) {
		return fmt.Errorf("SendInput sent %d of %d events: %w", n, len(inputs), err)
	}
	return nil
}
