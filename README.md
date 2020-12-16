# mandelbrot

Interactive Mandelbrot set visualizer.

Uses the [ebiten][3] game library to run an interactive window. To allow rendering in realtime with user input, some approximation is used to skip computations for some pixels, leading to frames improving in appearance as user input stops and re-rendering happens.

A still frame from the visualizer:
![still](https://user-images.githubusercontent.com/6401746/102289976-e16ea800-3ef4-11eb-9559-23161ad58e88.png)

Here is what the program looks like when used interactively, keep in mind the quality here is poor due to the low gif framerate and additional compression:
![live](https://user-images.githubusercontent.com/6401746/102290677-5bebf780-3ef6-11eb-90d6-be43bcebf90d.gif)

## Usage

To build and run from source:
```
make
```

If you are on MacOSX there is a [precompiled binary on GitHub][release] you can run:
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
[release]: https://github.com/keyan/mandelbrot/releases/tag/v1.0
