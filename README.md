# mandelbrot

Interactive Mandelbrot set visualizer. Uses the [ebiten][3] game library to run an interactive window.

## Usage

To build and run from source:
```
make
```

If you are on MacOSX there is a precompiled binary on GitHub you can run:
```
./mandelbrot
```

Controls are explained when the window first loads, but for completeness:
```
Arrow keys to move
I to zoom In
O to zoom Out
R to reset view
Escape to exit
```

## Resources

I found the youtube channel [fractalmath][1] to be helpful for better understanding complex plane dynamics. Lode Vandevenne also has a useful [tutorial][2] as well, but as with most of his articles it can be tough to follow. The ebiten [examples page][4] was invaluable in quickly using that library for the graphical/interactive portions.

[1]: https://www.youtube.com/channel/UCJ1i1TGHljQ6ETPgptchOZg
[2]: https://lodev.org/cgtutor/juliamandelbrot.html
[3]: https://ebiten.org/
[4]: https://ebiten.org/examples
