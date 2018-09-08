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
// TODO: Draw the animation
```

## Links
* [Spriter](https://brashmonkey.com)
* [SCML File format](http://www.brashmonkey.com/ScmlDocs/ScmlReference.html)
