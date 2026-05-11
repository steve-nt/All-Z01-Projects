
package main

import (
    "testing"
    "banner/banner"
)

func TestAsciiArtsGenerator(t *testing.T) {
    args := []string{"--color=red", "test", "thinkertoy.txt"}
    
    // Call the main function Ascii_Arts_Generator with simulated command-line arguments
    Ascii_Arts_Generator(args)
}

func TestBannerFunctions(t *testing.T) {
    // Testing a basic banner function to ensure import is working and the function behaves as expected
    color := banner.SetColor("green")
    if color != "\033[32m" {
        t.Errorf("banner.SetColor('green') = %q; expected %q", color, "\033[32m")
    }
}
