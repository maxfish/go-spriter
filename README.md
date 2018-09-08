# go-spriter
A simple Go importer and player for Spriter animation files (SCML). https://brashmonkey.com

## Usage

```go
// [...] Game init
spriterModel := spriter.NewSpriterModelFromFile("assets/hero/player.scml")
player := spriter.MakePlayer(spriterModel.GetEntityByName("Player"))
player.SetAnimationByName('run')

// [...] Game loop
player.Update(deltaMillisec)

// [...] Game render
// Draw the animation
```
For an example of how to draw the sprite, have a look at [drawer_example.go](https://github.com/maxfish/go-spriter/blob/master/drawer_example.go)

## Links
* [Spriter](https://brashmonkey.com)
* [SCML File format](http://www.brashmonkey.com/ScmlDocs/ScmlReference.html)
* [C+ Implementation](https://github.com/lucidspriter/SpriterPlusPlus)
* [Java implementation (Trixt0r)](https://github.com/Trixt0r/spriter)
