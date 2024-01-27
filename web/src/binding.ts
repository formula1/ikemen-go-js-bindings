import { GameLoop } from "./GameLoop"
class IkemenGoBinding {
  gameLoop = new GameLoop(60, ()=>{
    this.handleInput()
    this.logic()
    this.render();
  })

  constructor(){
    this.ensureStatsFile()
  }

  play(){
    this.gameLoop.play()
  }
  pause(){
    this.gameLoop.pause()
  }
  kill(){
    this.gameLoop.kill();
  }

}
